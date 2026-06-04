package gcfield

func NewFieldParser(allowedNames []string, parser func(name, value string) (parsedValue any, valid bool, err error), stringValue func(name string, value any) string) FieldParser {
	return FieldParser{
		allowedNames: allowedNames,
		parser:       parser,
		stringValue:  stringValue,
	}
}
