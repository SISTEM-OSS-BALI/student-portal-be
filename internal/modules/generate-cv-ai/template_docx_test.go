package generatecvai

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGenerateCVDocument_MergesTemplateAndUsesStudentName(t *testing.T) {
	t.Parallel()

	templateBytes := mustLoadTemplateDocx(t)
	templateBase64 := base64.StdEncoding.EncodeToString(templateBytes)
	studentName := "Ngurah mNaik Mahardika"
	studentID := "student-123"

	fileBase64, fileName, err := generateCVDocument(context.Background(), GenerateDTO{
		StudentID:      &studentID,
		StudentName:    &studentName,
		TemplateBase64: &templateBase64,
		Sections: []GenerateSectionDTO{
			{
				Key:   "personal",
				Label: "Personal Details",
				Items: []GenerateSectionItem{
					{Question: "First Name", Answer: "Ngurah"},
					{Question: "Last Name", Answer: "mNaik Mahardika"},
					{Question: "Current Address", Answer: "Jl. Sunset Road No. 10"},
					{Question: "City", Answer: "Denpasar"},
					{Question: "Province/State", Answer: "Bali"},
					{Question: "Postal Code", Answer: "80361"},
					{Question: "Country", Answer: "Indonesia"},
					{Question: "Date of Birth", Answer: "2001-01-01"},
					{Question: "Place of Birth", Answer: "Denpasar"},
					{Question: "Contact Number", Answer: "08123456789"},
					{Question: "Email Address", Answer: "ngurah@example.com"},
				},
			},
			{
				Key:   "formal-education",
				Label: "Formal Education",
				Items: []GenerateSectionItem{
					{Question: "Level 1", Answer: "Bachelor"},
					{Question: "School 1", Answer: "Universitas XYZ"},
					{Question: "Commence 1", Answer: "2022"},
					{Question: "Completed 1", Answer: "2026"},
				},
			},
			{
				Key:   "employment-history",
				Label: "Employment History",
				Items: []GenerateSectionItem{
					{Question: "Company 1", Answer: "PT Contoh"},
					{Question: "Type 1", Answer: "Internship"},
					{Question: "Position 1", Answer: "Frontend Developer"},
					{Question: "From 1", Answer: "2025"},
					{Question: "To 1", Answer: "2025"},
					{Question: "Contact Reference 1", Answer: "Budi"},
				},
			},
		},
	}, "")
	if err != nil {
		t.Fatalf("generateCVDocument returned error: %v", err)
	}

	if got, want := fileName, "Ngurah_mNaik_Mahardika_CV.docx"; got != want {
		t.Fatalf("generated file name mismatch: got %q want %q", got, want)
	}

	docxBytes, err := base64.StdEncoding.DecodeString(fileBase64)
	if err != nil {
		t.Fatalf("decode generated base64: %v", err)
	}

	documentXML := readDocumentXMLFromDocx(t, docxBytes)
	for _, expected := range []string{
		"Ngurah",
		"mNaik Mahardika",
		"Jl. Sunset Road No. 10",
		"Denpasar",
		"Universitas XYZ",
		"PT Contoh",
		"Frontend Developer",
	} {
		if !strings.Contains(documentXML, expected) {
			t.Fatalf("generated document.xml missing %q", expected)
		}
	}

	if strings.Contains(documentXML, "altChunk") {
		t.Fatalf("generated document.xml still contains altChunk markup")
	}
}

func mustLoadTemplateDocx(t *testing.T) []byte {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	templatePath := filepath.Clean(filepath.Join(
		filepath.Dir(currentFile),
		"..", "..", "..", "..",
		"student-portal-fe", "public", "assets", "file", "Template CV.docx",
	))

	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("read template docx: %v", err)
	}
	return content
}

func readDocumentXMLFromDocx(t *testing.T, docxBytes []byte) string {
	t.Helper()

	reader, err := zip.NewReader(bytes.NewReader(docxBytes), int64(len(docxBytes)))
	if err != nil {
		t.Fatalf("open generated docx: %v", err)
	}

	for _, file := range reader.File {
		if file.Name != "word/document.xml" {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			t.Fatalf("open word/document.xml: %v", err)
		}
		defer rc.Close()

		content, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("read word/document.xml: %v", err)
		}
		return string(content)
	}

	t.Fatal("word/document.xml not found in generated docx")
	return ""
}
