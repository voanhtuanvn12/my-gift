package sample

import "time"

// SampleModel is the GORM model for the samples table.
type SampleModel struct {
	ID          int32     `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"not null"`
	Description string
	CreatedBy   int32
	UpdatedBy   int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SampleModel) TableName() string {
	return "samples"
}

func (m *SampleModel) ToDomain() *Sample {
	return &Sample{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func fromDomain(s *Sample) *SampleModel {
	return &SampleModel{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		CreatedBy:   s.CreatedBy,
		UpdatedBy:   s.UpdatedBy,
	}
}
