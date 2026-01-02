package gcvalid

import (
	"testing"

	"github.com/Murilovisque/gocore/gctxt"
)

func TestEmailValidation(t *testing.T) {
	validator := EmailValidation(gctxt.Build("campo", nil))
	testArgs := []struct {
		email string
		valid bool
	}{
		{"", false},
		{" ", false},
		{"mail", false},
		{"mail@", false},
		{"@", false},
		{"@domain", false},
		{"mail@domain", false},
		{"  @  .com", false},
		{"mail@domain.com", true},
		{"mail@domain.com.br", true},
	}
	for _, ta := range testArgs {
		valid, _ := validator(ta.email)
		if valid != ta.valid {
			t.Fatalf("expected validation '%t' for argument '%s', but received '%t'", ta.valid, ta.email, valid)
		}
	}
}

func TestGreaterThanValidation(t *testing.T) {
	validator := GreaterThanValidation(gctxt.Build("campo", nil), 10)
	testArgs := []struct {
		value int
		valid bool
	}{
		{9, false},
		{10, false},
		{11, true},
	}
	for _, ta := range testArgs {
		valid, _ := validator(ta.value)
		if valid != ta.valid {
			t.Fatalf("expected validation '%t' for argument '%d', but received '%t'", ta.valid, ta.value, valid)
		}
	}
}

func TestGreaterOrEqualValidation(t *testing.T) {
	validator := GreaterOrEqualValidation(gctxt.Build("campo", nil), 10)
	testArgs := []struct {
		value int
		valid bool
	}{
		{9, false},
		{10, true},
		{11, true},
	}
	for _, ta := range testArgs {
		valid, _ := validator(ta.value)
		if valid != ta.valid {
			t.Fatalf("expected validation '%t' for argument '%d', but received '%t'", ta.valid, ta.value, valid)
		}
	}
}

func TestLessThanValidation(t *testing.T) {
	validator := LessThanValidation(gctxt.Build("campo", nil), 10)
	testArgs := []struct {
		value int
		valid bool
	}{
		{9, true},
		{10, false},
		{11, false},
	}
	for _, ta := range testArgs {
		valid, _ := validator(ta.value)
		if valid != ta.valid {
			t.Fatalf("expected validation '%t' for argument '%d', but received '%t'", ta.valid, ta.value, valid)
		}
	}
}

func TestLessOrEqualValidation(t *testing.T) {
	validator := LessOrEqualValidation(gctxt.Build("campo", nil), 10)
	testArgs := []struct {
		value int
		valid bool
	}{
		{9, true},
		{10, true},
		{11, false},
	}
	for _, ta := range testArgs {
		valid, _ := validator(ta.value)
		if valid != ta.valid {
			t.Fatalf("expected validation '%t' for argument '%d', but received '%t'", ta.valid, ta.value, valid)
		}
	}
}

func TestBetweenInclusiveValidation(t *testing.T) {
	validator := BetweenInclusiveValidation(gctxt.Build("campo", nil), 9, 11)
	testArgs := []struct {
		value int
		valid bool
	}{
		{8, false},
		{9, true},
		{10, true},
		{11, true},
		{12, false},
	}
	for _, ta := range testArgs {
		valid, _ := validator(ta.value)
		if valid != ta.valid {
			t.Fatalf("expected validation '%t' for argument '%d', but received '%t'", ta.valid, ta.value, valid)
		}
	}
}

func TestBetweenValidation(t *testing.T) {
	validator := BetweenValidation(gctxt.Build("campo", nil), 9, 11)
	testArgs := []struct {
		value int
		valid bool
	}{
		{8, false},
		{9, false},
		{10, true},
		{11, false},
		{12, false},
	}
	for _, ta := range testArgs {
		valid, _ := validator(ta.value)
		if valid != ta.valid {
			t.Fatalf("expected validation '%t' for argument '%d', but received '%t'", ta.valid, ta.value, valid)
		}
	}
}
