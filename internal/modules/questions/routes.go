package questions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/modules/auth"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	protected := rg.Group("")
	protected.Use(auth.AuthMiddleware())

	protected.POST("/question-bases", handler.CreateQuestionBase)
	protected.GET("/question-bases", handler.ListQuestionBases)
	protected.GET("/question-bases/:id", handler.GetQuestionBaseByID)
	protected.PUT("/question-bases/:id", handler.UpdateQuestionBase)
	protected.DELETE("/question-bases/:id", handler.DeleteQuestionBase)

	protected.POST("/questions", handler.CreateQuestion)
	protected.GET("/questions", handler.ListQuestions)
	protected.GET("/questions/:id", handler.GetQuestionByID)
	protected.PUT("/questions/:id", handler.UpdateQuestion)
	protected.DELETE("/questions/:id", handler.DeleteQuestion)

	protected.POST("/question-options", handler.CreateQuestionOption)
	protected.GET("/question-options", handler.ListQuestionOptions)
	protected.GET("/question-options/:id", handler.GetQuestionOptionByID)
	protected.PUT("/question-options/:id", handler.UpdateQuestionOption)
	protected.DELETE("/question-options/:id", handler.DeleteQuestionOption)

	protected.POST("/answer-submissions", handler.CreateAnswerSubmission)
	protected.GET("/answer-submissions", handler.ListAnswerSubmissions)
	protected.GET("/answer-submissions/:id", handler.GetAnswerSubmissionByID)
	protected.PUT("/answer-submissions/:id", handler.UpdateAnswerSubmission)
	protected.DELETE("/answer-submissions/:id", handler.DeleteAnswerSubmission)

	protected.POST("/answer-questions", handler.CreateAnswerQuestion)
	protected.GET("/answer-questions", handler.ListAnswerQuestions)
	protected.GET("/answer-questions/:id", handler.GetAnswerQuestionByID)
	protected.PUT("/answer-questions/:id", handler.UpdateAnswerQuestion)
	protected.DELETE("/answer-questions/:id", handler.DeleteAnswerQuestion)

	protected.POST("/answer-selected-options", handler.CreateAnswerSelectedOption)
	protected.GET("/answer-selected-options", handler.ListAnswerSelectedOptions)
	protected.GET("/answer-selected-options/:answer_id/:option_id", handler.GetAnswerSelectedOption)
	protected.DELETE("/answer-selected-options/:answer_id/:option_id", handler.DeleteAnswerSelectedOption)

	protected.POST("/answer-documents", handler.CreateAnswerDocuments)
	protected.GET("/answer-documents", handler.ListAnswerDocuments)
	protected.GET("/answer-documents/:id", handler.GetAnswerDocumentByID)
	protected.PUT("/answer-documents/:id", handler.UpdateAnswerDocument)
	protected.DELETE("/answer-documents/:id", handler.DeleteAnswerDocument)
}
