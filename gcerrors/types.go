package gcerrors

import (
	"context"
	"fmt"

	"github.com/Murilovisque/gocore/gctxt"
)

type ErrorApp struct {
	msg            gctxt.Txt
	msgDetails     gctxt.ListTxt
	httpStatusCode int
}

func (e *ErrorApp) Error() string {
	msg := e.msg.Content()
	if len(e.msgDetails) > 0 {
		return fmt.Sprintf("%s: %v", msg, e.msgDetails.Content())
	}
	return msg
}

func (e *ErrorApp) ContextError(ctx context.Context) string {
	msg := e.ContextGeneralError(ctx)
	if len(e.msgDetails) > 0 {
		return fmt.Sprintf("%s: %v", msg, e.ContextErrorDetails(ctx))
	}
	return msg
}

func (e *ErrorApp) GeneralError() string {
	return e.msg.Content()
}

func (e *ErrorApp) ContextGeneralError(ctx context.Context) string {
	return e.msg.ContentByContext(ctx)
}

func (e *ErrorApp) ErrorDetails() []string {
	return e.msgDetails.Content()
}

func (e *ErrorApp) ContextErrorDetails(ctx context.Context) []string {
	return e.msgDetails.ContentByContext(ctx)
}

func (e *ErrorApp) HttpStatus() int {
	return e.httpStatusCode
}
