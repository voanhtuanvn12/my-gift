package sample

import (
	"context"
	"errors"
	"time"

	samplev1 "my-gift/gen/proto/sample/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// GRPCHandler implements samplev1.SampleServiceServer.
type GRPCHandler struct {
	samplev1.UnimplementedSampleServiceServer
	svc SampleService
}

func NewGRPCHandler(svc SampleService) *GRPCHandler {
	return &GRPCHandler{svc: svc}
}

func (h *GRPCHandler) ListSamples(ctx context.Context, req *samplev1.ListSamplesRequest) (*samplev1.ListSamplesResponse, error) {
	result, err := h.svc.ListSamples(ctx, &ListSamplesRequest{
		Page:  int(req.Page),
		Limit: int(req.Limit),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	data := make([]*samplev1.Sample, len(result.Data))
	for i, s := range result.Data {
		data[i] = toProto(s)
	}

	return &samplev1.ListSamplesResponse{
		Data:  data,
		Total: result.Total,
		Page:  int32(result.Page),
		Limit: int32(result.Limit),
	}, nil
}

func (h *GRPCHandler) CreateSample(ctx context.Context, req *samplev1.CreateSampleRequest) (*samplev1.Sample, error) {
	result, err := h.svc.CreateSample(ctx, &CreateSampleRequest{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return toProto(result), nil
}

func (h *GRPCHandler) GetSample(ctx context.Context, req *samplev1.GetSampleRequest) (*samplev1.Sample, error) {
	result, err := h.svc.GetSample(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "sample %d not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return toProto(result), nil
}

func (h *GRPCHandler) UpdateSample(ctx context.Context, req *samplev1.UpdateSampleRequest) (*samplev1.UpdateSampleResponse, error) {
	err := h.svc.UpdateSample(ctx, req.Id, &UpdateSampleRequest{
		Name:        req.Name,
		Description: req.Description,
		UpdatedBy:   req.UpdatedBy,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "sample %d not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &samplev1.UpdateSampleResponse{}, nil
}

func (h *GRPCHandler) PatchSample(ctx context.Context, req *samplev1.PatchSampleRequest) (*samplev1.Sample, error) {
	result, err := h.svc.PatchSample(ctx, req.Id, &PatchSampleRequest{
		Name:        req.Name,
		Description: req.Description,
		UpdatedBy:   req.UpdatedBy,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "sample %d not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return toProto(result), nil
}

func (h *GRPCHandler) DeleteSample(ctx context.Context, req *samplev1.DeleteSampleRequest) (*samplev1.DeleteSampleResponse, error) {
	err := h.svc.DeleteSample(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "sample %d not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &samplev1.DeleteSampleResponse{}, nil
}

func toProto(s *Sample) *samplev1.Sample {
	return &samplev1.Sample{
		Id:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		CreatedBy:   s.CreatedBy,
		UpdatedBy:   s.UpdatedBy,
		CreatedAt:   s.CreatedAt.Unix(),
		UpdatedAt:   s.UpdatedAt.Unix(),
	}
}

// Ensure time.Time zero value doesn't cause issues
var _ = time.Now
