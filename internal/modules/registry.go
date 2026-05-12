package modules

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/modules/answer-approvals"
	"github.com/username/gin-gorm-api/internal/modules/answer-document-approvals"
	"github.com/username/gin-gorm-api/internal/modules/auth"
	"github.com/username/gin-gorm-api/internal/modules/chat"
	"github.com/username/gin-gorm-api/internal/modules/child-steps-management"
	"github.com/username/gin-gorm-api/internal/modules/country-management"
	"github.com/username/gin-gorm-api/internal/modules/country-steps-management"
	"github.com/username/gin-gorm-api/internal/modules/document-translations"
	"github.com/username/gin-gorm-api/internal/modules/documents-managament"
	generatecvai "github.com/username/gin-gorm-api/internal/modules/generate-cv-ai"
	generatesponsorletterai "github.com/username/gin-gorm-api/internal/modules/generate-sponsor-letter-ai"
	generatestatementletterai "github.com/username/gin-gorm-api/internal/modules/generate-statement-letter-ai"
	informationcountrymanagement "github.com/username/gin-gorm-api/internal/modules/information-country-management"
	"github.com/username/gin-gorm-api/internal/modules/notes-student"
	"github.com/username/gin-gorm-api/internal/modules/promo"
	"github.com/username/gin-gorm-api/internal/modules/questions"
	sponsorletteraiapprovals "github.com/username/gin-gorm-api/internal/modules/sponsor-letter-ai-approvals"
	"github.com/username/gin-gorm-api/internal/modules/stages-management"
	statementletteraiapprovals "github.com/username/gin-gorm-api/internal/modules/statement-letter-ai-approvals"
	"github.com/username/gin-gorm-api/internal/modules/steps-management"
	ticketmessage "github.com/username/gin-gorm-api/internal/modules/ticket-message"
	"github.com/username/gin-gorm-api/internal/modules/user"
	visatype "github.com/username/gin-gorm-api/internal/modules/visa-type-management"
)

func RegisterAll(rg *gin.RouterGroup, db *gorm.DB) {
	auth.RegisterRoutes(rg, db)
	user.RegisterRoutes(rg, db)
	protected := rg.Group("")
	protected.Use(auth.AuthMiddleware())
	user.RegisterProtectedRoutes(protected, db)
	country.RegisterRoutes(rg, db)
	documents.RegisterRoutes(rg, db)
	chat.RegisterRoutes(rg, db)
	generatecvai.RegisterRoutes(rg, db)
	generatesponsorletterai.RegisterRoutes(rg, db)
	sponsorletteraiapprovals.RegisterRoutes(rg, db)
	generatestatementletterai.RegisterRoutes(rg, db)
	statementletteraiapprovals.RegisterRoutes(rg, db)
	notesstudent.RegisterRoutes(rg, db)
	questions.RegisterRoutes(rg, db)
	answerapprovals.RegisterRoutes(rg, db)
	answerdocumentapprovals.RegisterRoutes(rg, db)
	documenttranslations.RegisterRoutes(rg, db)
	stages.RegisterRoutes(rg, db)
	steps.RegisterRoutes(rg, db)
	childsteps.RegisterRoutes(rg, db)
	countrysteps.RegisterRoutes(rg, db)
	ticketmessage.RegisterRoutes(rg, db)
	informationcountrymanagement.RegisterRoutes(rg, db)
	visatype.RegisterRoutes(rg, db)
	promo.RegisterRoutes(rg, db)
}
