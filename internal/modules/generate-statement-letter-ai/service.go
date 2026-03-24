package generatestatementletterai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrPromptRequired       = errors.New("prompt is required (or provide answers/sections payload)")
	ErrModelRequired        = errors.New("model is required")
	ErrOllamaRequestFailed  = errors.New("ollama request failed")
	ErrInvalidOllamaPayload = errors.New("invalid ollama response")
	ErrOllamaTimeout        = errors.New("ollama request timeout")
)

const defaultChecklistSource = "/assets/file/Kerangka GS 2026.pdf"

var gsChecklist2026 = []string{
	"Applicant identity: name, date of birth, nationality, and application ID.",
	"Purpose statement: clearly state the main goal to study in Australia and return home after completion.",
	"Course and institution: mention the chosen course/major and institution.",
	"Study destination: mention the city and country of study.",
	"Education background: latest school/university, major, and graduation years.",
	"Work or business experience: current or previous employment/business and its context.",
	"Organizational or event involvement, if any, and relevance to the chosen course.",
	"Home-country ties: family, property, assets, commitments, or relationship context that support return intention.",
	"Course consistency: explain why the chosen course matches prior education and work background.",
	"Entry requirements and evidence: academic results, IELTS or other prerequisites if available.",
	"Availability in home country: explain similar study options locally and why Australia is still chosen.",
	"Course detail: modules/subjects and skills to be gained, linked to future plan.",
}

type Service struct {
	baseURLs     []string
	defaultModel string
	timeout      time.Duration
	client       *http.Client
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaGenerateResponse struct {
	Model      string `json:"model"`
	Response   string `json:"response"`
	Done       bool   `json:"done"`
	DoneReason string `json:"done_reason"`
}

func NewService() *Service {
	baseURLs := resolveBaseURLs()
	defaultModel := strings.TrimSpace(getEnv("OLLAMA_MODEL", "qwen2.5:3b"))
	timeoutSec := getEnvAsInt("OLLAMA_TIMEOUT_SEC", 120)

	return &Service{
		baseURLs:     baseURLs,
		defaultModel: defaultModel,
		timeout:      time.Duration(timeoutSec) * time.Second,
		client: &http.Client{
			Timeout: 0,
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   20,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}
}

func (s *Service) Generate(ctx context.Context, input GenerateDTO) (GenerateResponseDTO, error) {
	prompt := strings.TrimSpace(input.Prompt)
	if prompt == "" {
		prompt = buildPromptFromStatementPayload(input)
	}
	if prompt == "" {
		return GenerateResponseDTO{}, ErrPromptRequired
	}

	model := s.defaultModel
	if input.Model != nil && strings.TrimSpace(*input.Model) != "" {
		model = strings.TrimSpace(*input.Model)
	}
	if model == "" {
		return GenerateResponseDTO{}, ErrModelRequired
	}

	requestBody, err := json.Marshal(ollamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	})
	if err != nil {
		return GenerateResponseDTO{}, err
	}

	reqCtx := ctx
	var cancel context.CancelFunc
	if _, hasDeadline := ctx.Deadline(); !hasDeadline && s.timeout > 0 {
		reqCtx, cancel = context.WithTimeout(ctx, s.timeout)
		defer cancel()
	}

	var lastErr error
	for _, baseURL := range s.baseURLs {
		parsed, err := s.generateOnce(reqCtx, baseURL, requestBody, model, len(prompt))
		if err == nil {
			responseText := strings.TrimSpace(parsed.Response)
			fileBase64, generatedFileName, fileErr := generateStatementLetterDocument(reqCtx, input, responseText)
			if fileErr != nil {
				return GenerateResponseDTO{}, fileErr
			}

			return GenerateResponseDTO{
				Model:             parsed.Model,
				Response:          responseText,
				Done:              parsed.Done,
				DoneReason:        parsed.DoneReason,
				FileBase64:        fileBase64,
				GeneratedFileName: generatedFileName,
				GeneratedMimeType: generatedWordMimeType,
				ChecklistVersion:  "GS 2026",
				ChecklistItems:    gsChecklist2026,
				ChecklistSource:   resolveChecklistSource(input),
				MissingIndicators: detectMissingIndicators(parsed.Response),
			}, nil
		}

		lastErr = err
		if !shouldRetryWithNextBase(err) {
			return GenerateResponseDTO{}, err
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("%w: no ollama base URL configured", ErrOllamaRequestFailed)
	}

	return GenerateResponseDTO{}, lastErr
}

func (s *Service) generateOnce(
	ctx context.Context,
	baseURL string,
	requestBody []byte,
	model string,
	promptLen int,
) (ollamaGenerateResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/api/generate",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return ollamaGenerateResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	log.Printf(
		"[statement-letter.service] ollama request url=%s model=%s prompt_len=%d timeout=%s",
		baseURL+"/api/generate",
		model,
		promptLen,
		s.timeout,
	)

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("[statement-letter.service] request error after=%s err=%v", time.Since(start), err)

		if errors.Is(err, context.DeadlineExceeded) || isNetTimeout(err) {
			return ollamaGenerateResponse{}, fmt.Errorf("%w: %v", ErrOllamaTimeout, err)
		}
		return ollamaGenerateResponse{}, fmt.Errorf("%w: %v", ErrOllamaRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, readErr := readAllOrScanner(resp.Body)
		if readErr != nil {
			log.Printf("[statement-letter.service] read error response body after=%s err=%v", time.Since(start), readErr)
		}
		return ollamaGenerateResponse{}, fmt.Errorf(
			"%w: status=%d body=%s",
			ErrOllamaRequestFailed,
			resp.StatusCode,
			string(raw),
		)
	}

	parsed, raw, err := readGenerateStream(resp.Body)
	if err != nil {
		log.Printf("[statement-letter.service] unmarshal error body=%s err=%v", truncate(string(raw), 2000), err)
		return ollamaGenerateResponse{}, fmt.Errorf("%w: %v", ErrInvalidOllamaPayload, err)
	}

	return parsed, nil
}

func buildPromptFromStatementPayload(input GenerateDTO) string {
	var builder strings.Builder
	destinationCountry := strings.TrimSpace(stringPtrValue(input.StudentCountry))
	campusName := strings.TrimSpace(stringPtrValue(input.CampusName))
	degree := strings.TrimSpace(stringPtrValue(input.Degree))

	builder.WriteString("You are a senior admissions writing assistant.\n")
	builder.WriteString("Write a Genuine Student (GS) / Statement Letter in professional English based strictly on the provided student data.\n")
	builder.WriteString("Do not invent facts. If information is missing, do not fabricate it.\n")
	builder.WriteString("After the statement letter, add a section titled 'Missing Information' with bullet points for any GS checklist items that are still unsupported by the provided data.\n")
	builder.WriteString("If explicit student profile metadata is provided for destination country, campus, or degree, use that metadata as the source of truth.\n")
	builder.WriteString("Use the GS 2026 checklist below as mandatory coverage guidance:\n")
	builder.WriteString("Checklist source reference: ")
	builder.WriteString(resolveChecklistSource(input))
	builder.WriteString("\n")
	for idx, item := range gsChecklist2026 {
		builder.WriteString(fmt.Sprintf("%d. %s\n", idx+1, item))
	}
	builder.WriteString("\nRequired output structure:\n")
	builder.WriteString("- Opening paragraph introducing the applicant and study purpose.\n")
	builder.WriteString("- Academic and career background in home country.\n")
	if destinationCountry != "" {
		builder.WriteString("- Reason for selecting the course, institution, city, and ")
		builder.WriteString(destinationCountry)
		builder.WriteString(".\n")
	} else {
		builder.WriteString("- Reason for selecting the course, institution, city, and destination country.\n")
	}
	builder.WriteString("- Evidence of home-country ties and intention to return.\n")
	builder.WriteString("- Future study and career plan after graduation.\n")
	builder.WriteString("- Missing Information section.\n")

	if input.StudentName != nil && strings.TrimSpace(*input.StudentName) != "" {
		builder.WriteString("\nStudent Name: ")
		builder.WriteString(strings.TrimSpace(*input.StudentName))
		builder.WriteString("\n")
	}
	if input.StudentID != nil && strings.TrimSpace(*input.StudentID) != "" {
		builder.WriteString("Student ID: ")
		builder.WriteString(strings.TrimSpace(*input.StudentID))
		builder.WriteString("\n")
	}
	if destinationCountry != "" {
		builder.WriteString("Destination Country: ")
		builder.WriteString(destinationCountry)
		builder.WriteString("\n")
	}
	if campusName != "" {
		builder.WriteString("Target Campus / Major: ")
		builder.WriteString(campusName)
		builder.WriteString("\n")
	}
	if degree != "" {
		builder.WriteString("Target Degree: ")
		builder.WriteString(degree)
		builder.WriteString("\n")
	}
	if input.Meta != nil {
		if strings.TrimSpace(input.Meta.LetterStatus) != "" {
			builder.WriteString("Letter Status: ")
			builder.WriteString(strings.TrimSpace(input.Meta.LetterStatus))
			builder.WriteString("\n")
		}
		if strings.TrimSpace(input.Meta.SubmittedAt) != "" {
			builder.WriteString("Submitted At: ")
			builder.WriteString(strings.TrimSpace(input.Meta.SubmittedAt))
			builder.WriteString("\n")
		}
		if strings.TrimSpace(input.Meta.AdmissionAt) != "" {
			builder.WriteString("Admission At: ")
			builder.WriteString(strings.TrimSpace(input.Meta.AdmissionAt))
			builder.WriteString("\n")
		}
	}

	if len(input.Answers) > 0 {
		builder.WriteString("\nCollected Answers:\n")
		for _, answer := range input.Answers {
			question := strings.TrimSpace(answer.Question)
			value := strings.TrimSpace(answer.Value)
			baseName := strings.TrimSpace(answer.BaseName)
			if question == "" || value == "" {
				continue
			}
			builder.WriteString("- ")
			if baseName != "" {
				builder.WriteString("[")
				builder.WriteString(baseName)
				builder.WriteString("] ")
			}
			builder.WriteString(question)
			builder.WriteString(": ")
			builder.WriteString(value)
			builder.WriteString("\n")
		}
	}

	if len(input.Sections) > 0 {
		builder.WriteString("\nGrouped Sections:\n")
		for _, section := range input.Sections {
			label := strings.TrimSpace(section.Label)
			if label == "" {
				label = strings.TrimSpace(section.Key)
			}
			if label == "" {
				label = "Section"
			}
			builder.WriteString("## ")
			builder.WriteString(label)
			builder.WriteString("\n")
			for _, item := range section.Items {
				question := strings.TrimSpace(item.Question)
				answer := strings.TrimSpace(item.Answer)
				if question == "" || answer == "" {
					continue
				}
				builder.WriteString("- ")
				builder.WriteString(question)
				builder.WriteString(": ")
				builder.WriteString(answer)
				builder.WriteString("\n")
			}
		}
	}

	return strings.TrimSpace(builder.String())
}

func detectMissingIndicators(response string) []string {
	normalized := strings.ToLower(response)
	indicators := make([]string, 0)
	for _, marker := range []string{
		"missing information",
		"not provided",
		"not available",
		"belum tersedia",
		"tidak tersedia",
		"tidak disebutkan",
	} {
		if strings.Contains(normalized, marker) {
			indicators = append(indicators, marker)
		}
	}
	return indicators
}

func resolveChecklistSource(input GenerateDTO) string {
	if input.ChecklistURL != nil && strings.TrimSpace(*input.ChecklistURL) != "" {
		return strings.TrimSpace(*input.ChecklistURL)
	}
	if input.ChecklistPath != nil && strings.TrimSpace(*input.ChecklistPath) != "" {
		return strings.TrimSpace(*input.ChecklistPath)
	}
	return defaultChecklistSource
}

func shouldRetryWithNextBase(err error) bool {
	if errors.Is(err, ErrOllamaTimeout) {
		return true
	}
	if errors.Is(err, ErrOllamaRequestFailed) {
		message := err.Error()
		return strings.Contains(message, "status=5")
	}
	return false
}

func resolveBaseURLs() []string {
	rawList := strings.TrimSpace(os.Getenv("OLLAMA_BASE_URLS"))
	if rawList == "" {
		rawList = getEnv("OLLAMA_BASE_URL", "https://llm.onestepsolutionbali.com")
	}

	seen := make(map[string]struct{})
	baseURLs := make([]string, 0)
	for _, item := range strings.Split(rawList, ",") {
		baseURL := strings.TrimRight(strings.TrimSpace(item), "/")
		if baseURL == "" {
			continue
		}
		if _, exists := seen[baseURL]; exists {
			continue
		}
		seen[baseURL] = struct{}{}
		baseURLs = append(baseURLs, baseURL)
	}
	if len(baseURLs) == 0 {
		return []string{"https://llm.onestepsolutionbali.com"}
	}
	return baseURLs
}

func isNetTimeout(err error) bool {
	var nerr net.Error
	return errors.As(err, &nerr) && nerr.Timeout()
}

func readGenerateStream(body io.Reader) (ollamaGenerateResponse, []byte, error) {
	raw, err := readAllOrScanner(body)
	if err != nil {
		return ollamaGenerateResponse{}, raw, err
	}

	lines := bytes.Split(raw, []byte("\n"))
	var aggregated ollamaGenerateResponse
	var responseBuilder strings.Builder
	var parsedAny bool

	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		var chunk ollamaGenerateResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			return ollamaGenerateResponse{}, raw, err
		}

		parsedAny = true
		if chunk.Model != "" {
			aggregated.Model = chunk.Model
		}
		if chunk.Response != "" {
			responseBuilder.WriteString(chunk.Response)
		}
		if chunk.Done {
			aggregated.Done = true
		}
		if chunk.DoneReason != "" {
			aggregated.DoneReason = chunk.DoneReason
		}
	}

	if !parsedAny {
		return ollamaGenerateResponse{}, raw, errors.New("empty ollama response body")
	}

	aggregated.Response = responseBuilder.String()
	return aggregated, raw, nil
}

func readAllOrScanner(body io.Reader) ([]byte, error) {
	var raw bytes.Buffer
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		raw.Write(scanner.Bytes())
		raw.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return raw.Bytes(), err
	}
	return bytes.TrimSpace(raw.Bytes()), nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "...(truncated)"
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
