package server

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
)

// configure gprc logger
//
//nolint:gocyclo // this is the only way to confgure the logger middleware
func customCodeToLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK:
		return zap.DebugLevel
	case codes.Canceled:
		return zap.DebugLevel
	case codes.Unknown:
		return zap.ErrorLevel
	case codes.InvalidArgument:
		return zap.DebugLevel
	case codes.DeadlineExceeded:
		return zap.ErrorLevel
	case codes.NotFound:
		return zap.DebugLevel
	case codes.AlreadyExists:
		return zap.DebugLevel
	case codes.PermissionDenied:
		return zap.ErrorLevel
	case codes.Unauthenticated:
		return zap.DebugLevel
	case codes.ResourceExhausted:
		return zap.ErrorLevel
	case codes.FailedPrecondition:
		return zap.ErrorLevel
	case codes.Aborted:
		return zap.ErrorLevel
	case codes.OutOfRange:
		return zap.ErrorLevel
	case codes.Unimplemented:
		return zap.ErrorLevel
	case codes.Internal:
		return zap.ErrorLevel
	case codes.Unavailable:
		return zap.ErrorLevel
	case codes.DataLoss:
		return zap.ErrorLevel
	default:
		return zap.ErrorLevel
	}
}
