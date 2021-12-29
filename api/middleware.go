package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/rs/zerolog/log"
)

const DefaultPageLimit = 50

type CtxKey string

const (
	CtxKeyLimit  CtxKey = "limit"
	CtxKeyOffset CtxKey = "offset"
	CtxKeyUser   CtxKey = "user"
)

var (
	urlHitCount *prometheus.CounterVec
	urlLatency  *prometheus.SummaryVec
)

func Logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			dur := fmt.Sprintf("%dms", time.Duration(time.Since(start).Milliseconds()))

			log.Info().
				Str("method", r.Method).
				Str("host", r.Host).
				Str("uri", r.RequestURI).
				Str("proto", r.Proto).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Str("duration", dur).Send()
		}()
		next.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}

type UserAccess interface {
	Auth(ctx context.Context, username, password string) bool
}

func Authenticate(ua UserAccess) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()

			if !ok {
				authErr(w)
				return
			}

			if !ua.Auth(r.Context(), username, password) {
				authErr(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func authErr(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func ConfigureMetrics() {
	urlHitCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_hit_count",
			Help: "Number of times the given url was hit",
		},
		[]string{"method", "url"},
	)
	urlLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "url_latency",
			Help:       "The latency quantiles for the given URL",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "url"},
	)

	prometheus.MustRegister(urlHitCount)
	prometheus.MustRegister(urlLatency)
}

func Metrics(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			ctx := chi.RouteContext(r.Context())

			if len(ctx.RoutePatterns) > 0 {
				dur := float64(time.Since(start).Milliseconds())
				urlLatency.WithLabelValues(ctx.RouteMethod, ctx.RoutePatterns[0]).Observe(dur)
				urlHitCount.WithLabelValues(ctx.RouteMethod, ctx.RoutePatterns[0]).Inc()
			}
		}()

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
