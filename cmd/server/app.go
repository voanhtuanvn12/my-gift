package main

import (
	"fmt"
	"net"
	"net/http"

	_ "my-gift/docs"
	"my-gift/configs"
	samplev1 "my-gift/gen/proto/sample/v1"
	"my-gift/internal/sample"

	"github.com/fullstorydev/grpchan"
	"github.com/fullstorydev/grpchan/httpgrpc"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	swag "github.com/swaggo/swag/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const scalarHTML = `<!DOCTYPE html>
<html>
<head>
  <title>My Gift API — Docs</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
  <script id="api-reference" data-url="/openapi.json"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

// App holds the iris application and its dependencies.
type App struct {
	iris    *iris.Application
	grpcSvr *grpc.Server
	cfg     *configs.Config
	logger  *zap.Logger
}

// NewApp assembles the Iris application with all routes.
func NewApp(
	cfg *configs.Config,
	logger *zap.Logger,
	sampleCtrl *sample.Controller,
	sampleGRPC *sample.GRPCHandler,
) *App {
	// ── Native gRPC server (HTTP/2, port GRPC_PORT) ──────────────────────────
	// Used by grpcurl and native gRPC clients.
	grpcSvr := grpc.NewServer()
	samplev1.RegisterSampleServiceServer(grpcSvr, sampleGRPC)
	reflection.Register(grpcSvr) // enables grpcurl --use-reflection

	// ── Iris HTTP server (port APP_PORT) ─────────────────────────────────────
	app := iris.New()
	app.Use(iris.Compression)

	app.Get("/health", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"status": "ok"})
	})

	app.Get("/openapi.json", func(ctx iris.Context) {
		spec, err := swag.ReadDoc()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
		ctx.ContentType("application/json")
		ctx.WriteString(spec) //nolint:errcheck
	})

	app.Get("/docs", func(ctx iris.Context) {
		ctx.ContentType("text/html")
		ctx.WriteString(scalarHTML) //nolint:errcheck
	})

	// REST routes
	api := app.Party("/api/v1")
	mvc.Configure(api.Party("/samples"), func(m *mvc.Application) {
		m.Handle(sampleCtrl)
	})

	// gRPC-over-HTTP/1.1 via grpchan — for clients that can't do HTTP/2.
	// Route: POST /grpc/<package>.<Service>/<Method>
	handlers := grpchan.HandlerMap{}
	samplev1.RegisterSampleServiceServer(handlers, sampleGRPC)
	grpcMux := http.NewServeMux()
	httpgrpc.HandleServices(grpcMux.HandleFunc, "/grpc/", handlers, nil, nil)
	app.Any("/grpc/{any:path}", iris.FromStd(grpcMux))

	return &App{iris: app, grpcSvr: grpcSvr, cfg: cfg, logger: logger}
}

// Run starts both the Iris HTTP server and the native gRPC server concurrently.
func (a *App) Run() error {
	httpAddr := fmt.Sprintf("%s:%d", a.cfg.App.Host, a.cfg.App.Port)
	grpcAddr := fmt.Sprintf("%s:%d", a.cfg.App.Host, a.cfg.App.GRPCPort)

	a.logger.Info("Starting HTTP server",
		zap.String("addr", httpAddr),
		zap.String("env", a.cfg.App.Env),
		zap.String("docs", fmt.Sprintf("http://localhost:%d/docs", a.cfg.App.Port)),
		zap.String("grpc-http1", fmt.Sprintf("http://localhost:%d/grpc/", a.cfg.App.Port)),
	)
	a.logger.Info("Starting native gRPC server",
		zap.String("addr", grpcAddr),
	)

	errCh := make(chan error, 2)

	// Native gRPC (HTTP/2) — for grpcurl and native clients
	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			errCh <- fmt.Errorf("grpc listen: %w", err)
			return
		}
		errCh <- a.grpcSvr.Serve(lis)
	}()

	// Iris HTTP
	go func() {
		errCh <- a.iris.Listen(httpAddr)
	}()

	return <-errCh
}
