// Package l10n provides locale helpers (LocaleRegistry) and predefined standard keys for Alexa skills
package l10n

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Default keys.
const (
	KeySkillName                string = "SKILL_Name"
	KeySkillDescription         string = "SKILL_Description"
	KeySkillSummary             string = "SKILL_Summary"
	KeySkillExamplePhrases      string = "SKILL_ExamplePhrases"
	KeySkillKeywords            string = "SKILL_Keywords"
	KeySkillSmallIconURI        string = "SKILL_SmallIconURI"
	KeySkillLargeIconURI        string = "SKILL_LargeIconURI"
	KeySkillTestingInstructions string = "SKILL_TestingInstructions"
	KeySkillInvocation          string = "SKILL_Invocation"
	KeySkillPrivacyPolicyURL    string = "SKILL_PrivacyPolicyURL"
	KeySkillTermsOfUseURL       string = "SKILL_TermsOfUse"
	KeyPostfixSamples           string = "_Samples"
	KeyPostfixValues            string = "_Values"
	KeyPostfixTitle             string = "_Title"
	KeyPostfixText              string = "_Text"
	KeyPostfixSSML              string = "_SSML"
	// fallback, standard error.
	KeyErrorTitle string = "Error_Title"
	KeyErrorText  string = "Error_Text"
	KeyErrorSSML  string = "Error_SSML"
	// unknown error.
	KeyErrorUnknownTitle string = "Error_Unknown_Title"
	KeyErrorUnknownText  string = "Error_Unknown_Text"
	KeyErrorUnknownSSML  string = "Error_Unknown_SSML"
	// not found error (slot, intent, ...), see request.go.
	KeyErrorNotFoundTitle string = "Error_NotFound_Title"
	KeyErrorNotFoundText  string = "Error_NotFound_Text"
	KeyErrorNotFoundSSML  string = "Error_NotFound_SSML"
	// missing locale error.
	KeyErrorLocaleNotFoundTitle string = "Error_LocaleNotFound_Title"
	KeyErrorLocaleNotFoundText  string = "Error_LocaleNotFound_Text"
	KeyErrorLocaleNotFoundSSML  string = "Error_LocaleNotFound_SSML"
	// generic translation error.
	KeyErrorTranslationTitle string = "Error_Translation_Title"
	KeyErrorTranslationText  string = "Error_Translation_Text"
	KeyErrorTranslationSSML  string = "Error_Translation_SSML"
	// missing translation error.
	KeyErrorNoTranslationTitle string = "Error_NoTranslation_Title"
	KeyErrorNoTranslationText  string = "Error_NoTranslation_Text"
	KeyErrorNoTranslationSSML  string = "Error_NoTranslation_SSML"
	// missing placeholder in translation error.
	KeyErrorMissingPlaceholderTitle string = "Error_MissingPlaceholder_Title"
	KeyErrorMissingPlaceholderText  string = "Error_MissingPlaceholder_Text"
	KeyErrorMissingPlaceholderSSML  string = "Error_MissingPlaceholder_SSML"
	// default intents.
	KeyLaunchTitle string = "Launch_Title"
	KeyLaunchText  string = "Launch_Text"
	KeyLaunchSSML  string = "Launch_SSML"
	KeyHelpTitle   string = "Help_Title"
	KeyHelpText    string = "Help_Text"
	KeyHelpSSML    string = "Help_SSML"
	KeyStopTitle   string = "Stop_Title"
	KeyStopText    string = "Stop_Text"
	KeyStopSSML    string = "Stop_SSML"
	KeyCancelTitle string = "Cancel_Title"
	KeyCancelText  string = "Cancel_Text"
	KeyCancelSSML  string = "Cancel_SSML"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// LocaleError defines the interface for locale errors.
type LocaleError interface {
	error
	GetLocale() string
	GetKey() string
	GetPlaceholder() string
}

// NoTranslationError defines a missing translation error.
type NoTranslationError struct {
	Locale      string
	Key         string
	Placeholder string
}

// Error returns a string of the error.
func (n NoTranslationError) Error() string {
	return fmt.Sprintf("locale %s: no translation for key '%s'", n.Locale, n.Key)
}

// GetLocale returns the locale of the error.
func (n NoTranslationError) GetLocale() string {
	return n.Locale
}

// GetKey returns the associated key.
func (n NoTranslationError) GetKey() string {
	return n.Key
}

// GetPlaceholder returns the placeholder concerned.
func (n NoTranslationError) GetPlaceholder() string {
	return n.Placeholder
}

// MissingPlaceholderError defines a generic not found error.
type MissingPlaceholderError struct {
	Locale      string
	Key         string
	Placeholder string
}

// GetLocale returns locale of the error.
func (e MissingPlaceholderError) GetLocale() string {
	return e.Locale
}

// GetKey returns the key of the error.
func (e MissingPlaceholderError) GetKey() string {
	return e.Key
}

// GetPlaceholder returns the placeholder concerned.
func (e MissingPlaceholderError) GetPlaceholder() string {
	return e.Placeholder
}

// Error returns a string representing the error including the key (and placeholder) missing.
func (e MissingPlaceholderError) Error() string {
	if e.Placeholder == "" {
		return fmt.Sprintf("locale %s: key '%s' is missing a placeholder in translation", e.Locale, e.Key)
	}
	return fmt.Sprintf("locale %s: key '%s' is missing placeholder '%s' in translation", e.Locale, e.Key, e.Placeholder)
}

// LocaleRegistry is the interface for an l10n registry.
type LocaleRegistry interface {
	Register(locale LocaleInstance, opts ...RegisterFunc) error
	Resolve(locale string) (LocaleInstance, error)
	GetDefault() LocaleInstance
	SetDefault(locale string) error
	GetLocales() map[string]LocaleInstance
}

// LocaleInstance is the interface for a specific locale.
type LocaleInstance interface {
	GetName() string
	Set(key string, values []string)
	Get(key string, args ...interface{}) string
	GetAny(key string, args ...interface{}) string
	GetAll(key string, args ...interface{}) []string
	GetErrors() []error
	ResetErrors()
}

// DefaultRegistry is the standard registry used.
var DefaultRegistry = NewRegistry()

// Config contains the options for Locale registration.
type Config struct {
	DefaultLocale bool
	FallbackFor   string
}

// RegisterFunc defines the functions to be passed to Register.
type RegisterFunc func(cfg *Config)

// AsDefault registers the given Locale as the default.
func AsDefault() RegisterFunc {
	return func(cfg *Config) {
		cfg.DefaultLocale = true
	}
}

// Registry is the Locale registry.
type Registry struct {
	defaultLocale string
	locales       map[string]LocaleInstance
}

// NewRegistry returns an empty Registry.
func NewRegistry() LocaleRegistry {
	return &Registry{locales: map[string]LocaleInstance{}}
}

// Register registers a new Locale in the DefaultRegistry.
func Register(locale LocaleInstance, opts ...RegisterFunc) error {
	return DefaultRegistry.Register(locale, opts...)
}

// GetDefault returns the default locale in the DefaultRegistry.
func GetDefault() LocaleInstance {
	return DefaultRegistry.GetDefault()
}

// SetDefault sets the default locale in the DefaultRegistry.
func SetDefault(locale string) error {
	return DefaultRegistry.SetDefault(locale)
}

// GetLocales returns the locales registered in the DefaultRegistry.
func GetLocales() map[string]LocaleInstance {
	return DefaultRegistry.GetLocales()
}

// Resolve returns the matching locale from the DefaultRegistry.
func Resolve(name string) (LocaleInstance, error) {
	return DefaultRegistry.Resolve(name)
}

// Register registers a new locale and fails if it already exists.
func (r *Registry) Register(l LocaleInstance, opts ...RegisterFunc) error {
	if l.GetName() == "" {
		return fmt.Errorf("cannot register locale with no name")
	}
	if _, ok := r.locales[l.GetName()]; ok {
		return fmt.Errorf("locale %s already registered", l.GetName())
	}

	// run all RegisterFuncs
	var cfg Config
	for _, opt := range opts {
		opt(&cfg)
	}

	// set locale as default
	if cfg.DefaultLocale || r.defaultLocale == "" {
		r.defaultLocale = l.GetName()
	}

	r.locales[l.GetName()] = l

	return nil
}

// GetDefault returns the default locale.
func (r *Registry) GetDefault() LocaleInstance {
	return r.locales[r.defaultLocale]
}

// SetDefault sets the default locale which must be registered.
func (r *Registry) SetDefault(locale string) error {
	if _, err := r.Resolve(locale); err != nil {
		return err
	}
	r.defaultLocale = locale
	return nil
}

// GetLocales returns all registered locales.
func (r *Registry) GetLocales() map[string]LocaleInstance {
	return r.locales
}

// Resolve returns the Locale matching the given name or an error.
func (r *Registry) Resolve(locale string) (LocaleInstance, error) {
	l, ok := r.locales[locale]
	if !ok {
		return nil, fmt.Errorf("locale '%s' not found", locale)
	}
	return l, nil
}

// Locale is a representation of keys in a specific language.
type Locale struct {
	Name         string // de-DE, en-US, ...
	TextSnippets Snippets
	errors       []error
}

// NewLocale creates a new, empty locale.
func NewLocale(locale string) *Locale {
	return &Locale{
		Name:         locale,
		TextSnippets: Snippets{},
	}
}

// GetName returns the name of the locale.
func (l *Locale) GetName() string {
	return l.Name
}

// Set sets the translations for a key.
func (l *Locale) Set(key string, values []string) {
	l.TextSnippets[key] = values
}

// Get returns the first translation.
func (l *Locale) Get(key string, args ...interface{}) string {
	t, err := l.TextSnippets.GetFirst(key, args...)
	if err != nil {
		var locaerr NoTranslationError
		switch {
		case errors.As(err, &locaerr):
			locaerr.Locale = l.GetName()
			l.errors = append(l.errors, locaerr)
		default:
			l.errors = append(l.errors, err)
		}
	}
	l.appendErrorMissingParam(key, []string{t})
	return t
}

// GetAny returns a random translation.
func (l *Locale) GetAny(key string, args ...interface{}) string {
	t, err := l.TextSnippets.GetAny(key, args...)
	if err != nil {
		var locaErr NoTranslationError
		switch {
		case errors.As(err, &locaErr):
			locaErr.Locale = l.GetName()
			l.errors = append(l.errors, locaErr)
		default:
			l.errors = append(l.errors, err)
		}
	}
	l.appendErrorMissingParam(key, []string{t})
	return t
}

// GetAll returns all translations.
func (l *Locale) GetAll(key string, args ...interface{}) []string {
	t, err := l.TextSnippets.GetAll(key, args...)
	if err != nil {
		var locaerr NoTranslationError
		switch {
		case errors.As(err, &locaerr):
			locaerr.Locale = l.GetName()
			l.errors = append(l.errors, locaerr)
		default:
			l.errors = append(l.errors, err)
		}
	}
	l.appendErrorMissingParam(key, t)
	return t
}

// GetErrors returns key lookup errors that occurred.
func (l *Locale) GetErrors() []error {
	return l.errors
}

// ResetErrors resets existing errors.
func (l *Locale) ResetErrors() {
	l.errors = nil
}

func (l *Locale) appendErrorMissingParam(key string, texts []string) {
	for _, t := range texts {
		if strings.Contains(t, "%!") &&
			strings.Contains(t, "(MISSING)") {
			l.errors = append(l.errors, MissingPlaceholderError{l.GetName(), key, ""})
		}
	}
}

// Snippets is the actual representation of key -> array of translations in a locale.
type Snippets map[string][]string

// GetFirst returns the first translation for the snippet.
func (s Snippets) GetFirst(key string, args ...interface{}) (string, error) {
	_, ok := s[key]
	if !ok || len(s[key]) == 0 {
		return "", NoTranslationError{"", key, ""}
	}
	return fmt.Sprintf(s[key][0], args...), nil
}

// GetAny returns a random translation for the snippet.
func (s Snippets) GetAny(key string, args ...interface{}) (string, error) {
	_, ok := s[key]
	if !ok || len(s[key]) == 0 {
		return "", NoTranslationError{"", key, ""}
	}
	if len(s[key]) == 1 {
		return fmt.Sprintf(s[key][0], args...), nil
	}
	l := len(s[key])
	r := rand.Intn(l) //nolint:gosec
	return fmt.Sprintf(s[key][r], args...), nil
}

// GetAll returns all translations of the snippet.
func (s Snippets) GetAll(key string, args ...interface{}) ([]string, error) {
	_, ok := s[key]
	if !ok || len(s[key]) == 0 {
		return []string{}, NoTranslationError{"", key, ""}
	}
	r := make([]string, len(s[key]))
	for i, v := range s[key] {
		r[i] = fmt.Sprintf(v, args...)
	}
	return r, nil
}
