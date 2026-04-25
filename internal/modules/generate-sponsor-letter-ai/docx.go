package generatesponsorletterai

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html"
	"strings"
)

const generatedWordMimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

func generateSponsorLetterDocument(_ context.Context, input GenerateDTO, content string) (string, string, error) {
	docxBytes, err := buildSponsorLetterDocxBytes(input, content)
	if err != nil {
		return "", "", err
	}

	fileName := buildGeneratedSponsorLetterFileName(input)
	return base64.StdEncoding.EncodeToString(docxBytes), fileName, nil
}

func buildGeneratedSponsorLetterFileName(input GenerateDTO) string {
	name := sanitizeFileName(stringPtrValue(input.StudentName))
	if name == "" {
		name = "student"
	}

	return fmt.Sprintf("%s_Sponsor_Letter.docx", name)
}

func buildSponsorLetterDocxBytes(input GenerateDTO, content string) ([]byte, error) {
	var output bytes.Buffer
	writer := zip.NewWriter(&output)

	files := map[string]string{
		"[Content_Types].xml":          buildSponsorLetterContentTypesXML(),
		"_rels/.rels":                  buildSponsorLetterRootRelsXML(),
		"word/document.xml":            buildSponsorLetterDocumentXML(input, content),
		"word/_rels/document.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`,
	}

	for name, body := range files {
		entry, err := writer.Create(name)
		if err != nil {
			return nil, fmt.Errorf("create docx entry %s: %w", name, err)
		}
		if _, err := entry.Write([]byte(body)); err != nil {
			return nil, fmt.Errorf("write docx entry %s: %w", name, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("finalize sponsor letter docx: %w", err)
	}

	return output.Bytes(), nil
}

func buildSponsorLetterContentTypesXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`
}

func buildSponsorLetterRootRelsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`
}

func buildSponsorLetterDocumentXML(input GenerateDTO, content string) string {
	title := "Sponsor Letter"
	if name := strings.TrimSpace(stringPtrValue(input.StudentName)); name != "" {
		title = fmt.Sprintf("Sponsor Letter - %s", name)
	}

	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	builder.WriteString(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">`)
	builder.WriteString(`<w:body>`)
	builder.WriteString(`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="30"/></w:rPr><w:t xml:space="preserve">`)
	builder.WriteString(html.EscapeString(title))
	builder.WriteString(`</w:t></w:r></w:p>`)
	builder.WriteString(`<w:p/>`)
	builder.WriteString(buildSponsorLetterParagraphsXML(content))
	builder.WriteString(`<w:sectPr><w:pgSz w:w="12240" w:h="15840"/><w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="720" w:footer="720" w:gutter="0"/></w:sectPr>`)
	builder.WriteString(`</w:body></w:document>`)
	return builder.String()
}

func buildSponsorLetterParagraphsXML(content string) string {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	var builder strings.Builder
	hasContent := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			builder.WriteString(`<w:p/>`)
			continue
		}

		hasContent = true
		builder.WriteString(`<w:p><w:pPr><w:spacing w:after="180" w:line="360" w:lineRule="auto"/></w:pPr><w:r><w:t xml:space="preserve">`)
		builder.WriteString(html.EscapeString(line))
		builder.WriteString(`</w:t></w:r></w:p>`)
	}

	if !hasContent {
		builder.WriteString(`<w:p><w:r><w:t xml:space="preserve">Sponsor letter content is empty.</w:t></w:r></w:p>`)
	}

	return builder.String()
}

func stringPtrValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func sanitizeFileName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	var b strings.Builder
	for _, r := range value {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune('_')
		}
	}

	return strings.Trim(b.String(), "_")
}
