package sample

import (
	"context"

	"gorm.io/gorm"
)

type sampleRepository struct {
	db *gorm.DB
}

func NewSampleRepository(db *gorm.DB) SampleRepository {
	return &sampleRepository{db: db}
}

func (r *sampleRepository) List(ctx context.Context, req *ListSamplesRequest) ([]*Sample, int64, error) {
	var models []*SampleModel
	var total int64

	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Model(&SampleModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	samples := make([]*Sample, len(models))
	for i, m := range models {
		samples[i] = m.ToDomain()
	}

	return samples, total, nil
}

func (r *sampleRepository) Create(ctx context.Context, sample *Sample) (*Sample, error) {
	model := fromDomain(sample)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *sampleRepository) GetByID(ctx context.Context, id int32) (*Sample, error) {
	var model SampleModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *sampleRepository) Update(ctx context.Context, sample *Sample) error {
	model := fromDomain(sample)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *sampleRepository) Delete(ctx context.Context, id int32) error {
	return r.db.WithContext(ctx).Delete(&SampleModel{}, id).Error
}
