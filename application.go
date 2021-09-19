// Package alexa enables parsing request and building responses.
package alexa

import (
	"errors"
	"fmt"

	"github.com/drpsychick/go-alexa-lambda/l10n"
	log "github.com/hamba/logger/v2"
	stats "github.com/hamba/statter"
)

// Application defines the interface used of the app.
type Application interface {
	Logger() log.Logger
	Statter() stats.Statter
}

// Response wraps the data needed for a skill response.
type Response struct {
	Title    string
	Text     string
	Speech   string
	Image    string
	Reprompt bool
	End      bool
}

// GetLocaleWithFallback falls back to default locale which must be considered carefully.
func GetLocaleWithFallback(registry l10n.LocaleRegistry, locale string) (l10n.LocaleInstance, Response) {
	loc, err := registry.Resolve(locale)
	if err != nil {
		loc = registry.GetDefault()
		if loc == nil {
			return nil, Response{
				Title: "Error",
				Text:  "No locale found!",
				End:   true,
			}
		}
	}
	return loc, Response{}
}

// ResponseError defines a response error.
type ResponseError interface {
	error

	Response(loc l10n.LocaleInstance) Response
}

// HandleError handles default and ResponseErrors. Returns true if the error was handled.
func HandleError(b *ResponseBuilder, loc l10n.LocaleInstance, err error) bool {
	var resp Response
	var respErr ResponseError
	var l10nErr l10n.LocaleError
	switch {
	case loc == nil:
		resp = Response{
			Title: "Error",
			Text:  "Locale not found!",
			End:   true,
		}
	case errors.As(err, &respErr):
		resp = respErr.Response(loc)
	case errors.As(err, &l10nErr):
		if l10nErr.GetPlaceholder() == "" {
			resp = Response{
				Title:  loc.GetAny(l10n.KeyErrorNoTranslationTitle),
				Text:   loc.GetAny(l10n.KeyErrorNoTranslationText, l10nErr.GetKey()),
				Speech: loc.GetAny(l10n.KeyErrorNoTranslationSSML, l10nErr.GetKey()),
				End:    true,
			}
		} else {
			resp = Response{
				Title:  loc.GetAny(l10n.KeyErrorMissingPlaceholderTitle),
				Text:   loc.GetAny(l10n.KeyErrorMissingPlaceholderText, l10nErr.GetPlaceholder()),
				Speech: loc.GetAny(l10n.KeyErrorMissingPlaceholderSSML, l10nErr.GetPlaceholder()),
				End:    true,
			}
		}
	default:
		return false
	}

	// locale errors during error processing.
	// if err := CheckForLocaleError(loc); err != nil {
	//
	// }

	b.With(resp)
	return true
}

// CheckForLocaleError returns a ResponseError for the last locale error.
func CheckForLocaleError(loc l10n.LocaleInstance) error {
	errs := loc.GetErrors()
	if len(errs) == 0 {
		return nil
	}

	lastErr := errs[len(errs)-1] //nolint:ifshort,nolintlint
	var l10nErr l10n.LocaleError
	if errors.As(lastErr, &l10nErr) {
		return TranslationError{l10nErr.GetLocale(), l10nErr.GetKey()}
	}
	return TextError{loc.GetName(), lastErr.Error()}
}

// TextError returns the error text to Alexa.
type TextError struct {
	Locale string
	Text   string
}

// Error returns the text of the error.
func (e TextError) Error() string {
	return e.Text
}

// Response returns and ApplicationResponse including the error text.
func (e TextError) Response(loc l10n.LocaleInstance) Response {
	return Response{
		Title:  loc.GetAny(l10n.KeyErrorTitle),
		Text:   e.Error(),
		Speech: loc.GetAny(l10n.KeyErrorSSML),
		End:    true,
	}
}

// TranslationError defines a missing translation error.
type TranslationError struct {
	Locale string
	Key    string
}

// Error returns a string representing the error including the key missing.
func (e TranslationError) Error() string {
	return fmt.Sprintf("locale %s: translation for key '%s' is missing", e.Locale, e.Key)
}

// Response returns a Response for the user.
func (e TranslationError) Response(loc l10n.LocaleInstance) Response {
	return Response{
		Title:  loc.GetAny(l10n.KeyErrorTranslationTitle),
		Text:   loc.GetAny(l10n.KeyErrorTranslationText),
		Speech: loc.GetAny(l10n.KeyErrorTranslationSSML),
		End:    true,
	}
}
