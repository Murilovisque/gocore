package gcfield

import "github.com/Murilovisque/gocore/gcopt"

func NewFieldParser(allowedNames []string, parser func(name, value string) (parsedValue any, valid bool, err error), stringValue func(name string, value any) string) FieldParser {
	return FieldParser{
		allowedNames: allowedNames,
		parser:       parser,
		stringValue:  stringValue,
	}
}

func NewFieldNameOrderedParser(parser func(name string) (parsedValue gcopt.Optional[FieldNameOrdered], err error)) FieldNameOrderedParser {
	return fieldNameOrderedParserImpl{
		parser: parser,
	}
}
