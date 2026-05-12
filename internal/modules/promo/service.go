package promo

import (
	"errors"
	"strings"

	"github.com/lucsky/cuid"

	"github.com/username/gin-gorm-api/internal/schema"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func (s *Service) Create(input CreateDTO) (schema.Promo, error) {
	code := strings.TrimSpace(input.Code)
	if code == "" {
		return schema.Promo{}, errors.New("code is required")
	}
	if input.ValidTo.Before(input.ValidFrom) {
		return schema.Promo{}, errors.New("valid_to must be greater than or equal to valid_from")
	}

	item := schema.Promo{
		ID:          cuid.New(),
		Code:        code,
		Description: normalizeOptionalString(input.Description),
		Discount:    input.Discount,
		ValidFrom:   input.ValidFrom,
		ValidTo:     input.ValidTo,
		IsActive:    input.IsActive == nil || *input.IsActive,
	}
	if err := s.repo.Create(&item); err != nil {
		return schema.Promo{}, err
	}
	return item, nil
}

func (s *Service) List(filter Filter) ([]schema.Promo, error) {
	return s.repo.List(filter)
}

func (s *Service) GetByID(id string) (schema.Promo, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id string, input UpdateDTO) (schema.Promo, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return schema.Promo{}, err
	}

	updates := map[string]interface{}{}
	nextValidFrom := item.ValidFrom
	nextValidTo := item.ValidTo

	if input.Code != nil {
		code := strings.TrimSpace(*input.Code)
		if code == "" {
			return schema.Promo{}, errors.New("code cannot be empty")
		}
		updates["code"] = code
	}
	if input.Description != nil {
		updates["description"] = normalizeOptionalString(input.Description)
	}
	if input.Discount != nil {
		if *input.Discount < 0 {
			return schema.Promo{}, errors.New("discount cannot be negative")
		}
		updates["discount"] = *input.Discount
	}
	if input.ValidFrom != nil {
		nextValidFrom = *input.ValidFrom
		updates["valid_from"] = *input.ValidFrom
	}
	if input.ValidTo != nil {
		nextValidTo = *input.ValidTo
		updates["valid_to"] = *input.ValidTo
	}
	if nextValidTo.Before(nextValidFrom) {
		return schema.Promo{}, errors.New("valid_to must be greater than or equal to valid_from")
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if len(updates) == 0 {
		return item, nil
	}
	return s.repo.Update(id, updates)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
