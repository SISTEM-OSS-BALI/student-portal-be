package generatecvai

type GenerateDTO struct {
	Prompt         string               `json:"prompt"`
	Model          *string              `json:"model"`
	StudentID      *string              `json:"student_id"`
	StudentName    *string              `json:"student_name"`
	TemplatePath   *string              `json:"template_path"`
	TemplateURL    *string              `json:"template_url"`
	TemplateBase64 *string              `json:"template_base64"`
	Answers        []GenerateAnswerDTO  `json:"answers"`
	Sections       []GenerateSectionDTO `json:"sections"`
	Meta           *GenerateMetaDTO     `json:"meta"`
}

type GenerateAnswerDTO struct {
	AnswerID          string   `json:"answer_id"`
	QuestionID        string   `json:"question_id"`
	Question          string   `json:"question"`
	InputType         string   `json:"input_type"`
	BaseID            string   `json:"base_id"`
	BaseName          string   `json:"base_name"`
	Value             string   `json:"value"`
	SelectedOptionIDs []string `json:"selected_option_ids"`
}

type GenerateSectionDTO struct {
	Key   string                `json:"key"`
	Label string                `json:"label"`
	Items []GenerateSectionItem `json:"items"`
}

type GenerateSectionItem struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type GenerateMetaDTO struct {
	CVStatus    string `json:"cv_status"`
	SubmittedAt string `json:"submitted_at"`
	AdmissionAt string `json:"admission_at"`
}

type GenerateResponseDTO struct {
	Model             string `json:"model"`
	Response          string `json:"response"`
	Done              bool   `json:"done"`
	DoneReason        string `json:"done_reason,omitempty"`
	FileURL           string `json:"file_url,omitempty"`
	FileBase64        string `json:"file_base64,omitempty"`
	GeneratedFileName string `json:"generated_file_name,omitempty"`
	GeneratedMimeType string `json:"generated_mime_type,omitempty"`
}
