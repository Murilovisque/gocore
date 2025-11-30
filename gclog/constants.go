package gclog

const (
	traceIdSize      = 8
	traceIdLoggerArg = "traceId"
)

var (
	loggerContextKey = loggerContext{}
)
