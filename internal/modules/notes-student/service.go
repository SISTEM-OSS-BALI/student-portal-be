package notesstudent

import "github.com/username/gin-gorm-api/internal/schema"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input CreateDTO) (schema.NoteStudent, error) {
	note := schema.NoteStudent{
		UserID:  input.UserID,
		Content: input.Content,
	}
	if err := s.repo.Create(&note); err != nil {
		return schema.NoteStudent{}, err
	}
	return note, nil
}

func (s *Service) List() ([]schema.NoteStudent, error) {
	return s.repo.List()
}

func (s *Service) ListByUserID(userID string) ([]schema.NoteStudent, error) {
	return s.repo.ListByUserID(userID)
}

func (s *Service) GetByID(id string) (schema.NoteStudent, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.NoteStudent, error) {
	note, err := s.repo.GetByID(id)
	if err != nil {
		return schema.NoteStudent{}, err
	}
	if input.Content != nil {
		note.Content = *input.Content
	}
	if err := s.repo.Update(&note); err != nil {
		return schema.NoteStudent{}, err
	}
	return note, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
