package l10n_test

import (
	"bou.ke/monkey"
	"github.com/drpsychick/go-alexa-lambda/l10n"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

const Greeting string = "greeting"
const WithParam string = "withparam"

var registry = l10n.NewRegistry()
var deDE = &l10n.Locale{
	Name: "de-DE",
	TextSnippets: l10n.Snippets{
		Greeting: []string{
			"Hi",
			"Hallo",
			"Howdi",
		},
		WithParam: []string{
			"Hello %s",
		},
	},
}
var enUS = &l10n.Locale{
	Name: "en-US",
	TextSnippets: l10n.Snippets{
		Greeting: []string{
			"Hi",
			"Hello",
		},
	},
}

// initialize registry
func init() {
	if err := registry.Register(deDE); err != nil {
		panic("register 'deDE' failed")
	}
	if err := registry.Register(enUS, l10n.AsDefault()); err != nil {
		panic("register 'enUS' failed")
	}
}

// Registry create new is covered.
func TestNewRegistry(t *testing.T) {
	r := l10n.NewRegistry()
	assert.Empty(t, r.GetLocales())
	assert.Nil(t, r.GetDefault())
}

// Registry define default locale is covered.
func TestRegistry_DefineDefault(t *testing.T) {
	// setup
	l10n.DefaultRegistry = l10n.NewRegistry()

	// fails as no locales are registered
	err := l10n.SetDefault("en-US")
	assert.Error(t, err)

	// first locale is default
	en := l10n.NewLocale("en-US")
	err = l10n.Register(en)
	assert.NoError(t, err)
	assert.Equal(t, en, l10n.GetDefault())

	// set default
	de := l10n.NewLocale("de-DE")
	err = l10n.Register(de)
	assert.NoError(t, err)
	err = l10n.SetDefault("de-DE")
	assert.NoError(t, err)
	assert.Equal(t, de, l10n.GetDefault())

	// register as default
	fr := l10n.NewLocale("fr-FR")
	err = l10n.Register(fr, l10n.AsDefault())
	assert.NoError(t, err)
	assert.Equal(t, fr, l10n.GetDefault())
}

// Registry register locale is covered.
func TestRegistry_Register(t *testing.T) {
	// setup
	loc := &l10n.Locale{}

	// fails: no name
	err := l10n.Register(loc)
	assert.Error(t, err)

	loc.Name = "my locale"
	err = l10n.Register(loc)
	assert.NoError(t, err)

	// fails: cannot register twice
	err = l10n.Register(loc)
	assert.Error(t, err)

	// register as default
	loc2 := l10n.NewLocale("foo")
	err = l10n.Register(loc2, l10n.AsDefault())
	assert.NoError(t, err)
	assert.Equal(t, loc2, l10n.GetDefault())
}

// Registry get locales is covered.
func TestRegistry_GetLocales(t *testing.T) {
	// setup
	l10n.DefaultRegistry = l10n.NewRegistry()
	r := l10n.NewRegistry()
	ls := l10n.GetLocales()
	assert.Empty(t, ls)
	ls = r.GetLocales()
	assert.Empty(t, ls)

	err := l10n.Register(l10n.NewLocale("foo"))
	assert.NoError(t, err)
	assert.Len(t, l10n.GetLocales(), 1)

	err = r.Register(l10n.NewLocale("bar"))
	assert.NoError(t, err)
	err = r.Register(l10n.NewLocale("foo"))
	assert.NoError(t, err)
	assert.Len(t, r.GetLocales(), 2)

	assert.Equal(t, l10n.DefaultRegistry.GetLocales(), l10n.GetLocales())
}

// Registry resolve locale is covered.
func TestRegistry_Resolve(t *testing.T) {
	// setup
	l10n.DefaultRegistry = l10n.NewRegistry()
	r := l10n.NewRegistry()

	f := l10n.NewLocale("foo")
	err := l10n.Register(f)
	assert.NoError(t, err)
	f2, err := l10n.Resolve("foo")
	assert.NoError(t, err)
	assert.Equal(t, f, f2)

	err = r.Register(f)
	assert.NoError(t, err)
	f2, err = r.Resolve("foo")
	assert.NoError(t, err)
	assert.Equal(t, f, f2)

	f2, err = l10n.Resolve("bar")
	assert.Error(t, err)
	assert.Equal(t, "locale 'bar' not found", err.Error())
	assert.Nil(t, f2)
}

// Registry no default is covered.
func TestRegistry_ErrorsIfNoDefault(t *testing.T) {
	r := l10n.NewRegistry()
	l := r.GetDefault()
	assert.Nil(t, l)
}

// Locale create new is covered.
func TestNewLocale(t *testing.T) {
	l := l10n.NewLocale("fo-BA")
	assert.IsType(t, &l10n.Locale{}, l)
	assert.Equal(t, "fo-BA", l.GetName())
	assert.Empty(t, l.Get("not exists"))
}

// Locale set is covered.
func TestLocale_Set(t *testing.T) {
	// setup
	l := l10n.NewLocale("fo-BA")
	vals := []string{"bar1", "bar2"}

	l.Set("foo", vals)
	v2 := l.GetAll("foo")
	assert.Equal(t, vals, v2)
}

// Locale get random key is covered.
func TestLocale_GetAny(t *testing.T) {
	// requires registry setup
	assert.NotNil(t, registry)

	patch := monkey.Patch(rand.Intn, func(i int) int {
		return 1
	})
	defer patch.Unpatch()

	l, err := registry.Resolve("de-DE")
	assert.NoError(t, err)
	assert.Contains(t, deDE.TextSnippets[Greeting], l.GetAny(Greeting))
}

// Locale get with param is covered.
func TestLocale_GetWithParam(t *testing.T) {
	// requires registry setup
	assert.NotNil(t, registry)

	l, err := registry.Resolve("de-DE")
	assert.NoError(t, err)

	assert.Equal(t, "Hello there", l.Get(WithParam, "there"))
	assert.Equal(t, "Hello there", l.GetAny(WithParam, "there"))
	assert.Equal(t, "Hello there", l.GetAll(WithParam, "there")[0])
}

// Locale key does not exist is covered.
func TestLocale_KeyNotExists(t *testing.T) {
	// setup
	l10n.DefaultRegistry = l10n.NewRegistry()
	err := l10n.Register(l10n.NewLocale("de-DE"))
	assert.NoError(t, err)
	l, err := l10n.Resolve("de-DE")
	assert.NoError(t, err)

	tx := l.Get("not exists")
	assert.Empty(t, tx)
	tx = l.GetAny("not exists")
	assert.Empty(t, tx)
	txs := l.GetAll("not exists")
	assert.Empty(t, txs)
	assert.Len(t, l.GetErrors(), 3)
	assert.Equal(t, "locale de-DE: no translation for key 'not exists'", l.GetErrors()[1].Error())
}

// Locale empty locale is covered.
func TestLocale_Errors(t *testing.T) {
	l := &l10n.Locale{}
	assert.Empty(t, l.GetName())
	assert.Empty(t, l.Get("foo"))
	assert.Len(t, l.GetErrors(), 1)
}

// Locale with no param
func TestLocale_ErrorNoParam(t *testing.T) {
	// requires registry setup
	assert.NotNil(t, registry)
	l, err := registry.Resolve("de-DE")
	assert.NoError(t, err)

	assert.Equal(t, "Hello %!s(MISSING)", l.Get(WithParam))
	assert.Len(t, l.GetErrors(), 1)
	assert.NotEmpty(t, l.GetErrors())
	assert.Equal(t, "locale de-DE: key '"+WithParam+"' is missing a placeholder in translation", l.GetErrors()[0].Error())
}
