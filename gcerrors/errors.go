package gcerrors

import (
	"net/http"

	"github.com/Murilovisque/gocore/gctxt"
)

// Generic errors
var (
	ErrBadRequest   = newGenericErr(gctxt.Build("parâmetro(s) inválido(s)", nil), http.StatusBadRequest)
	ErrConflict     = newGenericErr(gctxt.Build("já existe", nil), http.StatusConflict)
	ErrNotFound     = newGenericErr(gctxt.Build("não encontrado", nil), http.StatusNotFound)
	ErrUnauthorized = newGenericErr(gctxt.Build("credenciais inválidas", nil), http.StatusUnauthorized)
	ErrForbidden    = newGenericErr(gctxt.Build("sem permissão para o recurso ou ação", nil), http.StatusForbidden)
	ErrInternal     = newGenericErr(gctxt.Build("falha interna", nil), http.StatusInternalServerError)
)

func NewErrWith(errApp *ErrorApp, msgsDetails gctxt.ListTxt) *ErrorApp {
	return &ErrorApp{msg: errApp.msg, httpStatusCode: errApp.httpStatusCode, msgDetails: msgsDetails}
}

func newGenericErr(msgs gctxt.Txt, httpStatus int) *ErrorApp {
	return &ErrorApp{msg: msgs, httpStatusCode: httpStatus, msgDetails: make(gctxt.ListTxt)}
}
