//go:build wireinject
// +build wireinject

// wire.go defines Wire injectors.
// Run `make wire` to regenerate wire_gen.go.

package main

import (
	"my-gift/configs"
	"my-gift/internal/infra"
	"my-gift/internal/sample"

	"github.com/google/wire"
)

var infraSet = wire.NewSet(
	infra.NewLogger,
	infra.NewDatabase,
)

var loggerSet = wire.NewSet(
	infra.NewLogger,
)

var sampleSet = wire.NewSet(
	sample.ProvideRepository,
	sample.ProvideService,
	sample.ProvideController,
	sample.ProvideGRPCHandler,
)

var dummySampleSet = wire.NewSet(
	sample.NewSampleRepositoryDummy,
	sample.ProvideService,
	sample.ProvideController,
	sample.ProvideGRPCHandler,
)

// InitializeApp wires the full app with PostgreSQL.
func InitializeApp(cfg *configs.Config) (*App, error) {
	wire.Build(infraSet, sampleSet, NewApp)
	return nil, nil
}

// InitializeAppDummy wires the app with an in-memory repository (no DB).
func InitializeAppDummy(cfg *configs.Config) (*App, error) {
	wire.Build(loggerSet, dummySampleSet, NewApp)
	return nil, nil
}
