package gcvalidation

import (
	"github.com/Murilovisque/gocore/gcerrors"
	"github.com/Murilovisque/gocore/gctxt"
)

type Validation[T any] func(T) (bool, gctxt.Txt)

type ValidationTwoFields[T, T2 any] func(T, T2) (bool, gctxt.Txt)

type RequestValidator struct {
	msgs gctxt.ListTxt
}

func (p *RequestValidator) AddErrorMsg(msg gctxt.Txt) {
	if p.msgs == nil {
		p.msgs = make(gctxt.ListTxt)
	}
	for k, v := range msg {
		if msgs, ok := p.msgs[k]; ok {
			p.msgs[k] = append(msgs, v)
		} else {
			p.msgs[k] = []string{v}
		}
	}
}

func (p RequestValidator) HasErrors() bool {
	return len(p.msgs) > 0
}

func (p RequestValidator) ValidErrors() error {
	if p.HasErrors() {
		return gcerrors.NewErrWith(gcerrors.ErrBadRequest, p.msgs)
	}
	return nil
}
