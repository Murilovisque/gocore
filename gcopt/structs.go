package gcopt

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Optional[T any] struct {
	value  T
	exists bool
}

func (o Optional[T]) IsPresent() bool {
	return o.exists
}

func (o Optional[T]) IfPresent(f func(T)) {
	if o.exists {
		f(o.value)
	}
}

func (o Optional[T]) Take() (T, bool) {
	return o.value, o.exists
}

func (o Optional[T]) MustTake() T {
	if o.exists {
		return o.value
	}
	panic("gcopt: attempt to MustTake() value from invalid Optional")
}

func (o Optional[T]) TakeOr(other T) T {
	if o.exists {
		return o.value
	}
	return other
}

func (o Optional[T]) TakeOrElse(orElseFunc func() T) T {
	if o.exists {
		return o.value
	}
	return orElseFunc()
}

func (o Optional[T]) TakeOrError(err error) (T, error) {
	if o.exists {
		return o.value, nil
	}
	return o.value, err
}

func (o Optional[T]) TakeOrErrorElse(orElseErrFunc func() error) (T, error) {
	if o.exists {
		return o.value, nil
	}
	return o.value, orElseErrFunc()
}

func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if !o.exists {
		return o
	} else if predicate(o.value) {
		return o
	}
	return Optional[T]{exists: false}
}

// json interface

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.exists {
		return jsonNull, nil
	}
	return json.Marshal(o.MustTake())
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, jsonNull) {
		*o = Empty[T]()
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*o = Of(v)
	return nil
}

// string interface
func (o Optional[T]) String() string {
	if !o.exists {
		return "Optional.Empty"
	}
	return fmt.Sprintf("Optional[%v]", o.value)
}

// SQL interfaces
func (o *Optional[T]) Scan(value any) error {
	var nt sql.Null[T]
	if err := nt.Scan(value); err != nil {
		return err
	}
	o.value = nt.V
	o.exists = nt.Valid
	return nil
}

func (o Optional[T]) Value() (driver.Value, error) {
	if !o.exists {
		return nil, nil
	}
	return o.value, nil
}
