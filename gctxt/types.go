package gctxt

import "context"

type Language int

type Txt map[Language]string

func (t Txt) Content() string {
	c, exists := t[DefaultLang]
	if exists {
		return c
	}
	c, exists = t[LangBrazil]
	if exists {
		return c
	}
	return ""
}

func (t Txt) ContentByLanguage(lang Language) string {
	c, exists := t[lang]
	if exists {
		return c
	}
	return t.Content()
}

func (t Txt) ContentByContext(ctx context.Context) string {
	return t.ContentByLanguage(ContextValueLanguage(ctx))
}

func (t Txt) MapContent(mapFunc func(Language, string) string) Txt {
	newTxt := make(Txt)
	for k, v := range t {
		newTxt[k] = mapFunc(k, v)
	}
	return newTxt
}

type ListTxt map[Language][]string

func (t ListTxt) Content() []string {
	c, exists := t[DefaultLang]
	if exists {
		return c
	}
	c, exists = t[LangBrazil]
	if exists {
		return c
	}
	return []string{}
}

func (t ListTxt) ContentByLanguage(lang Language) []string {
	c, exists := t[lang]
	if exists {
		return c
	}
	return t.Content()
}

func (t ListTxt) ContentByContext(ctx context.Context) []string {
	return t.ContentByLanguage(ContextValueLanguage(ctx))
}

func (t ListTxt) MapContent(mapFunc func(Language, []string) []string) ListTxt {
	newTxt := make(ListTxt)
	for k, v := range t {
		newTxt[k] = mapFunc(k, v)
	}
	return newTxt
}

type languageTypeContext struct{}
