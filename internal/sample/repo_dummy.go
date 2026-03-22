package sample

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// sampleRepositoryDummy is an in-memory implementation of SampleRepository for testing.
type sampleRepositoryDummy struct {
	mu      sync.RWMutex
	store   map[int32]*Sample
	counter int32
}

func NewSampleRepositoryDummy() SampleRepository {
	return &sampleRepositoryDummy{
		store: make(map[int32]*Sample),
	}
}

func (r *sampleRepositoryDummy) List(_ context.Context, req *ListSamplesRequest) ([]*Sample, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	all := make([]*Sample, 0, len(r.store))
	for _, s := range r.store {
		all = append(all, s)
	}

	total := int64(len(all))
	start := (page - 1) * limit
	if start >= len(all) {
		return []*Sample{}, total, nil
	}
	end := start + limit
	if end > len(all) {
		end = len(all)
	}

	return all[start:end], total, nil
}

func (r *sampleRepositoryDummy) Create(_ context.Context, sample *Sample) (*Sample, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter++
	sample.ID = r.counter
	sample.CreatedAt = time.Now()
	sample.UpdatedAt = time.Now()

	cp := *sample
	r.store[cp.ID] = &cp

	return &cp, nil
}

func (r *sampleRepositoryDummy) GetByID(_ context.Context, id int32) (*Sample, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("sample %d not found", id)
	}

	cp := *s
	return &cp, nil
}

func (r *sampleRepositoryDummy) Update(_ context.Context, sample *Sample) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[sample.ID]; !ok {
		return fmt.Errorf("sample %d not found", sample.ID)
	}

	sample.UpdatedAt = time.Now()
	cp := *sample
	r.store[cp.ID] = &cp

	return nil
}

func (r *sampleRepositoryDummy) Delete(_ context.Context, id int32) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("sample %d not found", id)
	}

	delete(r.store, id)
	return nil
}
