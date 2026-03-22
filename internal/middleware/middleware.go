package middleware

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

// ─── 1. WrapRouter ───────────────────────────────────────────────────────────
// Low-level nhất, nhận http.ResponseWriter và http.Request gốc.
// Chạy trước tất cả mọi thứ kể cả Iris router.
// Dùng cho: CORS, rate limiting, chặn IP...

var requestCount int64

func WrapRouter(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
	atomic.AddInt64(&requestCount, 1)

	// CORS headers cho mọi request
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Xử lý OPTIONS preflight mà không đi vào Iris router
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	router(w, r)
}

// ─── 2. UseRouter ────────────────────────────────────────────────────────────
// Chạy sau WrapRouter, trước khi Iris tìm route phù hợp.
// Có đầy đủ iris.Context.
// Dùng cho: request ID, logging, recover from panic...

func UseRouter(logger *zap.Logger) iris.Handler {
	return func(ctx iris.Context) {
		start := time.Now()
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx.Values().Set("requestID", requestID)
		ctx.Header("X-Request-ID", requestID)

		logger.Info("[UseRouter] incoming request",
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.Path()),
			zap.String("requestID", requestID),
			zap.String("ip", ctx.RemoteAddr()),
		)

		ctx.Next()

		logger.Info("[UseRouter] request completed",
			zap.String("requestID", requestID),
			zap.Int("status", ctx.GetStatusCode()),
			zap.Duration("duration", time.Since(start)),
		)
	}
}

// ─── 3. UseGlobal ────────────────────────────────────────────────────────────
// Chạy cho MỌI route kể cả error pages (4xx, 5xx).
// Dùng cho: auth global, logging toàn bộ...

func UseGlobal(logger *zap.Logger) iris.Handler {
	return func(ctx iris.Context) {
		logger.Debug("[UseGlobal] before handler",
			zap.String("path", ctx.Path()),
		)

		ctx.Next()

		logger.Debug("[UseGlobal] after handler",
			zap.Int("status", ctx.GetStatusCode()),
		)
	}
}

// ─── 4. Use ──────────────────────────────────────────────────────────────────
// Chỉ chạy cho route thường, KHÔNG chạy cho error handler (4xx, 5xx).
// Dùng cho: auth, business logic middleware...

func Use(logger *zap.Logger) iris.Handler {
	return func(ctx iris.Context) {
		logger.Debug("[Use] middleware running",
			zap.String("path", ctx.Path()),
		)

		// Ví dụ: kiểm tra API key đơn giản
		// apiKey := ctx.GetHeader("X-API-Key")
		// if apiKey == "" {
		// 	ctx.StopWithStatus(iris.StatusUnauthorized)
		// 	return
		// }

		ctx.Next()
	}
}

// ─── 5. UseError ─────────────────────────────────────────────────────────────
// Chỉ chạy cho error handler (OnErrorCode).
// Dùng cho: format lỗi thống nhất, error logging...

func UseError(logger *zap.Logger) iris.Handler {
	return func(ctx iris.Context) {
		ctx.Next()

		logger.Warn("[UseError] error response",
			zap.Int("status", ctx.GetStatusCode()),
			zap.String("path", ctx.Path()),
		)
	}
}

// ─── 6. Done ─────────────────────────────────────────────────────────────────
// Chạy SAU handler, chỉ cho route thường (không phải error handler).
// Dùng cho: cleanup resources, audit log...

func Done(logger *zap.Logger) iris.Handler {
	return func(ctx iris.Context) {
		logger.Debug("[Done] cleanup after normal handler",
			zap.String("path", ctx.Path()),
			zap.Int("status", ctx.GetStatusCode()),
		)
	}
}

// ─── 7. DoneGlobal ───────────────────────────────────────────────────────────
// Chạy SAU handler, cho MỌI route kể cả error handler.
// Dùng cho: cleanup toàn bộ, metrics...

func DoneGlobal(logger *zap.Logger) iris.Handler {
	return func(ctx iris.Context) {
		requestID := ctx.Values().GetString("requestID")
		logger.Debug("[DoneGlobal] cleanup for all routes",
			zap.String("requestID", requestID),
			zap.String("path", ctx.Path()),
			zap.Int64("totalRequests", atomic.LoadInt64(&requestCount)),
		)
	}
}
