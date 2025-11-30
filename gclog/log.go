package gclog

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"strings"
)

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func ContextValueLogger(ctx context.Context) *slog.Logger {
	vl, ok := ctx.Value(loggerContextKey).(*slog.Logger)
	if ok {
		return vl
	}
	return slog.Default()
}

func ContextWithNewTracedLogger(ctx context.Context) (context.Context, *slog.Logger) {
	logger := slog.With(traceIdLoggerArg, generateTraceId())
	return ContextWithLogger(ctx, logger), logger
}

func ContextWithNewTracedLoggerArgs(ctx context.Context, args ...any) (context.Context, *slog.Logger) {
	loggerArgs := make([]any, 0, len(args)+2)
	loggerArgs = append(loggerArgs, traceIdLoggerArg, generateTraceId())
	loggerArgs = append(loggerArgs, args...)
	logger := slog.With(loggerArgs...)
	return ContextWithLogger(ctx, logger), logger
}

func generateTraceId() string {
	cs := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	rand.Shuffle(len(cs), func(i, j int) {
		cs[i], cs[j] = cs[j], cs[i]
	})
	return strings.Join(cs[:traceIdSize], "")

}
