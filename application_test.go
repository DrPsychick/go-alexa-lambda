package alexa

import (
	"errors"
	"fmt"
	"github.com/drpsychick/alexa-go-lambda/l10n"
	"reflect"
	"testing"
)

func initL10n() {
	en := &l10n.Locale{Name: "en-US", TextSnippets: l10n.Snippets{"key": {"value1", "value2"}}}
	de := &l10n.Locale{Name: "de-DE", TextSnippets: l10n.Snippets{"foo": {"bar", "foo"}}}
	l10n.DefaultRegistry.Register(en) // default
	l10n.DefaultRegistry.Register(de)
}

func TestCheckForLocaleError(t *testing.T) {
	type args struct {
		loc l10n.LocaleInstance
	}
	initL10n()
	en, _ := l10n.DefaultRegistry.Resolve("en-US")
	de, _ := l10n.DefaultRegistry.Resolve("de-DE")
	de.Get("key") // causes error
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "LocaleNoError", args: args{en}, wantErr: false},
		{name: "LocaleError", args: args{de}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckForLocaleError(tt.args.loc); (err != nil) != tt.wantErr {
				t.Errorf("CheckForLocaleError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetLocaleWithFallback(t *testing.T) {
	type args struct {
		registry l10n.LocaleRegistry
		locale   string
	}
	en, _ := l10n.DefaultRegistry.Resolve("en-US")
	tests := []struct {
		name  string
		args  args
		want  l10n.LocaleInstance
		want1 Response
	}{
		{"Fallback", args{l10n.DefaultRegistry, "fr-FR"}, en, Response{}},
		{"NoFallback", args{l10n.NewRegistry(), "fr-FR"}, nil, Response{
			Title: "Error", Text: "No locale found!", End: true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetLocaleWithFallback(tt.args.registry, tt.args.locale)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLocaleWithFallback() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetLocaleWithFallback() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHandleError(t *testing.T) {
	type args struct {
		b   *ResponseBuilder
		loc l10n.LocaleInstance
		err error
	}

	initL10n()
	b := &ResponseBuilder{}
	loc, _ := l10n.DefaultRegistry.Resolve("en-US")
	myTextErr := TextError{"en-US", "Text Error"}
	myTransErr := TranslationError{"en-US", "foo"}
	myNoTransErr := l10n.NoTranslationError{Locale: "en-US", Key: "foo", Placeholder: ""}
	myNoTransPlceholderErr := l10n.NoTranslationError{Locale: "en-US", Key: "key", Placeholder: "placeholder"}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"HandleNoError", args{b, loc, nil}, false},
		{"HandleAnyError", args{b, loc, errors.New("foo")}, false},
		{"HandleNilLocale", args{b, nil, nil}, true},
		{"HandleTextError", args{b, loc, myTextErr}, true},
		{"HandleTransError", args{b, loc, myTransErr}, true},
		{"HandleNoTransError", args{b, loc, myNoTransErr}, true},
		{"HandlePlaceholderError", args{b, loc, myNoTransPlceholderErr}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandleError(tt.args.b, tt.args.loc, tt.args.err); got != tt.want {
				t.Errorf("HandleError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextError_Error(t *testing.T) {
	type fields struct {
		Locale string
		Text   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"TextError", fields{"en-US", "text"}, "text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := TextError{
				Locale: tt.fields.Locale,
				Text:   tt.fields.Text,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextError_Response(t *testing.T) {
	type fields struct {
		Locale string
		Text   string
	}
	type args struct {
		loc l10n.LocaleInstance
	}
	initL10n()
	loc, _ := l10n.DefaultRegistry.Resolve("en-US")
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Response
	}{
		{"TextErrorResponse", fields{"en-US", "text"}, args{loc}, Response{
			Title:  loc.GetAny(l10n.KeyErrorTitle),
			Text:   "text",
			Speech: loc.GetAny(l10n.KeyErrorSSML),
			End:    true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := TextError{
				Locale: tt.fields.Locale,
				Text:   tt.fields.Text,
			}
			if got := e.Response(tt.args.loc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslationError_Error(t *testing.T) {
	type fields struct {
		Locale string
		Key    string
	}
	key := "key"
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"TransError", fields{"en-US", key}, fmt.Sprintf(
			"locale %s: translation for key '%s' is missing", "en-US", key,
		)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := TranslationError{
				Locale: tt.fields.Locale,
				Key:    tt.fields.Key,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslationError_Response(t *testing.T) {
	type fields struct {
		Locale string
		Key    string
	}
	type args struct {
		loc l10n.LocaleInstance
	}
	initL10n()
	en, _ := l10n.DefaultRegistry.Resolve("en-US")
	key := "key"
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Response
	}{
		{"TransErrorResponse", fields{"en-US", key}, args{en}, Response{
			Title:  en.GetAny(l10n.KeyErrorTranslationTitle),
			Text:   en.GetAny(l10n.KeyErrorTranslationText),
			Speech: en.GetAny(l10n.KeyErrorTranslationSSML),
			End:    true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := TranslationError{
				Locale: tt.fields.Locale,
				Key:    tt.fields.Key,
			}
			if got := e.Response(tt.args.loc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response() = %v, want %v", got, tt.want)
			}
		})
	}
}
