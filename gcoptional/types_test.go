package gcoptional

import (
	"encoding/json"
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
	opn := EmtpyValue[int]()
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

func TestJsonMarshal(t *testing.T) {
	testArgs := []struct {
		jsonArgument any
		jsonText     string
	}{
		{
			jsonArgument: struct {
				Name Optional[string] `json:"name"`
			}{
				Name: FromValue("Murilo"),
			},
			jsonText: `{"name":"Murilo"}`,
		},
		{
			jsonArgument: struct {
				Age Optional[int] `json:"age"`
			}{
				Age: FromValue(2),
			},
			jsonText: `{"age":2}`,
		},
		{
			jsonArgument: struct {
				Age Optional[int] `json:"age"`
			}{
				Age: EmtpyValue[int](),
			},
			jsonText: `{"age":null}`,
		},
		{
			jsonArgument: testComplexJson{
				Name: "teste",
				SubValue: FromValue(testSubComplex{
					Age: 10,
				}),
			},
			jsonText: `{"name":"teste","subValue":{"age":10}}`,
		},
		{
			jsonArgument: []testComplexJson{
				{
					Name: "teste",
					SubValue: FromValue(testSubComplex{
						Age: 10,
					}),
				},
			},
			jsonText: `[{"name":"teste","subValue":{"age":10}}]`,
		},
	}
	for _, a := range testArgs {
		txtJson, err := json.Marshal(&a.jsonArgument)
		if err != nil {
			t.Fatalf("the error is not expected, but '%s'", err.Error())
		}
		result := string(txtJson)
		if result != a.jsonText {
			t.Fatalf("expected '%s', but '%s'", a.jsonText, result)
		}
	}
}

func TestJsonUnmarshal(t *testing.T) {
	testArgs := []struct {
		jsonStructExpected testComplexJson
		jsonStructEmpty    testComplexJson
		jsonString         string
	}{
		{
			jsonStructExpected: testComplexJson{
				Name: "teste",
				SubValue: FromValue(testSubComplex{
					Age: 10,
				}),
			},
			jsonStructEmpty: testComplexJson{},
			jsonString:      `{"name":"teste","subValue":{"age":10}}`,
		},
		{
			jsonStructExpected: testComplexJson{
				Name:     "",
				SubValue: EmtpyValue[testSubComplex](),
			},
			jsonStructEmpty: testComplexJson{},
			jsonString:      `{"name":"","subValue":null}`,
		},
	}
	for _, a := range testArgs {
		err := json.Unmarshal([]byte(a.jsonString), &a.jsonStructEmpty)
		if err != nil {
			t.Fatalf("the error is not expected, but '%s'", err.Error())
		}
		if !reflect.DeepEqual(a.jsonStructExpected, a.jsonStructEmpty) {
			t.Fatalf("expected '%#v', but '%#v'", a.jsonStructExpected, a.jsonStructEmpty)
		}
	}
}

type testComplexJson struct {
	Name     string                   `json:"name"`
	SubValue Optional[testSubComplex] `json:"subValue"`
}

type testSubComplex struct {
	Age int `json:"age"`
}
