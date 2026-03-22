package sample

import (
	"context"
	"time"
)

// Domain models

type Sample struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   int32     `json:"created_by"`
	UpdatedBy   int32     `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserDetails struct {
	UserID int32
	Email  string
}

// Request / Response DTOs

type ListSamplesRequest struct {
	Page  int `json:"page" url:"page"`
	Limit int `json:"limit" url:"limit"`
}

type ListSamplesResponse struct {
	Data  []*Sample `json:"data"`
	Total int64     `json:"total"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
}

type CreateSampleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	CreatedBy   int32  `json:"created_by"`
}

type UpdateSampleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	UpdatedBy   int32  `json:"updated_by"`
}

type PatchSampleRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	UpdatedBy   int32   `json:"updated_by"`
}

// Service interface

type SampleService interface {
	ListSamples(ctx context.Context, req *ListSamplesRequest) (*ListSamplesResponse, error)
	CreateSample(ctx context.Context, req *CreateSampleRequest) (*Sample, error)
	GetSample(ctx context.Context, id int32) (*Sample, error)
	UpdateSample(ctx context.Context, id int32, req *UpdateSampleRequest) error
	PatchSample(ctx context.Context, id int32, req *PatchSampleRequest) (*Sample, error)
	DeleteSample(ctx context.Context, id int32) error
}

// Repository interface

type SampleRepository interface {
	List(ctx context.Context, req *ListSamplesRequest) ([]*Sample, int64, error)
	Create(ctx context.Context, sample *Sample) (*Sample, error)
	GetByID(ctx context.Context, id int32) (*Sample, error)
	Update(ctx context.Context, sample *Sample) error
	Delete(ctx context.Context, id int32) error
}
