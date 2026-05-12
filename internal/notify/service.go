package notify

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
	"gorm.io/gorm"
)

type Service struct {
	db          *gorm.DB
	host        string
	port        int
	username    string
	password    string
	fromName    string
	fromAddress string
	frontendURL string
}

func NewService(db *gorm.DB) *Service {
	port, _ := strconv.Atoi(strings.TrimSpace(os.Getenv("MAIL_PORT")))
	if port == 0 {
		port = 587
	}

	frontendURL := strings.TrimSpace(os.Getenv("FRONTEND_URL"))
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	return &Service{
		db:          db,
		host:        strings.TrimSpace(os.Getenv("MAIL_HOST")),
		port:        port,
		username:    strings.TrimSpace(os.Getenv("MAIL_USERNAME")),
		password:    strings.TrimSpace(os.Getenv("MAIL_PASSWORD")),
		fromName:    strings.TrimSpace(os.Getenv("MAIL_FROM_NAME")),
		fromAddress: strings.TrimSpace(os.Getenv("MAIL_FROM_ADDRESS")),
		frontendURL: frontendURL,
	}
}

func (s *Service) Enabled() bool {
	return s.host != "" &&
		s.username != "" &&
		s.password != "" &&
		s.fromAddress != ""
}

func (s *Service) frontendBaseURL() string {
	return strings.TrimRight(s.frontendURL, "/")
}

func (s *Service) SendForgotPasswordEmail(toEmail, userName, otp string, expiresAt time.Time) error {
	subject := "Kode OTP Reset Password Student Portal"
	body := fmt.Sprintf(
		"Halo %s,\n\nKami menerima permintaan reset password akun student Anda.\n\nKode OTP Anda: %s\n\nKode berlaku sampai: %s\n\nJika Anda tidak meminta reset password, abaikan email ini.",
		strings.TrimSpace(userName),
		otp,
		expiresAt.Format(time.RFC1123),
	)
	return s.send(toEmail, subject, body)
}

func (s *Service) SendMentionEmail(message schema.ChatMessage, mentionUserIDs []string) error {
	if !s.Enabled() || len(mentionUserIDs) == 0 {
		return nil
	}

	uniqueIDs := make(map[string]struct{}, len(mentionUserIDs))
	for _, id := range mentionUserIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		uniqueIDs[id] = struct{}{}
	}
	if len(uniqueIDs) == 0 {
		return nil
	}

	ids := make([]string, 0, len(uniqueIDs))
	for id := range uniqueIDs {
		ids = append(ids, id)
	}

	var users []schema.User
	if err := s.db.Select("id", "name", "email").Where("id IN ?", ids).Find(&users).Error; err != nil {
		return err
	}

	senderName := "Someone"
	if message.Sender != nil && strings.TrimSpace(message.Sender.Name) != "" {
		senderName = strings.TrimSpace(message.Sender.Name)
	}

	preview := strings.TrimSpace(ptrToString(message.Text))
	if preview == "" {
		preview = "Anda mendapatkan mention baru di percakapan."
	}

	for _, user := range users {
		if strings.TrimSpace(user.Email) == "" {
			continue
		}
		subject := "Anda di-mention di percakapan Student Portal"
		body := fmt.Sprintf(
			"Halo %s,\n\n%s menandai Anda di percakapan.\n\nPesan:\n%s\n\nSilakan buka inbox pada Student Portal untuk menindaklanjuti.",
			strings.TrimSpace(user.Name),
			senderName,
			preview,
		)
		if err := s.send(user.Email, subject, body); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) SendStudentUploadedDocumentEmail(studentID string, docs []schema.AnswerDocument) error {
	if !s.Enabled() || strings.TrimSpace(studentID) == "" || len(docs) == 0 {
		return nil
	}

	var student schema.User
	if err := s.db.Select("id", "name", "email").Where("id = ?", studentID).First(&student).Error; err != nil {
		return err
	}

	var admissions []schema.User
	if err := s.db.Select("id", "name", "email").
		Where("role = ?", schema.UserRoleAdmission).
		Find(&admissions).Error; err != nil {
		return err
	}

	if len(admissions) == 0 {
		return nil
	}

	var names []string
	for _, doc := range docs {
		if strings.TrimSpace(ptrToString(doc.FileName)) != "" {
			names = append(names, ptrToString(doc.FileName))
			continue
		}
		names = append(names, fmt.Sprintf("DocumentID:%s", doc.DocumentID))
	}

	docSummary := strings.Join(names, ", ")
	for _, adm := range admissions {
		if strings.TrimSpace(adm.Email) == "" {
			continue
		}
		subject := "Student upload dokumen baru"
		body := fmt.Sprintf(
			"Halo %s,\n\nStudent %s baru saja upload dokumen.\n\nDaftar dokumen: %s\n\nSilakan cek dashboard admission untuk review.",
			strings.TrimSpace(adm.Name),
			strings.TrimSpace(student.Name),
			docSummary,
		)
		if err := s.send(adm.Email, subject, body); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) send(toEmail, subject, body string) error {
	if !s.Enabled() {
		return nil
	}

	toEmail = strings.TrimSpace(toEmail)
	if toEmail == "" {
		return nil
	}

	fromName := s.fromName
	if fromName == "" {
		fromName = "Student Portal"
	}

	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, s.fromAddress))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", toEmail))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	return smtp.SendMail(addr, auth, s.fromAddress, []string{toEmail}, msg.Bytes())
}

func ptrToString(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}
