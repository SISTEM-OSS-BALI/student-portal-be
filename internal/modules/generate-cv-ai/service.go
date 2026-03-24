package generatecvai

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

type Service struct {
	baseURLs     []string
	defaultModel string
	timeout      time.Duration
	client       *http.Client
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

func (s *Service) Generate(ctx context.Context, input GenerateDTO) (GenerateResponseDTO, error) {
	prompt := strings.TrimSpace(input.Prompt)
	if prompt == "" {
		prompt = buildPromptFromCVPayload(input)
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

	payload := ollamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	}

	requestBody, err := json.Marshal(payload)
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
			fileBase64, fileName, fileErr := generateCVDocument(reqCtx, input, parsed.Response)
			if fileErr != nil {
				return GenerateResponseDTO{}, fileErr
			}

			return GenerateResponseDTO{
				Model:             parsed.Model,
				Response:          parsed.Response,
				Done:              parsed.Done,
				DoneReason:        parsed.DoneReason,
				FileURL:           pickFileURL(input),
				FileBase64:        fileBase64,
				GeneratedFileName: fileName,
				GeneratedMimeType: generatedWordMimeType,
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
		"[ai.service] ollama request url=%s model=%s prompt_len=%d timeout=%s",
		baseURL+"/api/generate",
		model,
		promptLen,
		s.timeout,
	)

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("[ai.service] ollama request error after=%s err=%v", time.Since(start), err)

		if errors.Is(err, context.DeadlineExceeded) || isNetTimeout(err) {
			return ollamaGenerateResponse{}, fmt.Errorf("%w: %v", ErrOllamaTimeout, err)
		}
		return ollamaGenerateResponse{}, fmt.Errorf("%w: %v", ErrOllamaRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, readErr := readAllOrScanner(resp.Body)
		if readErr != nil {
			log.Printf("[ai.service] read error response body after=%s err=%v", time.Since(start), readErr)
		}
		log.Printf("[ai.service] ollama non-2xx body=%s", truncate(string(raw), 2000))
		return ollamaGenerateResponse{}, fmt.Errorf(
			"%w: status=%d body=%s",
			ErrOllamaRequestFailed,
			resp.StatusCode,
			string(raw),
		)
	}

	parsed, raw, err := readGenerateStream(resp.Body)
	if err != nil {
		log.Printf("[ai.service] unmarshal error body=%s err=%v", truncate(string(raw), 2000), err)
		return ollamaGenerateResponse{}, fmt.Errorf("%w: %v", ErrInvalidOllamaPayload, err)
	}

	log.Printf(
		"[ai.service] ollama response status=%d after=%s body_len=%d url=%s",
		resp.StatusCode,
		time.Since(start),
		len(raw),
		baseURL+"/api/generate",
	)

	return parsed, nil
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

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "...(truncated)"
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
			if !parsedAny {
				return ollamaGenerateResponse{}, raw, err
			}
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

func buildPromptFromCVPayload(input GenerateDTO) string {
	var builder strings.Builder

	builder.WriteString("Anda adalah asisten penyusun Statement Latter untuk ke luar negeri yang  profesional.\n")
	builder.WriteString("Tugas Anda adalah menyusun konten CV ringkas, rapi, dan profesional berdasarkan data berikut.\n")
	builder.WriteString("Gunakan bahasa Indonesia yang formal.\n")
	builder.WriteString("Format output: Ringkasan Profil, Pendidikan, Pengalaman, Skill, Sertifikat.\n\n")

	if input.StudentID != nil && strings.TrimSpace(*input.StudentID) != "" {
		builder.WriteString("Student ID: ")
		builder.WriteString(strings.TrimSpace(*input.StudentID))
		builder.WriteString("\n")
	}
	if input.Meta != nil {
		if strings.TrimSpace(input.Meta.CVStatus) != "" {
			builder.WriteString("Status CV: ")
			builder.WriteString(strings.TrimSpace(input.Meta.CVStatus))
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
		builder.WriteString("\nData Jawaban:\n")
		for _, answer := range input.Answers {
			question := strings.TrimSpace(answer.Question)
			value := strings.TrimSpace(answer.Value)
			baseName := strings.TrimSpace(answer.BaseName)
			if question == "" || value == "" {
				continue
			}

			if baseName != "" {
				builder.WriteString("- [")
				builder.WriteString(baseName)
				builder.WriteString("] ")
			} else {
				builder.WriteString("- ")
			}
			builder.WriteString(question)
			builder.WriteString(": ")
			builder.WriteString(value)
			builder.WriteString("\n")
		}
	}

	if len(input.Sections) > 0 {
		builder.WriteString("\nData Per Section:\n")
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

func pickFileURL(input GenerateDTO) string {
	if input.TemplateURL != nil && strings.TrimSpace(*input.TemplateURL) != "" {
		return strings.TrimSpace(*input.TemplateURL)
	}
	if input.TemplatePath != nil && strings.TrimSpace(*input.TemplatePath) != "" {
		path := strings.TrimSpace(*input.TemplatePath)
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			return path
		}
	}
	return ""
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
