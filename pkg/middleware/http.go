package middleware

import (
	"context"
	"expvar"
	"net/http"
	"net/http/pprof"
	"path"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// ErrorHandler convert grpc status error to http error.
func ErrorHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler,
	w http.ResponseWriter, r *http.Request, err error,
) {
	w.Header().Set("Content-Type", "application/json")
	runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, r, err)
}

// AllowCORS add cors to http handler.
func AllowCORS(h http.Handler, origins []string, customHeaders ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if len(origins) != 0 {
			if origin == "" {
				w.WriteHeader(http.StatusForbidden)

				return
			}
			if !checkOrigin(origin, origins) {
				w.WriteHeader(http.StatusForbidden)

				return
			}
		} else {
			origin = "*"
		}

		headers := []string{
			"Accept",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"ResponseType",
		}

		headers = append(headers, customHeaders...)

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)

			return
		}

		h.ServeHTTP(w, r)
	})
}

// SetRuntimeAsRootHandler set runtime mux as root handler http server mux
func SetRuntimeAsRootHandler(mux *http.ServeMux, rMux *runtime.ServeMux) *http.ServeMux {
	mux.Handle("/", rMux)

	return mux
}

// SwaggerHandler add swagger file embedded to http handler
// path, swaggerFileName (swagger.json or swagger.yaml and etc).
func SwaggerHandler(mux *http.ServeMux, swaggerFileName string, swagger []byte) *http.ServeMux {
	mux.HandleFunc("/"+swaggerFileName, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(swagger)
	})

	return mux
}

// DebuggerHandler add pprof handlers to http server mux
func DebuggerHandler(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

func checkOrigin(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || originMatchesPattern(origin, allowedOrigin) {
			return true
		}
	}

	return false
}

func originMatchesPattern(origin, pattern string) bool {
	matched, err := path.Match(pattern, origin)

	return err == nil && matched
}
