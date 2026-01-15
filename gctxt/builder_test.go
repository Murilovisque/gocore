package gctxt

import (
	"context"
	"slices"
	"testing"
)

func TestTxt(t *testing.T) {
	DefaultLang = LangBrazil
	expectedBrazilMsg := "Mensagem"
	expectedEnglishMsg := "Message"

	// Build test
	txt := Build(expectedBrazilMsg, Txt{LangEnglish: expectedEnglishMsg})
	if txt.Content() != expectedBrazilMsg {
		t.Fatalf("expected '%s', but '%s'", expectedBrazilMsg, txt.Content())
	}
	if txt.ContentByLanguage(LangEnglish) != expectedEnglishMsg {
		t.Fatalf("expected '%s', but '%s'", expectedEnglishMsg, txt.ContentByLanguage(LangEnglish))
	}

	// MapContent
	txt = txt.MapContent(func(language Language, s string) string {
		switch language {
		case LangBrazil:
			return "Essa " + s
		case LangEnglish:
			return "This " + s
		default:
			return s
		}
	})
	expectedBrazilMsg = "Essa " + expectedBrazilMsg
	expectedEnglishMsg = "This " + expectedEnglishMsg
	if txt.Content() != expectedBrazilMsg {
		t.Fatalf("expected '%s', but '%s'", expectedBrazilMsg, txt.Content())
	}
	if txt.ContentByLanguage(LangEnglish) != expectedEnglishMsg {
		t.Fatalf("expected '%s', but '%s'", expectedEnglishMsg, txt.ContentByLanguage(LangEnglish))
	}

	// Context
	ctx := context.Background()
	if ContextValueLanguage(ctx) != DefaultLang {
		t.Fatalf("expected '%v', but '%v'", DefaultLang, ContextValueLanguage(ctx))
	}
	ctx = ContextWithLanguage(ctx, LangEnglish)
	if ContextValueLanguage(ctx) != LangEnglish {
		t.Fatalf("expected '%v', but '%v'", LangEnglish, ContextValueLanguage(ctx))
	}
	if txt.ContentByContext(ctx) != expectedEnglishMsg {
		t.Fatalf("expected '%v', but '%v'", expectedEnglishMsg, ContextValueLanguage(ctx))
	}

	DefaultLang = LangEnglish
	if txt.Content() != expectedEnglishMsg {
		t.Fatalf("expected '%s', but '%s'", expectedEnglishMsg, txt.Content())
	}
}

func TestListTxt(t *testing.T) {
	DefaultLang = LangBrazil
	expectedBrazilMsg := []string{"Mensagem"}
	expectedEnglishMsg := []string{"Message"}

	// Build test
	txt := BuildList(expectedBrazilMsg, ListTxt{LangEnglish: expectedEnglishMsg})
	if !slices.Equal(txt.Content(), expectedBrazilMsg) {
		t.Fatalf("expected '%v', but '%v'", expectedBrazilMsg, txt.Content())
	}
	if !slices.Equal(txt.ContentByLanguage(LangEnglish), expectedEnglishMsg) {
		t.Fatalf("expected '%v', but '%v'", expectedEnglishMsg, txt.ContentByLanguage(LangEnglish))
	}

	// MapContent
	txt = txt.MapContent(func(language Language, s []string) []string {
		switch language {
		case LangBrazil:
			return []string{"Essa " + s[0]}
		case LangEnglish:
			return []string{"This " + s[0]}
		default:
			return s
		}
	})
	expectedBrazilMsg = []string{"Essa " + expectedBrazilMsg[0]}
	expectedEnglishMsg = []string{"This " + expectedEnglishMsg[0]}
	if !slices.Equal(txt.Content(), expectedBrazilMsg) {
		t.Fatalf("expected '%v', but '%v'", expectedBrazilMsg, txt.Content())
	}
	if !slices.Equal(txt.ContentByLanguage(LangEnglish), expectedEnglishMsg) {
		t.Fatalf("expected '%v', but '%v'", expectedEnglishMsg, txt.ContentByLanguage(LangEnglish))
	}

	// Context
	ctx := context.Background()
	if ContextValueLanguage(ctx) != DefaultLang {
		t.Fatalf("expected '%v', but '%v'", DefaultLang, ContextValueLanguage(ctx))
	}
	ctx = ContextWithLanguage(ctx, LangEnglish)
	if ContextValueLanguage(ctx) != LangEnglish {
		t.Fatalf("expected '%v', but '%v'", LangEnglish, ContextValueLanguage(ctx))
	}
	if !slices.Equal(txt.ContentByContext(ctx), expectedEnglishMsg) {
		t.Fatalf("expected '%v', but '%v'", expectedEnglishMsg, ContextValueLanguage(ctx))
	}

	DefaultLang = LangEnglish
	if !slices.Equal(txt.Content(), expectedEnglishMsg) {
		t.Fatalf("expected '%v', but '%v'", expectedEnglishMsg, txt.Content())
	}
}
