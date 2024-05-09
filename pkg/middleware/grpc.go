package middleware

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/akbariandev/pacassistant/pkg/logger"
	"github.com/getsentry/sentry-go"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type validator interface {
	ValidateAll() error
}

type validatorLegacy interface {
	Validate() error
}

// New create chained middleware for grpc.
func New(middlewares ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		middlewares...,
	)
}

// GrpcRecovery recovery panics
func GrpcRecovery(logger logger.Logger) grpc.UnaryServerInterceptor {
	rec := func(p interface{}) (err error) {
		err = status.Errorf(codes.Unknown, "%v", p)
		logger.Error(true, "recovery: panic triggered", "error", err)

		return
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(rec),
	}

	return grpc_recovery.UnaryServerInterceptor(opts...)
}

// GrpcValidator validate your message fields, for user validator
// please check https://github.com/envoyproxy/protoc-gen-validate
func GrpcValidator() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		switch in := req.(type) {
		case validator:
			if err = in.ValidateAll(); err != nil {
				return nil, err
			}
		case validatorLegacy:
			if err = in.Validate(); err != nil {
				return nil, err
			}
		}

		return handler(ctx, req)
	}
}

func GRPCLogging(logger logger.Logger) grpc.UnaryServerInterceptor {
	logFunc := logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		logger.Log(ctx, false, slog.Level(lvl), msg, fields...)
	})

	opts := []logging.Option{
		logging.WithLevels(logging.DefaultServerCodeToLevel),
		logging.WithLogOnEvents(logging.FinishCall, logging.StartCall),
	}

	return logging.UnaryServerInterceptor(logFunc, opts...)
}

// GrpcSentryPerformance track request performance in sentry performance
func GrpcSentryPerformance(client *sentry.Client, opts ...Option) grpc.UnaryServerInterceptor {
	o := newConfig(opts)
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		hub := sentry.NewHub(client, sentry.NewScope())
		ctx = sentry.SetHubOnContext(ctx, hub)

		md, _ := metadata.FromIncomingContext(ctx) // nil check in continueFromGrpcMetadata
		span := sentry.StartSpan(ctx, "grpc.server", continueFromGrpcMetadata(md),
			sentry.WithTransactionName(info.FullMethod))
		ctx = span.Context()
		defer span.Finish()

		reqBytes, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		hub.Scope().SetRequestBody(reqBytes)

		resp, err := handler(ctx, req)
		if err != nil && o.ReportOn(err) {
			tags := grpc_tags.Extract(ctx)
			for k, v := range tags.Values() {
				hub.Scope().SetTag(k, v.(string))
			}
		}
		span.Status = toSpanStatus(status.Code(err))

		return resp, err
	}
}
