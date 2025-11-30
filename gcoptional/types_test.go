package gcoptional

import (
	"errors"
	"reflect"
	"testing"
)

func TestOptional(t *testing.T) {
	var st struct {
		id   int
		name string
	}
	st.id = 1
	st.name = "tester"
	op := FromValue(st)
	if !op.IsPresent() || !reflect.DeepEqual(st, op.MustTake()) {
		t.Fatalf("expected '%v', but '%v'", st, op.MustTake())
	}
	vl, ok := op.Take()
	if !ok || !reflect.DeepEqual(st, vl) {
		t.Fatalf("expected '%v', but '%v'", st, vl)
	}
}

func TestNone(t *testing.T) {
	// FromPointer
	var st *struct{}
	op := FromPointer(st)
	if op.IsPresent() {
		t.Fatalf("expected '%v', but '%v'", false, op)
	}
	vl, ok := op.Take()
	if ok {
		t.Fatalf("expected '%v', but '%v'", false, vl)
	}

	// None
	opn := None[int]()
	if opn.IsPresent() {
		t.Fatalf("expected '%v', but '%v'", false, opn)
	}
	const expectedValue = 10
	vln := opn.TakeOrElse(func() int { return expectedValue })
	if vln != expectedValue {
		t.Fatalf("expected '%v', but '%v'", expectedValue, vln)
	}
	_, err := opn.TakeOrError(func() error { return errors.New("failed") })
	if err == nil {
		t.Fatalf("expected not nil %s", err)
	}
}
