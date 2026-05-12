package promo

import (
	"time"

	"github.com/username/gin-gorm-api/internal/schema"
)

type CreateDTO struct {
	Code        string    `json:"code" binding:"required"`
	Description *string   `json:"description"`
	Discount    float64   `json:"discount" binding:"required,gte=0"`
	ValidFrom   time.Time `json:"valid_from" binding:"required"`
	ValidTo     time.Time `json:"valid_to" binding:"required"`
	IsActive    *bool     `json:"is_active"`
}

type UpdateDTO struct {
	Code        *string    `json:"code"`
	Description *string    `json:"description"`
	Discount    *float64   `json:"discount"`
	ValidFrom   *time.Time `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
	IsActive    *bool      `json:"is_active"`
}

type Filter struct {
	ActiveOnly bool
}

type ResponseDTO struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Description *string   `json:"description,omitempty"`
	Discount    float64   `json:"discount"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidTo     time.Time `json:"valid_to"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewResponseDTO(item schema.Promo) ResponseDTO {
	return ResponseDTO{
		ID:          item.ID,
		Code:        item.Code,
		Description: item.Description,
		Discount:    item.Discount,
		ValidFrom:   item.ValidFrom,
		ValidTo:     item.ValidTo,
		IsActive:    item.IsActive,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func NewResponseListDTO(items []schema.Promo) []ResponseDTO {
	out := make([]ResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, NewResponseDTO(item))
	}
	return out
}
