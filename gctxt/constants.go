package gctxt

const (
	LangBrazil Language = iota
	LangEnglish
)

var (
	DefaultLang = LangBrazil
)

var (
	languageTypeContextKey = languageTypeContext{}
)
