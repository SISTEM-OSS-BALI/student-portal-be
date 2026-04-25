package generatesponsorletterai

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

const defaultChecklistSource = "/assets/file/Sponsor Letter Checklist.pdf"
const sponsorChecklistVersion = "Sponsor Letter Checklist"

var sponsorLetterChecklist = []string{
	"Sponsor identity: full name, nationality, and relationship to the student.",
	"Clear statement that the sponsor is willing to financially support the student.",
	"Study destination details: country, institution/campus, and degree or course.",
	"Funding scope: tuition fee, living expenses, accommodation, travel, or other covered costs.",
	"Sponsor financial background: occupation, business, employment, or income source.",
	"Reason the sponsor is supporting the student and confidence in the study plan.",
	"Confirmation that the support will remain available for the duration of study.",
	"Any relevant family context or obligation that strengthens the sponsorship commitment.",
	"Closing statement, date context, and readiness to provide supporting evidence if requested.",
}

func sponsorChecklistItems() []string {
	items := make([]string, len(sponsorLetterChecklist))
	copy(items, sponsorLetterChecklist)
	return items
}

func sponsorChecklistSource() string {
	return defaultChecklistSource
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
		prompt = buildPromptFromSponsorPayload(input)
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

	requestBody, err := json.Marshal(ollamaGenerateRequest{Model: model, Prompt: prompt, Stream: true})
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
			fileBase64, generatedFileName, fileErr := generateSponsorLetterDocument(reqCtx, input, responseText)
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
				ChecklistVersion:  sponsorChecklistVersion,
				ChecklistItems:    sponsorChecklistItems(),
				ChecklistSource:   sponsorChecklistSource(),
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

func (s *Service) generateOnce(ctx context.Context, baseURL string, requestBody []byte, model string, promptLen int) (ollamaGenerateResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/generate", bytes.NewReader(requestBody))
	if err != nil {
		return ollamaGenerateResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	log.Printf("[sponsor-letter.service] ollama request url=%s model=%s prompt_len=%d timeout=%s", baseURL+"/api/generate", model, promptLen, s.timeout)

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("[sponsor-letter.service] request error after=%s err=%v", time.Since(start), err)
		if errors.Is(err, context.DeadlineExceeded) || isNetTimeout(err) {
			return ollamaGenerateResponse{}, fmt.Errorf("%w: %v", ErrOllamaTimeout, err)
		}
		return ollamaGenerateResponse{}, fmt.Errorf("%w: %v", ErrOllamaRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, readErr := readAllOrScanner(resp.Body)
		if readErr != nil {
			log.Printf("[sponsor-letter.service] read error response body after=%s err=%v", time.Since(start), readErr)
		}
		return ollamaGenerateResponse{}, fmt.Errorf("%w: status=%d body=%s", ErrOllamaRequestFailed, resp.StatusCode, string(raw))
	}

	parsed, raw, err := readGenerateStream(resp.Body)
	if err != nil {
		log.Printf("[sponsor-letter.service] unmarshal error body=%s err=%v", truncate(string(raw), 2000), err)
		return ollamaGenerateResponse{}, fmt.Errorf("%w: %v", ErrInvalidOllamaPayload, err)
	}

	return parsed, nil
}

func buildPromptFromSponsorPayload(input GenerateDTO) string {
	var builder strings.Builder
	destinationCountry := strings.TrimSpace(stringPtrValue(input.StudentCountry))
	campusName := strings.TrimSpace(stringPtrValue(input.CampusName))
	degree := strings.TrimSpace(stringPtrValue(input.Degree))

	builder.WriteString("You are a senior admissions writing assistant.\n")
	builder.WriteString("Write a Sponsor Letter in professional English based strictly on the provided student and sponsor data.\n")
	builder.WriteString("The sponsor letter should be written from the sponsor perspective, confirming financial support for the student.\n")
	builder.WriteString("Do not invent facts. If information is missing, do not fabricate it.\n")
	builder.WriteString("After the sponsor letter, add a section titled 'Missing Information' with bullet points for unsupported sponsor-letter requirements.\n")
	builder.WriteString("Use any explicit destination country, campus, or degree metadata as source of truth.\n")
	builder.WriteString("Use the checklist below as mandatory coverage guidance:\n")
	builder.WriteString("Checklist source reference: ")
	builder.WriteString(sponsorChecklistSource())
	builder.WriteString("\n")
	for idx, item := range sponsorLetterChecklist {
		builder.WriteString(fmt.Sprintf("%d. %s\n", idx+1, item))
	}

	builder.WriteString("\nRequired output structure:\n")
	builder.WriteString("- Opening sponsor declaration.\n")
	builder.WriteString("- Sponsor relationship and financial background.\n")
	builder.WriteString("- Commitment to cover study and living expenses.\n")
	builder.WriteString("- Confidence in the student's study plan and future.\n")
	builder.WriteString("- Closing commitment and readiness to provide evidence.\n")
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
				continue
			}
			builder.WriteString("\nSection: ")
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

func detectMissingIndicators(content string) []string {
	lower := strings.ToLower(content)
	var missing []string

	checks := map[string][]string{
		"sponsor relationship":    {"relationship", "father", "mother", "guardian", "sponsor"},
		"financial support scope": {"tuition", "living expenses", "accommodation", "financial support"},
		"funding source":          {"business", "employment", "income", "financial"},
		"duration of sponsorship": {"duration", "entire study", "until completion", "throughout"},
		"supporting evidence":     {"evidence", "bank statement", "supporting document", "if requested"},
	}

	for label, keywords := range checks {
		matched := false
		for _, keyword := range keywords {
			if strings.Contains(lower, keyword) {
				matched = true
				break
			}
		}
		if !matched {
			missing = append(missing, label)
		}
	}

	return missing
}

func resolveBaseURLs() []string {
	primary := strings.TrimSpace(getEnv("OLLAMA_BASE_URL", "http://127.0.0.1:11434"))
	rawList := strings.TrimSpace(getEnv("OLLAMA_BASE_URLS", ""))

	seen := map[string]struct{}{}
	var out []string
	add := func(v string) {
		v = strings.TrimRight(strings.TrimSpace(v), "/")
		if v == "" {
			return
		}
		if _, ok := seen[v]; ok {
			return
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}

	add(primary)
	if rawList != "" {
		for _, item := range strings.Split(rawList, ",") {
			add(item)
		}
	}

	if len(out) == 0 {
		add("http://127.0.0.1:11434")
	}

	return out
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	v := strings.TrimSpace(getEnv(key, ""))
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func readGenerateStream(r io.Reader) (ollamaGenerateResponse, []byte, error) {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var raw bytes.Buffer
	var final ollamaGenerateResponse
	var combined strings.Builder
	seen := false

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		raw.Write(line)
		raw.WriteByte('\n')

		var part ollamaGenerateResponse
		if err := json.Unmarshal(line, &part); err != nil {
			return ollamaGenerateResponse{}, raw.Bytes(), err
		}
		seen = true
		if part.Model != "" {
			final.Model = part.Model
		}
		if part.Response != "" {
			combined.WriteString(part.Response)
		}
		if part.Done {
			final.Done = true
			final.DoneReason = part.DoneReason
		}
	}
	if err := scanner.Err(); err != nil {
		return ollamaGenerateResponse{}, raw.Bytes(), err
	}
	if !seen {
		all, err := io.ReadAll(r)
		if err == nil && len(all) > 0 {
			raw.Write(all)
		}
		return ollamaGenerateResponse{}, raw.Bytes(), io.EOF
	}
	final.Response = combined.String()
	return final, raw.Bytes(), nil
}

func readAllOrScanner(r io.Reader) ([]byte, error) {
	b, err := io.ReadAll(r)
	if err == nil {
		return b, nil
	}
	var raw bytes.Buffer
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		raw.Write(scanner.Bytes())
		raw.WriteByte('\n')
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return raw.Bytes(), scanErr
	}
	return raw.Bytes(), err
}

func isNetTimeout(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func shouldRetryWithNextBase(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrOllamaTimeout) {
		return true
	}
	msg := strings.ToLower(err.Error())
	for _, part := range []string{"connection refused", "no such host", "i/o timeout", "bad gateway", "502", "503", "504"} {
		if strings.Contains(msg, part) {
			return true
		}
	}
	return false
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max] + "..."
}
