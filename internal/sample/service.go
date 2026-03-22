package sample

import "context"

type sampleService struct {
	repo SampleRepository
}

func NewSampleService(repo SampleRepository) SampleService {
	return &sampleService{repo: repo}
}

func (s *sampleService) ListSamples(ctx context.Context, req *ListSamplesRequest) (*ListSamplesResponse, error) {
	data, total, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &ListSamplesResponse{
		Data:  data,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func (s *sampleService) CreateSample(ctx context.Context, req *CreateSampleRequest) (*Sample, error) {
	sample := &Sample{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
	}
	return s.repo.Create(ctx, sample)
}

func (s *sampleService) GetSample(ctx context.Context, id int32) (*Sample, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *sampleService) UpdateSample(ctx context.Context, id int32, req *UpdateSampleRequest) error {
	sample := &Sample{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		UpdatedBy:   req.UpdatedBy,
	}
	return s.repo.Update(ctx, sample)
}

func (s *sampleService) PatchSample(ctx context.Context, id int32, req *PatchSampleRequest) (*Sample, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	existing.UpdatedBy = req.UpdatedBy

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *sampleService) DeleteSample(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}
