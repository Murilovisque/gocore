package gcoptional

import (
	"bytes"
	"encoding/json"
)

type Optional[T any] struct {
	vl     T
	exists bool
}

func (o Optional[T]) IsPresent() bool {
	return o.exists
}

func (o Optional[T]) Take() (T, bool) {
	return o.vl, o.exists
}

func (o Optional[T]) MustTake() T {
	if o.exists {
		return o.vl
	}
	panic("gcoptional: attempt to MustTake() value from invalid Optional")
}

func (o Optional[T]) TakeOrElse(orElseFunc func() T) T {
	if o.exists {
		return o.vl
	}
	return orElseFunc()
}

func (o Optional[T]) TakeOrError(orElseErrFunc func() error) (T, error) {
	if o.exists {
		return o.vl, nil
	}
	return o.vl, orElseErrFunc()
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.IsPresent() {
		return jsonNull, nil
	}
	marshal, err := json.Marshal(o.MustTake())
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, jsonNull) {
		*o = EmtpyValue[T]()
		return nil
	}
	var v T
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*o = FromValue(v)
	return nil
}
