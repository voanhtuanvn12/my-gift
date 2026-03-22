package main

import (
	"fmt"
	_ "my-gift/docs"
	"my-gift/configs"
	"my-gift/internal/sample"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	swag "github.com/swaggo/swag/v2"
	"go.uber.org/zap"
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
	iris   *iris.Application
	cfg    *configs.Config
	logger *zap.Logger
}

// NewApp assembles the Iris application with all routes.
func NewApp(cfg *configs.Config, logger *zap.Logger, sampleCtrl *sample.Controller) *App {
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

	api := app.Party("/api/v1")
	mvc.Configure(api.Party("/samples"), func(m *mvc.Application) {
		m.Handle(sampleCtrl)
	})

	return &App{iris: app, cfg: cfg, logger: logger}
}

// Run starts the HTTP server.
func (a *App) Run() error {
	addr := fmt.Sprintf("%s:%d", a.cfg.App.Host, a.cfg.App.Port)
	a.logger.Info("Starting server",
		zap.String("addr", addr),
		zap.String("env", a.cfg.App.Env),
		zap.String("docs", fmt.Sprintf("http://localhost:%d/docs", a.cfg.App.Port)),
	)
	return a.iris.Listen(addr)
}
