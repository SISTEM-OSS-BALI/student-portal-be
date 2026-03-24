package generatecvai

import (
	"context"
	"fmt"
	"html"
	"regexp"
	"strings"
)

const generatedWordMimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

type cvRow map[string]string

type cvDocumentData struct {
	personal    map[string]string
	formal      []cvRow
	nonformal   []cvRow
	employment  []cvRow
	achievement []cvRow
	travel      []cvRow
	nextOfKin   map[string]string
}

type rowFieldSchema struct {
	key     string
	aliases []string
}

var rowIndexPattern = regexp.MustCompile(`(?:^|[^0-9])([1-9][0-9]?)(?:[^0-9]|$)`)

func generateCVDocument(ctx context.Context, input GenerateDTO, aiSummary string) (string, string, error) {
	docxBytes, err := mergeCVTemplate(ctx, input, aiSummary)
	if err != nil {
		return "", "", err
	}

	fileName := buildGeneratedWordFileName(input)
	return encodeFileToBase64(docxBytes), fileName, nil
}

func buildGeneratedWordFileName(input GenerateDTO) string {
	name := sanitizeFileName(stringPtrValue(input.StudentName))
	if name == "" {
		name = sanitizeFileName(guessStudentName(input))
	}
	if name == "" {
		name = "student"
	}

	return fmt.Sprintf("%s_CV.docx", name)
}

func buildCVHTML(input GenerateDTO, aiSummary string) string {
	_ = aiSummary

	data := extractCVDocumentData(input)

	var b strings.Builder
	b.WriteString(`<!DOCTYPE html><html xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:w="urn:schemas-microsoft-com:office:word" xmlns="http://www.w3.org/TR/REC-html40"><head><meta charset="utf-8"><meta http-equiv="Content-Type" content="text/html; charset=utf-8"><style>`)
	b.WriteString(`@page{size:A4;margin:14mm 12mm;}body{font-family:"Times New Roman",serif;color:#000;font-size:11pt;line-height:1.15;} .page{width:820px;margin:0 auto;} h1{margin:10px 0 26px;text-align:center;font-size:29pt;font-weight:400;letter-spacing:0.2px;} table{border-collapse:collapse;table-layout:fixed;margin:0 0 28px 46px;} td{border:1.2px solid #111;padding:5px 8px;vertical-align:middle;} .section-title{background:#d9d9d9;font-weight:700;font-size:11.7pt;padding:8px 11px;text-align:left;} .label{background:#efefef;font-weight:700;} .value{background:#fff;} .top{vertical-align:top;} .center{text-align:center;} .small{font-size:10.8pt;} </style></head><body><div class="page">`)
	b.WriteString(`<h1>CURRICULUM VITAE</h1>`)
	b.WriteString(renderPersonalDetailsTable(data.personal))
	b.WriteString(renderSimpleGridTable("2. Formal Education", []string{"No.", "Level", "School", "Commence", "Completed"}, []int{14, 14, 18, 27, 27}, data.formal, 3, 438))
	b.WriteString(renderSimpleGridTable("3. Nonformal Education", []string{"No.", "Level", "School", "Commence", "Completed"}, []int{14, 14, 18, 27, 27}, data.nonformal, 1, 446))
	b.WriteString(renderSimpleGridTable("4. Employment History", []string{"No.", "Company", "Type", "Position", "From", "To", "Contact Reference"}, []int{8, 22, 12, 16, 10, 10, 22}, data.employment, 3, 620))
	b.WriteString(renderSimpleGridTable("5. Achievement", []string{"No.", "Company", "Type", "Position", "From", "To"}, []int{10, 26, 14, 24, 13, 13}, data.achievement, 2, 560))
	b.WriteString(renderSimpleGridTable("6. Travel History", []string{"No.", "Country", "Purpose", "Start", "End"}, []int{10, 22, 30, 19, 19}, data.travel, 1, 470))
	b.WriteString(renderNextOfKinTable(data.nextOfKin))
	b.WriteString(`</div></body></html>`)
	return b.String()
}

func renderPersonalDetailsTable(personal map[string]string) string {
	var b strings.Builder
	b.WriteString(`<table style="width:533px;" border="1" cellspacing="0" cellpadding="0"><colgroup><col style="width:19%"><col style="width:13%"><col style="width:36%"><col style="width:13%"><col style="width:19%"></colgroup>`)
	b.WriteString(`<tr><td colspan="5" class="section-title">1. Personal Details</td></tr>`)
	b.WriteString(`<tr><td class="label">First Name</td><td colspan="2" class="value">` + esc(personal["first_name"]) + `</td><td class="label">Last Name</td><td class="value">` + esc(personal["last_name"]) + `</td></tr>`)
	b.WriteString(`<tr><td rowspan="3" class="label top">Current Address</td><td colspan="4" class="value top" style="height:76px;">` + esc(personal["current_address"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">City</td><td class="value">` + esc(personal["city"]) + `</td><td class="label">Province/State</td><td class="value">` + esc(personal["province"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">Postal Code</td><td class="value">` + esc(personal["postal_code"]) + `</td><td class="label">Country</td><td class="value">` + esc(personal["country"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">Date of Birth</td><td colspan="2" class="value">` + esc(personal["date_of_birth"]) + `</td><td class="label">Place of Birth</td><td class="value">` + esc(personal["place_of_birth"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">Age</td><td colspan="2" class="value">` + esc(personal["age"]) + `</td><td class="label">Nationality</td><td class="value">` + esc(personal["nationality"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">Gender</td><td colspan="2" class="value">` + esc(personal["gender"]) + `</td><td class="label">Marital Status</td><td class="value">` + esc(personal["marital_status"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">Height</td><td colspan="2" class="value">` + esc(personal["height"]) + `</td><td class="label">Weight</td><td class="value">` + esc(personal["weight"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">ID Number</td><td colspan="4" class="value">` + esc(personal["id_number"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">Passport Number</td><td colspan="4" class="value">` + esc(personal["passport_number"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">Contact Number</td><td colspan="4" class="value">` + esc(personal["contact_number"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">Email Address</td><td colspan="4" class="value">` + esc(personal["email_address"]) + `</td></tr>`)
	b.WriteString(`</table>`)
	return b.String()
}

func renderSimpleGridTable(title string, headers []string, widths []int, rows []cvRow, minRows int, tableWidth int) string {
	if len(rows) < minRows {
		padding := make([]cvRow, 0, minRows-len(rows))
		for i := len(rows); i < minRows; i++ {
			padding = append(padding, cvRow{})
		}
		rows = append(rows, padding...)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<table style="width:%dpx;" border="1" cellspacing="0" cellpadding="0"><colgroup>`, tableWidth))
	for _, width := range widths {
		b.WriteString(fmt.Sprintf(`<col style="width:%d%%">`, width))
	}
	b.WriteString(`</colgroup>`)
	b.WriteString(`<tr><td colspan="` + fmt.Sprintf("%d", len(headers)) + `" class="section-title">` + esc(title) + `</td></tr>`)
	b.WriteString(`<tr>`)
	for _, header := range headers {
		b.WriteString(`<td class="label center nowrap small">` + esc(header) + `</td>`)
	}
	b.WriteString(`</tr>`)
	for i, row := range rows {
		b.WriteString(`<tr>`)
		for idx, header := range headers {
			if idx == 0 {
				b.WriteString(`<td class="center value" style="height:38px;">` + esc(fmt.Sprintf("%d.", i+1)) + `</td>`)
				continue
			}
			key := normalizeKey(header)
			b.WriteString(`<td class="value top" style="height:38px;">` + esc(row[key]) + `</td>`)
		}
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</table>`)
	return b.String()
}

func renderNextOfKinTable(data map[string]string) string {
	var b strings.Builder
	b.WriteString(`<table style="width:533px;" border="1" cellspacing="0" cellpadding="0"><colgroup><col style="width:19%"><col style="width:13%"><col style="width:36%"><col style="width:13%"><col style="width:19%"></colgroup>`)
	b.WriteString(`<tr><td colspan="5" class="section-title">7. Next of Kin/Parents</td></tr>`)
	b.WriteString(`<tr><td class="label">Name</td><td colspan="2" class="value">` + esc(data["name"]) + `</td><td class="label">Relationship</td><td class="value">` + esc(data["relationship"]) + `</td></tr>`)
	b.WriteString(`<tr><td rowspan="3" class="label top">Address</td><td colspan="4" class="value top" style="height:58px;">` + esc(data["address"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">City</td><td class="value">` + esc(data["city"]) + `</td><td class="label">Province</td><td class="value">` + esc(data["province"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">Postal Code</td><td class="value">` + esc(data["postal_code"]) + `</td><td class="label">Country</td><td class="value">` + esc(data["country"]) + `</td></tr>`)
	b.WriteString(`<tr><td class="label">Contact No.</td><td colspan="2" class="value">` + esc(data["contact_no"]) + `</td><td class="label">Email</td><td class="value">` + esc(data["email"]) + `</td></tr>`)
	b.WriteString(`</table>`)
	return b.String()
}

func extractCVDocumentData(input GenerateDTO) cvDocumentData {
	data := cvDocumentData{
		personal:  make(map[string]string),
		nextOfKin: make(map[string]string),
	}

	allItems := make([]GenerateSectionItem, 0)
	for _, section := range input.Sections {
		allItems = append(allItems, section.Items...)

		sectionKey := normalizeLabel(section.Label + " " + section.Key)
		switch {
		case isPersonalSection(sectionKey):
			assignPersonalFromItems(data.personal, section.Items)
		case isFormalEducationSection(sectionKey):
			data.formal = buildRows(section.Items, []rowFieldSchema{
				{key: "level", aliases: []string{"level", "jenjang"}},
				{key: "school", aliases: []string{"school", "sekolah", "university", "campus", "institution"}},
				{key: "commence", aliases: []string{"commence", "start", "mulai", "from", "tahun masuk"}},
				{key: "completed", aliases: []string{"completed", "finish", "end", "selesai", "graduation", "tahun selesai"}},
			})
		case isNonformalEducationSection(sectionKey):
			data.nonformal = buildRows(section.Items, []rowFieldSchema{
				{key: "level", aliases: []string{"level", "jenjang", "course"}},
				{key: "school", aliases: []string{"school", "provider", "institution", "course", "academy"}},
				{key: "commence", aliases: []string{"commence", "start", "mulai", "from"}},
				{key: "completed", aliases: []string{"completed", "finish", "end", "selesai"}},
			})
		case isEmploymentSection(sectionKey):
			data.employment = buildRows(section.Items, []rowFieldSchema{
				{key: "company", aliases: []string{"company", "perusahaan", "employer"}},
				{key: "type", aliases: []string{"type", "jenis", "industry", "bidang"}},
				{key: "position", aliases: []string{"position", "jabatan", "role"}},
				{key: "from", aliases: []string{"from", "start", "mulai"}},
				{key: "to", aliases: []string{"to", "end", "selesai"}},
				{key: "contact_reference", aliases: []string{"contact reference", "reference", "referensi", "supervisor"}},
			})
		case isAchievementSection(sectionKey):
			data.achievement = buildRows(section.Items, []rowFieldSchema{
				{key: "company", aliases: []string{"company", "organizer", "issuer", "penyelenggara"}},
				{key: "type", aliases: []string{"type", "jenis", "category"}},
				{key: "position", aliases: []string{"position", "title", "achievement", "prestasi", "award"}},
				{key: "from", aliases: []string{"from", "start", "mulai"}},
				{key: "to", aliases: []string{"to", "end", "selesai", "year"}},
			})
		case isTravelSection(sectionKey):
			data.travel = buildRows(section.Items, []rowFieldSchema{
				{key: "country", aliases: []string{"country", "negara"}},
				{key: "purpose", aliases: []string{"purpose", "tujuan"}},
				{key: "start", aliases: []string{"start", "from", "mulai"}},
				{key: "end", aliases: []string{"end", "to", "selesai"}},
			})
		case isNextOfKinSection(sectionKey):
			assignNextOfKinFromItems(data.nextOfKin, section.Items)
		}
	}

	if data.personal["first_name"] == "" && data.personal["last_name"] == "" {
		assignPersonalFromItems(data.personal, allItems)
	}
	if len(data.formal) == 0 {
		data.formal = buildRows(allItems, []rowFieldSchema{
			{key: "level", aliases: []string{"education level", "level", "jenjang"}},
			{key: "school", aliases: []string{"school", "sekolah", "university", "campus"}},
			{key: "commence", aliases: []string{"commence", "start", "mulai", "tahun masuk"}},
			{key: "completed", aliases: []string{"completed", "graduation", "tahun selesai", "end"}},
		})
	}
	if len(data.nonformal) == 0 {
		data.nonformal = buildRows(allItems, []rowFieldSchema{
			{key: "level", aliases: []string{"course level", "nonformal level", "level"}},
			{key: "school", aliases: []string{"course", "provider", "institution", "school"}},
			{key: "commence", aliases: []string{"commence", "start", "mulai"}},
			{key: "completed", aliases: []string{"completed", "finish", "end"}},
		})
	}
	if len(data.employment) == 0 {
		data.employment = buildRows(allItems, []rowFieldSchema{
			{key: "company", aliases: []string{"company", "perusahaan"}},
			{key: "type", aliases: []string{"type", "jenis", "industry"}},
			{key: "position", aliases: []string{"position", "jabatan", "role"}},
			{key: "from", aliases: []string{"from", "start", "mulai"}},
			{key: "to", aliases: []string{"to", "end", "selesai"}},
			{key: "contact_reference", aliases: []string{"contact reference", "reference", "referensi"}},
		})
	}
	if len(data.achievement) == 0 {
		data.achievement = buildRows(allItems, []rowFieldSchema{
			{key: "company", aliases: []string{"organizer", "issuer", "company"}},
			{key: "type", aliases: []string{"type", "jenis"}},
			{key: "position", aliases: []string{"achievement", "prestasi", "award", "title"}},
			{key: "from", aliases: []string{"from", "start", "mulai"}},
			{key: "to", aliases: []string{"to", "end", "selesai", "year"}},
		})
	}
	if len(data.travel) == 0 {
		data.travel = buildRows(allItems, []rowFieldSchema{
			{key: "country", aliases: []string{"country", "negara"}},
			{key: "purpose", aliases: []string{"purpose", "tujuan"}},
			{key: "start", aliases: []string{"start", "from", "mulai"}},
			{key: "end", aliases: []string{"end", "to", "selesai"}},
		})
	}
	if len(data.nextOfKin) == 0 {
		assignNextOfKinFromItems(data.nextOfKin, allItems)
	}

	ensurePersonalNames(data.personal, guessStudentName(input))
	return data
}

func assignPersonalFromItems(target map[string]string, items []GenerateSectionItem) {
	for _, item := range items {
		assignPersonalValue(target, item.Question, item.Answer)
	}
}

func assignPersonalValue(target map[string]string, question string, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}

	key := matchPersonalField(question)
	if key == "" {
		return
	}

	if key == "full_name" {
		ensurePersonalNames(target, value)
		return
	}
	if target[key] == "" {
		target[key] = value
	}
}

func assignNextOfKinFromItems(target map[string]string, items []GenerateSectionItem) {
	for _, item := range items {
		question := normalizeLabel(item.Question)
		value := strings.TrimSpace(item.Answer)
		if value == "" {
			continue
		}
		switch {
		case containsAny(question, "relationship", "hubungan"):
			setIfEmpty(target, "relationship", value)
		case containsAny(question, "address", "alamat"):
			setIfEmpty(target, "address", value)
		case containsAny(question, "city", "kota"):
			setIfEmpty(target, "city", value)
		case containsAny(question, "province", "state", "provinsi"):
			setIfEmpty(target, "province", value)
		case containsAny(question, "postal", "kode pos"):
			setIfEmpty(target, "postal_code", value)
		case containsAny(question, "country", "negara"):
			setIfEmpty(target, "country", value)
		case containsAny(question, "contact", "phone", "telp", "no hp"):
			setIfEmpty(target, "contact_no", value)
		case containsAny(question, "email"):
			setIfEmpty(target, "email", value)
		case containsAny(question, "father", "mother", "parent", "kin", "wali", "orang tua", "name", "nama"):
			setIfEmpty(target, "name", value)
		}
	}
}

func ensurePersonalNames(target map[string]string, fullName string) {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return
	}
	if target["first_name"] == "" && target["last_name"] == "" {
		parts := strings.Fields(fullName)
		if len(parts) == 1 {
			target["first_name"] = parts[0]
			return
		}
		target["first_name"] = parts[0]
		target["last_name"] = strings.Join(parts[1:], " ")
	}
}

func matchPersonalField(question string) string {
	key := normalizeLabel(question)
	switch {
	case containsAny(key, "full name", "nama lengkap"):
		return "full_name"
	case containsAny(key, "first name", "nama depan"):
		return "first_name"
	case containsAny(key, "last name", "surname", "nama belakang"):
		return "last_name"
	case containsAny(key, "current address", "address", "alamat"):
		return "current_address"
	case containsAny(key, "city", "kota"):
		return "city"
	case containsAny(key, "province", "state", "provinsi"):
		return "province"
	case containsAny(key, "postal code", "zip", "kode pos"):
		return "postal_code"
	case containsAny(key, "country", "negara"):
		return "country"
	case containsAny(key, "date of birth", "birth date", "tanggal lahir"):
		return "date_of_birth"
	case containsAny(key, "place of birth", "birth place", "tempat lahir"):
		return "place_of_birth"
	case containsAny(key, "age", "umur"):
		return "age"
	case containsAny(key, "nationality", "kebangsaan"):
		return "nationality"
	case containsAny(key, "gender", "jenis kelamin"):
		return "gender"
	case containsAny(key, "marital status", "status perkawinan", "status pernikahan"):
		return "marital_status"
	case containsAny(key, "height", "tinggi"):
		return "height"
	case containsAny(key, "weight", "berat"):
		return "weight"
	case containsAny(key, "id number", "ktp", "nik", "identity"):
		return "id_number"
	case containsAny(key, "passport number", "passport"):
		return "passport_number"
	case containsAny(key, "contact number", "phone", "mobile", "whatsapp", "contact", "no hp", "nomor telepon"):
		return "contact_number"
	case containsAny(key, "email"):
		return "email_address"
	default:
		return ""
	}
}

func buildRows(items []GenerateSectionItem, schema []rowFieldSchema) []cvRow {
	rows := make([]cvRow, 1)
	rows[0] = cvRow{}
	currentIndex := 0

	for _, item := range items {
		question := strings.TrimSpace(item.Question)
		answer := strings.TrimSpace(item.Answer)
		if question == "" || answer == "" {
			continue
		}

		column := matchRowField(question, schema)
		if column == "" {
			continue
		}

		if explicitIdx, ok := extractExplicitRowIndex(question); ok {
			for len(rows) <= explicitIdx {
				rows = append(rows, cvRow{})
			}
			assignRowValue(rows[explicitIdx], column, answer)
			currentIndex = explicitIdx
			continue
		}

		if rows[currentIndex][column] != "" {
			rows = append(rows, cvRow{})
			currentIndex++
		}
		assignRowValue(rows[currentIndex], column, answer)
	}

	filtered := make([]cvRow, 0, len(rows))
	for _, row := range rows {
		if len(row) > 0 {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func assignRowValue(row cvRow, key string, value string) {
	if existing := strings.TrimSpace(row[key]); existing != "" && !strings.EqualFold(existing, value) {
		row[key] = existing + "; " + value
		return
	}
	row[key] = value
}

func matchRowField(question string, schema []rowFieldSchema) string {
	normalized := normalizeLabel(question)
	for _, field := range schema {
		for _, alias := range field.aliases {
			if containsAny(normalized, alias) {
				return field.key
			}
		}
	}
	return ""
}

func extractExplicitRowIndex(question string) (int, bool) {
	matches := rowIndexPattern.FindStringSubmatch(strings.ToLower(question))
	if len(matches) < 2 {
		return 0, false
	}

	idxText := strings.TrimSpace(matches[1])
	if idxText == "" {
		return 0, false
	}
	if strings.Contains(question, "202") {
		return 0, false
	}

	switch idxText {
	case "1":
		return 0, true
	case "2":
		return 1, true
	case "3":
		return 2, true
	case "4":
		return 3, true
	case "5":
		return 4, true
	default:
		return 0, false
	}
}

func isPersonalSection(value string) bool {
	return containsAny(value, "personal", "data pribadi", "biodata")
}

func isFormalEducationSection(value string) bool {
	return containsAny(value, "formal education", "pendidikan formal")
}

func isNonformalEducationSection(value string) bool {
	return containsAny(value, "nonformal education", "non formal", "pendidikan nonformal", "pendidikan non formal")
}

func isEmploymentSection(value string) bool {
	return containsAny(value, "employment", "working experience", "pengalaman kerja", "work history")
}

func isAchievementSection(value string) bool {
	return containsAny(value, "achievement", "prestasi", "award")
}

func isTravelSection(value string) bool {
	return containsAny(value, "travel", "perjalanan")
}

func isNextOfKinSection(value string) bool {
	return containsAny(value, "next of kin", "parents", "orang tua", "wali", "keluarga")
}

func guessStudentName(input GenerateDTO) string {
	var firstName string
	var lastName string

	for _, answer := range input.Answers {
		question := strings.ToLower(strings.TrimSpace(answer.Question))
		value := strings.TrimSpace(answer.Value)
		if value == "" {
			continue
		}

		switch {
		case strings.Contains(question, "full name"), strings.Contains(question, "nama lengkap"), question == "name":
			return value
		case strings.Contains(question, "first name"), strings.Contains(question, "nama depan"):
			firstName = value
		case strings.Contains(question, "last name"), strings.Contains(question, "nama belakang"):
			lastName = value
		}
	}

	fullName := strings.TrimSpace(strings.Join([]string{firstName, lastName}, " "))
	if fullName != "" {
		return fullName
	}

	return stringPtrValue(input.StudentID)
}

func normalizeLabel(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.NewReplacer("_", " ", "-", " ", "/", " ", ".", " ", ",", " ", ":", " ").Replace(value)
	return strings.Join(strings.Fields(value), " ")
}

func normalizeKey(value string) string {
	value = normalizeLabel(value)
	value = strings.ReplaceAll(value, " ", "_")
	value = strings.ReplaceAll(value, ".", "")
	return value
}

func containsAny(value string, aliases ...string) bool {
	for _, alias := range aliases {
		if strings.Contains(value, normalizeLabel(alias)) {
			return true
		}
	}
	return false
}

func setIfEmpty(target map[string]string, key string, value string) {
	if strings.TrimSpace(target[key]) == "" {
		target[key] = value
	}
}

func esc(value string) string {
	return html.EscapeString(strings.TrimSpace(value))
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
