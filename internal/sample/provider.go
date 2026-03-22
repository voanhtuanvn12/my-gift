package sample

import "gorm.io/gorm"

// ProvideRepository creates a new SampleRepository.
func ProvideRepository(db *gorm.DB) SampleRepository {
	return NewSampleRepository(db)
}

// ProvideService creates a new SampleService.
func ProvideService(repo SampleRepository) SampleService {
	return NewSampleService(repo)
}

// ProvideController creates a new HTTP Controller.
func ProvideController(svc SampleService) *Controller {
	return &Controller{SampleSvc: svc}
}

// ProvideGRPCHandler creates a new gRPC Handler.
func ProvideGRPCHandler(svc SampleService) *GRPCHandler {
	return NewGRPCHandler(svc)
}
