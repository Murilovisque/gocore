package gctxt

const (
	LangBrazil  Language = iota
	LangEnglish Language = iota
)

var (
	DefaultLang = LangBrazil
)

var (
	languageTypeContextKey = languageTypeContext{}
)
