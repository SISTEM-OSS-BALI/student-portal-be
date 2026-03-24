package generatecvai

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	tableXMLPattern = regexp.MustCompile(`(?s)<w:tbl>.*?</w:tbl>`)
	rowXMLPattern   = regexp.MustCompile(`(?s)<w:tr\b.*?</w:tr>`)
	cellXMLPattern  = regexp.MustCompile(`(?s)<w:tc\b.*?</w:tc>`)
	tcPrXMLPattern  = regexp.MustCompile(`(?s)<w:tcPr>.*?</w:tcPr>`)
)

func mergeCVTemplate(ctx context.Context, input GenerateDTO, aiSummary string) ([]byte, error) {
	_ = aiSummary

	templateBytes, err := resolveTemplateDocxBytes(ctx, input)
	if err != nil {
		return nil, err
	}

	data := extractCVDocumentData(input)
	return mergeTemplateDocxBytes(templateBytes, data)
}

func resolveTemplateDocxBytes(ctx context.Context, input GenerateDTO) ([]byte, error) {
	if input.TemplateBase64 != nil && strings.TrimSpace(*input.TemplateBase64) != "" {
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(*input.TemplateBase64))
		if err != nil {
			return nil, fmt.Errorf("decode template_base64: %w", err)
		}
		return decoded, nil
	}

	if path := strings.TrimSpace(stringPtrValue(input.TemplatePath)); path != "" {
		if content, err := os.ReadFile(path); err == nil {
			return content, nil
		}

		candidates := []string{
			filepath.Join(".", path),
			filepath.Join(".", "public", path),
			filepath.Join("..", "student-portal-fe", "public", path),
		}
		for _, candidate := range candidates {
			candidate = filepath.Clean(candidate)
			if content, err := os.ReadFile(candidate); err == nil {
				return content, nil
			}
		}
	}

	if templateURL := strings.TrimSpace(stringPtrValue(input.TemplateURL)); templateURL != "" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, templateURL, nil)
		if err != nil {
			return nil, fmt.Errorf("build template request: %w", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetch template_url: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("fetch template_url: status=%d", resp.StatusCode)
		}

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read template_url: %w", err)
		}
		return content, nil
	}

	return nil, fmt.Errorf("template docx tidak tersedia")
}

func mergeTemplateDocxBytes(templateBytes []byte, data cvDocumentData) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(templateBytes), int64(len(templateBytes)))
	if err != nil {
		return nil, fmt.Errorf("open template docx: %w", err)
	}

	files := make(map[string][]byte, len(reader.File))
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("open template entry %s: %w", file.Name, err)
		}
		content, readErr := io.ReadAll(rc)
		closeErr := rc.Close()
		if readErr != nil {
			return nil, fmt.Errorf("read template entry %s: %w", file.Name, readErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("close template entry %s: %w", file.Name, closeErr)
		}
		files[file.Name] = content
	}

	documentXML, ok := files["word/document.xml"]
	if !ok {
		return nil, fmt.Errorf("template docx missing word/document.xml")
	}

	mergedDocumentXML, err := mergeTemplateDocumentXML(string(documentXML), data)
	if err != nil {
		return nil, err
	}
	files["word/document.xml"] = []byte(mergedDocumentXML)

	var output bytes.Buffer
	writer := zip.NewWriter(&output)
	for _, file := range reader.File {
		header := file.FileHeader
		w, err := writer.CreateHeader(&header)
		if err != nil {
			return nil, fmt.Errorf("create output entry %s: %w", file.Name, err)
		}
		if _, err := w.Write(files[file.Name]); err != nil {
			return nil, fmt.Errorf("write output entry %s: %w", file.Name, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("finalize docx: %w", err)
	}

	return output.Bytes(), nil
}

func mergeTemplateDocumentXML(documentXML string, data cvDocumentData) (string, error) {
	tableIndices := tableXMLPattern.FindAllStringIndex(documentXML, -1)
	if len(tableIndices) < 7 {
		return "", fmt.Errorf("template table structure mismatch: found %d tables", len(tableIndices))
	}

	tableContents := make([]string, 0, len(tableIndices))
	for _, idx := range tableIndices {
		tableContents = append(tableContents, documentXML[idx[0]:idx[1]])
	}

	tableContents[0] = fillPersonalTableXML(tableContents[0], data.personal)
	tableContents[1] = fillGridTableXML(tableContents[1], data.formal, []string{"level", "school", "commence", "completed"})
	tableContents[2] = fillGridTableXML(tableContents[2], data.nonformal, []string{"level", "school", "commence", "completed"})
	tableContents[3] = fillGridTableXML(tableContents[3], data.employment, []string{"company", "type", "position", "from", "to", "contact_reference"})
	tableContents[4] = fillGridTableXML(tableContents[4], data.achievement, []string{"company", "type", "position", "from", "to"})
	tableContents[5] = fillGridTableXML(tableContents[5], data.travel, []string{"country", "purpose", "start", "end"})
	tableContents[6] = fillNextOfKinTableXML(tableContents[6], data.nextOfKin)

	var builder strings.Builder
	last := 0
	for index, idx := range tableIndices {
		builder.WriteString(documentXML[last:idx[0]])
		builder.WriteString(tableContents[index])
		last = idx[1]
	}
	builder.WriteString(documentXML[last:])

	return builder.String(), nil
}

func fillPersonalTableXML(tableXML string, personal map[string]string) string {
	tableXML = setTableRowCellText(tableXML, 1, 1, personal["first_name"])
	tableXML = setTableRowCellText(tableXML, 1, 3, personal["last_name"])
	tableXML = setTableRowCellText(tableXML, 2, 1, personal["current_address"])
	tableXML = setTableRowCellText(tableXML, 3, 1, personal["city"])
	tableXML = setTableRowCellText(tableXML, 3, 3, personal["province"])
	tableXML = setTableRowCellText(tableXML, 4, 1, personal["postal_code"])
	tableXML = setTableRowCellText(tableXML, 4, 3, personal["country"])
	tableXML = setTableRowCellText(tableXML, 5, 1, personal["date_of_birth"])
	tableXML = setTableRowCellText(tableXML, 5, 3, personal["place_of_birth"])
	tableXML = setTableRowCellText(tableXML, 6, 1, personal["age"])
	tableXML = setTableRowCellText(tableXML, 6, 3, personal["nationality"])
	tableXML = setTableRowCellText(tableXML, 7, 1, personal["gender"])
	tableXML = setTableRowCellText(tableXML, 7, 3, personal["marital_status"])
	tableXML = setTableRowCellText(tableXML, 8, 1, personal["height"])
	tableXML = setTableRowCellText(tableXML, 8, 3, personal["weight"])
	tableXML = setTableRowCellText(tableXML, 9, 1, personal["id_number"])
	tableXML = setTableRowCellText(tableXML, 10, 1, personal["passport_number"])
	tableXML = setTableRowCellText(tableXML, 11, 1, personal["contact_number"])
	tableXML = setTableRowCellText(tableXML, 12, 1, personal["email_address"])
	return tableXML
}

func fillNextOfKinTableXML(tableXML string, data map[string]string) string {
	tableXML = setTableRowCellText(tableXML, 1, 1, data["name"])
	tableXML = setTableRowCellText(tableXML, 1, 3, data["relationship"])
	tableXML = setTableRowCellText(tableXML, 2, 1, data["address"])
	tableXML = setTableRowCellText(tableXML, 3, 1, data["city"])
	tableXML = setTableRowCellText(tableXML, 3, 3, data["province"])
	tableXML = setTableRowCellText(tableXML, 4, 1, data["postal_code"])
	tableXML = setTableRowCellText(tableXML, 4, 3, data["country"])
	tableXML = setTableRowCellText(tableXML, 5, 1, data["contact_no"])
	tableXML = setTableRowCellText(tableXML, 5, 3, data["email"])
	return tableXML
}

func fillGridTableXML(tableXML string, rows []cvRow, keys []string) string {
	for rowIndex, rowData := range rows {
		templateRowIndex := rowIndex + 2
		for keyIndex, key := range keys {
			tableXML = setTableRowCellText(tableXML, templateRowIndex, keyIndex+1, rowData[key])
		}
	}
	return tableXML
}

func setTableRowCellText(tableXML string, rowIndex int, cellIndex int, value string) string {
	if strings.TrimSpace(value) == "" {
		return tableXML
	}

	rowMatches := rowXMLPattern.FindAllStringIndex(tableXML, -1)
	if rowIndex < 0 || rowIndex >= len(rowMatches) {
		return tableXML
	}

	rowXML := tableXML[rowMatches[rowIndex][0]:rowMatches[rowIndex][1]]
	updatedRowXML := setRowCellText(rowXML, cellIndex, value)
	if updatedRowXML == rowXML {
		return tableXML
	}

	return tableXML[:rowMatches[rowIndex][0]] + updatedRowXML + tableXML[rowMatches[rowIndex][1]:]
}

func setRowCellText(rowXML string, cellIndex int, value string) string {
	cellMatches := cellXMLPattern.FindAllStringIndex(rowXML, -1)
	if cellIndex < 0 || cellIndex >= len(cellMatches) {
		return rowXML
	}

	cellXML := rowXML[cellMatches[cellIndex][0]:cellMatches[cellIndex][1]]
	updatedCellXML := rebuildCellXML(cellXML, value)
	if updatedCellXML == cellXML {
		return rowXML
	}

	return rowXML[:cellMatches[cellIndex][0]] + updatedCellXML + rowXML[cellMatches[cellIndex][1]:]
}

func rebuildCellXML(cellXML string, value string) string {
	if strings.TrimSpace(value) == "" {
		return cellXML
	}

	openTagEnd := strings.Index(cellXML, ">")
	closeTagStart := strings.LastIndex(cellXML, "</w:tc>")
	if openTagEnd == -1 || closeTagStart == -1 {
		return cellXML
	}

	var builder strings.Builder
	builder.WriteString(cellXML[:openTagEnd+1])

	if tcPr := tcPrXMLPattern.FindString(cellXML); tcPr != "" {
		builder.WriteString(tcPr)
	}

	builder.WriteString(buildCellParagraphsXML(value))
	builder.WriteString(`</w:tc>`)

	return builder.String()
}

func buildCellParagraphsXML(value string) string {
	lines := strings.Split(strings.TrimSpace(value), "\n")
	if len(lines) == 0 {
		return `<w:p/>`
	}

	var builder strings.Builder
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			builder.WriteString(`<w:p/>`)
			continue
		}

		builder.WriteString(`<w:p><w:pPr><w:widowControl w:val="0"/></w:pPr><w:r><w:t xml:space="preserve">`)
		builder.WriteString(html.EscapeString(line))
		builder.WriteString(`</w:t></w:r></w:p>`)
	}

	return builder.String()
}

func encodeFileToBase64(content []byte) string {
	return base64.StdEncoding.EncodeToString(content)
}
