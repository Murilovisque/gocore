package gcerror

import (
	"errors"
	"slices"
	"testing"

	"github.com/Murilovisque/gocore/gctxt"
)

func TestErrors(t *testing.T) {
	gctxt.DefaultLang = gctxt.LangBrazil
	expectedDetails := []string{"Mensagem"}
	err := NewErrWith(ErrBadRequest, gctxt.NewList(expectedDetails, nil))

	if err.GeneralError() != ErrBadRequest.GeneralError() {
		t.Fatalf("expected '%s', but '%s", ErrBadRequest.GeneralError(), err.GeneralError())
	}
	if !slices.Equal(err.ErrorDetails(), expectedDetails) {
		t.Fatalf("expected '%s', but '%s", ErrBadRequest.GeneralError(), err.GeneralError())
	}

	var specificErr *ErrorApp
	if !errors.As(err, &specificErr) {
		t.Fatalf("expected specific error, but '%s", specificErr.Error())
	} else if err.Error() != specificErr.Error() {
		t.Fatalf("expected '%s', but '%s", err.Error(), specificErr.Error())
	}
}
