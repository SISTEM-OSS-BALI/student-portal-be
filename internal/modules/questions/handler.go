package questions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/username/gin-gorm-api/internal/modules/auth"
	"github.com/username/gin-gorm-api/internal/httpx"
	"github.com/username/gin-gorm-api/internal/schema"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateQuestionBase(c *gin.Context) {
	var input QuestionBaseCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	base, err := h.service.CreateQuestionBase(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewQuestionBaseResponseDTO(base))
}

func (h *Handler) ListQuestionBases(c *gin.Context) {
	var filter QuestionBaseFilter
	if value := c.Query("type"); value != "" {
		filter.TypeCountry = &value
	}
	if value := c.Query("country_id"); value != "" {
		filter.CountryID = &value
	}
	if value := c.Query("active"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			httpx.RespondError(c, http.StatusBadRequest, "validation_error", "invalid active flag", nil)
			return
		}
		filter.Active = &parsed
	}

	bases, err := h.service.ListQuestionBases(filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewQuestionBaseResponseListDTO(bases))
}

func (h *Handler) GetQuestionBaseByID(c *gin.Context) {
	base, err := h.service.GetQuestionBaseByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "question base not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewQuestionBaseResponseDTO(base))
}

func (h *Handler) UpdateQuestionBase(c *gin.Context) {
	var input QuestionBaseUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	base, err := h.service.UpdateQuestionBase(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewQuestionBaseResponseDTO(base))
}

func (h *Handler) DeleteQuestionBase(c *gin.Context) {
	if err := h.service.DeleteQuestionBase(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateQuestion(c *gin.Context) {
	var input QuestionCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	question, err := h.service.CreateQuestion(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewQuestionResponseDTO(question))
}

func (h *Handler) ListQuestions(c *gin.Context) {
	var filter QuestionFilter
	if value := c.Query("base_id"); value != "" {
		filter.BaseID = &value
	}
	if value := c.Query("active"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			httpx.RespondError(c, http.StatusBadRequest, "validation_error", "invalid active flag", nil)
			return
		}
		filter.Active = &parsed
	}

	questions, err := h.service.ListQuestions(filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewQuestionResponseListDTO(questions))
}

func (h *Handler) GetQuestionByID(c *gin.Context) {
	question, err := h.service.GetQuestionByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "question not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewQuestionResponseDTO(question))
}

func (h *Handler) UpdateQuestion(c *gin.Context) {
	var input QuestionUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	question, err := h.service.UpdateQuestion(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewQuestionResponseDTO(question))
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	if err := h.service.DeleteQuestion(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateQuestionOption(c *gin.Context) {
	var input QuestionOptionCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	option, err := h.service.CreateQuestionOption(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewQuestionOptionResponseDTO(option))
}

func (h *Handler) ListQuestionOptions(c *gin.Context) {
	var filter QuestionOptionFilter
	if value := c.Query("question_id"); value != "" {
		filter.QuestionID = &value
	}
	if value := c.Query("active"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			httpx.RespondError(c, http.StatusBadRequest, "validation_error", "invalid active flag", nil)
			return
		}
		filter.Active = &parsed
	}

	options, err := h.service.ListQuestionOptions(filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewQuestionOptionResponseListDTO(options))
}

func (h *Handler) GetQuestionOptionByID(c *gin.Context) {
	option, err := h.service.GetQuestionOptionByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "question option not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewQuestionOptionResponseDTO(option))
}

func (h *Handler) UpdateQuestionOption(c *gin.Context) {
	var input QuestionOptionUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	option, err := h.service.UpdateQuestionOption(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewQuestionOptionResponseDTO(option))
}

func (h *Handler) DeleteQuestionOption(c *gin.Context) {
	if err := h.service.DeleteQuestionOption(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateAnswerSubmission(c *gin.Context) {
	var input AnswerSubmissionCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	submission, err := h.service.CreateAnswerSubmission(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewAnswerSubmissionResponseDTO(submission))
}

func (h *Handler) ListAnswerSubmissions(c *gin.Context) {
	var filter AnswerSubmissionFilter
	if value := c.Query("base_id"); value != "" {
		filter.BaseID = &value
	}
	if value := c.Query("student_id"); value != "" {
		filter.StudentID = &value
	}
	if value := c.Query("status"); value != "" {
		filter.Status = &value
	}

	submissions, err := h.service.ListAnswerSubmissions(filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewAnswerSubmissionResponseListDTO(submissions))
}

func (h *Handler) GetAnswerSubmissionByID(c *gin.Context) {
	submission, err := h.service.GetAnswerSubmissionByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "answer submission not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewAnswerSubmissionResponseDTO(submission))
}

func (h *Handler) UpdateAnswerSubmission(c *gin.Context) {
	var input AnswerSubmissionUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	submission, err := h.service.UpdateAnswerSubmission(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewAnswerSubmissionResponseDTO(submission))
}

func (h *Handler) DeleteAnswerSubmission(c *gin.Context) {
	if err := h.service.DeleteAnswerSubmission(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateAnswerQuestion(c *gin.Context) {
	var input AnswerQuestionCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	answer, err := h.service.CreateAnswerQuestion(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	selectedIDs, err := h.loadSelectedOptionIDs(answer.ID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusCreated, NewAnswerQuestionResponseDTO(answer, selectedIDs))
}

func (h *Handler) ListAnswerQuestions(c *gin.Context) {
	var filter AnswerQuestionFilter
	if value := c.Query("submission_id"); value != "" {
		filter.SubmissionID = &value
	}
	if value := c.Query("question_id"); value != "" {
		filter.QuestionID = &value
	}
	if value := c.Query("student_id"); value != "" {
		filter.UserID = &value
	}

	answers, err := h.service.ListAnswerQuestions(filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}

	selectedMap, err := h.loadSelectedOptionIDsByAnswers(answers)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewAnswerQuestionResponseListDTO(answers, selectedMap))
}

func (h *Handler) GetAnswerQuestionByID(c *gin.Context) {
	answer, err := h.service.GetAnswerQuestionByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "answer not found", nil)
		return
	}
	selectedIDs, err := h.loadSelectedOptionIDs(answer.ID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewAnswerQuestionResponseDTO(answer, selectedIDs))
}

func (h *Handler) UpdateAnswerQuestion(c *gin.Context) {
	var input AnswerQuestionUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	answer, err := h.service.UpdateAnswerQuestion(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	selectedIDs, err := h.loadSelectedOptionIDs(answer.ID)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewAnswerQuestionResponseDTO(answer, selectedIDs))
}

func (h *Handler) DeleteAnswerQuestion(c *gin.Context) {
	if err := h.service.DeleteAnswerQuestion(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateAnswerSelectedOption(c *gin.Context) {
	var input AnswerSelectedOptionCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	selected, err := h.service.CreateAnswerSelectedOption(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusCreated, NewAnswerSelectedOptionResponseDTO(selected))
}

func (h *Handler) ListAnswerSelectedOptions(c *gin.Context) {
	var filter AnswerSelectedOptionFilter
	if value := c.Query("answer_id"); value != "" {
		filter.AnswerID = &value
	}
	if value := c.Query("option_id"); value != "" {
		filter.OptionID = &value
	}

	selected, err := h.service.ListAnswerSelectedOptions(filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewAnswerSelectedOptionResponseListDTO(selected))
}

func (h *Handler) GetAnswerSelectedOption(c *gin.Context) {
	answerID := c.Param("answer_id")
	optionID := c.Param("option_id")
	selected, err := h.service.GetAnswerSelectedOption(answerID, optionID)
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "answer selected option not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewAnswerSelectedOptionResponseDTO(selected))
}

func (h *Handler) DeleteAnswerSelectedOption(c *gin.Context) {
	answerID := c.Param("answer_id")
	optionID := c.Param("option_id")
	if err := h.service.DeleteAnswerSelectedOption(answerID, optionID); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateAnswerDocuments(c *gin.Context) {
	var input []AnswerDocumentCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}
	if len(input) == 0 {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", "documents is required", nil)
		return
	}

	claims, ok := getAuthClaims(c)
	if !ok {
		return
	}

	for i := range input {
		if input[i].StudentID == nil || *input[i].StudentID == "" {
			input[i].StudentID = &claims.UserID
		}
	}

	docs, err := h.service.CreateAnswerDocuments(input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "create_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusCreated, NewAnswerDocumentResponseListDTO(docs))
}

func (h *Handler) ListAnswerDocuments(c *gin.Context) {
	var filter AnswerDocumentFilter
	if value := c.Query("submission_id"); value != "" {
		filter.SubmissionID = &value
	}
	if value := c.Query("student_id"); value != "" {
		filter.StudentID = &value
	}
	if value := c.Query("document_id"); value != "" {
		filter.DocumentID = &value
	}

	docs, err := h.service.ListAnswerDocuments(filter)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "list_failed", err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, NewAnswerDocumentResponseListDTO(docs))
}

func (h *Handler) GetAnswerDocumentByID(c *gin.Context) {
	doc, err := h.service.GetAnswerDocumentByID(c.Param("id"))
	if err != nil {
		httpx.RespondError(c, http.StatusNotFound, "not_found", "answer document not found", nil)
		return
	}
	c.JSON(http.StatusOK, NewAnswerDocumentResponseDTO(doc))
}

func (h *Handler) UpdateAnswerDocument(c *gin.Context) {
	var input AnswerDocumentUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	doc, err := h.service.UpdateAnswerDocument(c.Param("id"), input)
	if err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "update_failed", err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, NewAnswerDocumentResponseDTO(doc))
}

func (h *Handler) DeleteAnswerDocument(c *gin.Context) {
	if err := h.service.DeleteAnswerDocument(c.Param("id")); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error(), nil)
		return
	}
	c.Status(http.StatusNoContent)
}

func getAuthClaims(c *gin.Context) (*auth.Claims, bool) {
	claimsAny, ok := c.Get("auth")
	if !ok {
		httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "missing auth claims", nil)
		return nil, false
	}
	claims, ok := claimsAny.(*auth.Claims)
	if !ok {
		httpx.RespondError(c, http.StatusUnauthorized, "unauthorized", "invalid auth claims", nil)
		return nil, false
	}
	return claims, true
}

func (h *Handler) loadSelectedOptionIDs(answerID string) ([]string, error) {
	filter := AnswerSelectedOptionFilter{AnswerID: &answerID}
	selected, err := h.service.ListAnswerSelectedOptions(filter)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(selected))
	for _, item := range selected {
		ids = append(ids, item.OptionID)
	}
	return ids, nil
}

func (h *Handler) loadSelectedOptionIDsByAnswers(answers []schema.AnswerQuestion) (map[string][]string, error) {
	answerIDs := make([]string, 0, len(answers))
	for _, answer := range answers {
		answerIDs = append(answerIDs, answer.ID)
	}

	result := make(map[string][]string, len(answers))
	if len(answerIDs) == 0 {
		return result, nil
	}

	selected, err := h.service.ListAnswerSelectedOptions(AnswerSelectedOptionFilter{AnswerIDs: answerIDs})
	if err != nil {
		return nil, err
	}
	for _, item := range selected {
		result[item.AnswerID] = append(result[item.AnswerID], item.OptionID)
	}
	return result, nil
}
