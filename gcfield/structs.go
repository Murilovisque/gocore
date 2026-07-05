package gcfield

import "github.com/Murilovisque/gocore/gcopt"

type FieldNameOrdered string

func (f FieldNameOrdered) String() string {
	return string(f)
}

type fieldNameOrderedParserImpl struct {
	parser func(name string) (parsedValue gcopt.Optional[FieldNameOrdered], err error)
}

func (f fieldNameOrderedParserImpl) ParseFieldNameOrdered(name string) (gcopt.Optional[FieldNameOrdered], error) {
	return f.parser(name)
}

type FieldParser struct { //todo: REMOVE
	allowedNames []string
	currentName  string
	currentValue any
	parser       func(name, value string) (parsedValue any, valid bool, err error)
	stringValue  func(name string, value any) string
}

func (f *FieldParser) AllowedNames() []string {
	return f.allowedNames
}

func (f *FieldParser) Name() string {
	return f.currentName
}

func (f *FieldParser) IsValid() bool {
	return f.currentName != ""
}

func (f *FieldParser) Parse(name, value string) (bool, error) {
	pv, ok, err := f.parser(name, value)
	if err != nil || !ok {
		return ok, err
	}
	f.currentName = name
	f.currentValue = pv
	return true, nil
}

func (f *FieldParser) Value() any {
	return f.currentValue
}

func (f *FieldParser) String() string {
	if f.IsValid() {
		return f.stringValue(f.currentName, f.currentValue)
	}
	return "FieldParser.Empty"
}
