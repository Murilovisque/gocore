package gctxt

import (
	"context"
	"maps"
)

func Build(brazilMsg string, anotherLanguages Txt) Txt {
	msg := Txt{
		LangBrazil: brazilMsg,
	}
	if anotherLanguages != nil {
		maps.Copy(msg, anotherLanguages)
	}
	return msg
}

func BuildList(brazilMsgs []string, anotherLanguages ListTxt) ListTxt {
	msgs := ListTxt{
		LangBrazil: brazilMsgs,
	}
	if anotherLanguages != nil {
		maps.Copy(msgs, anotherLanguages)
	}
	return msgs
}

func ContextWithLanguage(ctx context.Context, lang Language) context.Context {
	return context.WithValue(ctx, languageTypeContextKey, lang)
}

func ContextValueLanguage(ctx context.Context) Language {
	vl, ok := ctx.Value(languageTypeContextKey).(Language)
	if ok {
		return vl
	}
	return DefaultLang
}
