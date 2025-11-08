package middleware

import (
	"context"
	"strings"

	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/logger"
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
	"github.com/WithSoull/platform_common/pkg/tokens"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TokenVerifierUnaryInterceptor(verifier tokens.TokenVerifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		// skip health checks
		logger.Warn(ctx, info.FullMethod)
		if strings.HasPrefix(info.FullMethod, "/grpc.health.v1.Health") ||
			strings.HasPrefix(info.FullMethod, "/grpc.reflection") {
			return handler(ctx, req)
		}

		ctx, err = VerifyToken(ctx, verifier)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	wrappedCtx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.wrappedCtx
}

func TokenVerifierStreamInterceptor(verifier tokens.TokenVerifier) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// skip health checks
		logger.Warn(ss.Context(), info.FullMethod)
		if strings.HasPrefix(info.FullMethod, "/grpc.health.v1.Health") ||
			strings.HasPrefix(info.FullMethod, "/grpc.reflection") {
			return handler(srv, ss)
		}

		ctx, err := VerifyToken(ss.Context(), verifier)
		if err != nil {
			return err
		}

		wrapped := &wrappedServerStream{
			ServerStream: ss,
			wrappedCtx:   ctx,
		}

		return handler(srv, wrapped)
	}
}

func VerifyToken(ctx context.Context, verifier tokens.TokenVerifier) (context.Context, error) {
	logger.Warn(ctx, "Verify token has been called")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, sys.NewCommonError("metadata not provided", codes.Unauthenticated)
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return ctx, sys.NewCommonError("authorization header not provided", codes.Unauthenticated)
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := verifier.VerifyAccessToken(ctx, token)
	if err != nil {
		return ctx, sys.NewCommonError("invalid token", codes.Unauthenticated)
	}

	ctxWithEmail := claimsctx.InjectUserEmail(ctx, claims.Email)
	ctxWithEmailAndUserID := claimsctx.InjectUserID(ctxWithEmail, claims.UserId)

	return ctxWithEmailAndUserID, nil
}
