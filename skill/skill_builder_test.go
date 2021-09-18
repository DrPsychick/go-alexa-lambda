package skill_test

import (
	"encoding/json"
	"fmt"
	"github.com/drpsychick/go-alexa-lambda/l10n"
	"github.com/drpsychick/go-alexa-lambda/skill"
	"github.com/drpsychick/go-alexa-lambda/ssml"
	"github.com/stretchr/testify/assert"
	"testing"
)

var registry = l10n.NewRegistry()

var enUS = &l10n.Locale{
	Name: "en-US",
	TextSnippets: map[string][]string{
		// Skill
		l10n.KeySkillTestingInstructions: {"Initial instructions"},
		"Skill_Instructions":             {"My instructions"},
		l10n.KeySkillName:                {"SkillName"},
		l10n.KeySkillDescription:         {"SkillDescription"},
		l10n.KeySkillSummary:             {"SkillSummary"},
		l10n.KeySkillKeywords:            {"Keyword1", "Keyword2"},
		l10n.KeySkillExamplePhrases:      {"start me", "boot me up"},
		l10n.KeySkillSmallIconURI:        {"https://small"},
		l10n.KeySkillLargeIconURI:        {"https://large"},
		l10n.KeySkillPrivacyPolicyURL:    {"https://policy"},
		//l10n.KeySkillTermsOfUseURL:       {"https://toc"},
		l10n.KeySkillInvocation: {"call me"},
		"Name":                  {"name"},
		"Summary":               {"summary"},
		"Description":           {"description"},
		"Keywords":              {"key", "words"},
		"Examples":              {"say", "something"},
		"SmallIcon":             {"https://small.icon"},
		"LargeIcon":             {"https://large.icon"},
		"Privacy":               {"https://privacy.url"},
		"Terms":                 {"https://terms.url"},
		// Model
		// Intents
		"MyIntent_Samples":                 {"say one", "say two"},
		"MyIntent_Title":                   {"Title"},
		"MyIntent_Text":                    {"Text1", "Text2"},
		"MyIntent_SSML":                    {ssml.Speak("SSML one"), ssml.Speak("SSML two")},
		"SlotIntent_Samples":               {"what about slot {SlotName}"},
		"SlotIntent_Title":                 {"Test intent with slot"},
		"SlotIntent_Text":                  {"it seems to work"},
		"SlotIntent_SlotName_Samples":      {"of {SlotName}", "{SlotName}"},
		"SlotIntent_SlotName_Elicit_Text":  {"Which slot did you mean?", "I did not understand, which slot?"},
		"SlotIntent_SlotName_Elicit_SSML":  {ssml.Speak("I'm sorry, which slot did you mean?")},
		"SlotIntent_SlotName_Confirm_SSML": {ssml.Speak("Are you sure you know what you're doing?")},
		// Types
		"MyType_Values": {"Value 1", "Value 2"},
	},
}

func init() {
	if err := registry.Register(enUS, l10n.AsDefault()); err != nil {
		panic("something went horribly wrong")
	}
}

func TestSetup(t *testing.T) {
	assert.NotEmpty(t, registry.GetLocales())
	assert.NotEmpty(t, registry.GetDefault())
	assert.Equal(t, enUS, registry.GetDefault())
}

// SkillBuilder Build with registry is covered.
func TestSkillBuilder_WithLocaleRegistry(t *testing.T) {
	sb := skill.NewSkillBuilder().
		WithLocaleRegistry(registry).
		WithCategory(skill.CategoryKnowledgeAndTrivia)

	_, err := sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
}

// SkillBuilder Lookup locale keys is covered.
func TestSkillBuilder_LocaleKeyLookups(t *testing.T) {
	// setup
	sb := skill.NewSkillBuilder().
		WithLocaleRegistry(registry).
		WithCategory(skill.CategoryNews)
	us := enUS.TextSnippets

	// takes default keys
	sk, err := sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
	assert.Equal(t, us[l10n.KeySkillTestingInstructions][0], sk.Manifest.Publishing.TestingInstructions)

	// define our own keys
	sb.WithTestingInstructions("Skill_Instructions")
	sk, err = sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
	assert.Equal(t, us["Skill_Instructions"][0], sk.Manifest.Publishing.TestingInstructions)
}

// SkillBuilder Defining default locale is covered.
func TestSkillBuilder_DefineDefaultLocale(t *testing.T) {
	// setup: skill with first local as default
	sb := skill.NewSkillBuilder().
		WithCategory(skill.CategoryDeviceTracking).
		AddLocale("en-US").
		WithDefaultLocaleTestingInstructions("test it")
	sb.Locale("en-US").
		WithLocaleName("name").
		WithLocaleSummary("summary").
		WithLocaleDescription("description").
		WithLocaleSmallIcon("small").
		WithLocaleLargeIcon("large")
	sk, err := sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
	assert.Equal(t, "test it", sk.Manifest.Publishing.TestingInstructions)

	// add FR as default
	sb.AddLocale("fr-FR", l10n.AsDefault()).
		WithDefaultLocaleTestingInstructions("test le").
		Locale("fr-FR").
		WithLocaleName("nom").
		WithLocaleSummary("sommaire").
		WithLocaleDescription("description").
		WithLocaleSmallIcon("petit").
		WithLocaleLargeIcon("gros")
	sk, err = sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
	assert.Equal(t, "test le", sk.Manifest.Publishing.TestingInstructions)

	// add DE and later make it default
	sb.AddLocale("de-DE").Locale("de-DE").
		WithLocaleName("name").
		WithLocaleSummary("zusammenfassung").
		WithLocaleDescription("beschreibung").
		WithLocaleSmallIcon("klein").
		WithLocaleLargeIcon("groß")
	sk, err = sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
	assert.Equal(t, "test le", sk.Manifest.Publishing.TestingInstructions)

	sb.WithDefaultLocale("de-DE").
		WithDefaultLocaleTestingInstructions("teste es")
	sk, err = sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
	assert.Equal(t, "teste es", sk.Manifest.Publishing.TestingInstructions)
}

// SkillBuilder Defining countries is covered.
func TestSkillBuilder_DefineCountries(t *testing.T) {
	// setup
	sb := skill.NewSkillBuilder().
		WithLocaleRegistry(registry).
		WithCategory(skill.CategoryTvGuides)

	// add
	sb.AddCountry("US").AddCountries([]string{"US", "CA"})
	sk, err := sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
	assert.Equal(t, []string{"US", "US", "CA"}, sk.Manifest.Publishing.Countries)

	// set/overwrite
	sb.WithCountries([]string{"FR", "DE"})
	sk, err = sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
	assert.Equal(t, []string{"FR", "DE"}, sk.Manifest.Publishing.Countries)
}

// SkillBuilder Privacy flags are covered.
func TestSkillBuilder_WithPrivacyFlag(t *testing.T) {
	// setup
	sb := skill.NewSkillBuilder().
		WithLocaleRegistry(registry).
		WithCategory(skill.CategoryUnitConverters).
		WithPrivacyFlag(skill.FlagIsExportCompliant, true).
		WithPrivacyFlag(skill.FlagContainsAds, true).
		WithPrivacyFlag(skill.FlagIsChildDirected, true).
		WithPrivacyFlag(skill.FlagAllowsPurchases, true).
		WithPrivacyFlag(skill.FlagUsesPersonalInfo, true)

	sk, err := sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
	assert.Equal(t, true, sk.Manifest.Privacy.IsExportCompliant)
	assert.Equal(t, true, sk.Manifest.Privacy.ContainsAds)
	assert.Equal(t, true, sk.Manifest.Privacy.IsChildDirected)
	assert.Equal(t, true, sk.Manifest.Privacy.AllowsPurchases)
	assert.Equal(t, true, sk.Manifest.Privacy.UsesPersonalInfo)
}

// SkillBuilder Model is covered.
func TestSkillBuilder_WithModel(t *testing.T) {
	sb := skill.NewSkillBuilder().
		WithLocaleRegistry(registry).
		WithModel()

	ms1, err1 := sb.BuildModels()
	ms2, err2 := sb.Model().Build()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, ms1, ms2)
}

// SkillBuilder Missing elements are covered.
func TestSkillBuilder_ErrorsIfElementsAreMissing(t *testing.T) {
	// setup
	sb := skill.NewSkillBuilder().
		AddLocale("en-US")

	// no category
	_, err := sb.Build()
	assert.Error(t, err)
}

// SkillBuilder To many elements are covered.
func TestSkillBuilder_ErrorsIfTooManyElements(t *testing.T) {
	// setup
	sb := skill.NewSkillBuilder().
		WithCategory(skill.CategoryFriendsAndFamily).
		AddLocale("en-US").
		WithDefaultLocaleTestingInstructions("test it")
	sb.Locale("en-US").
		WithLocaleName("Name").
		WithLocaleDescription("Description").
		WithLocaleSummary("Summary").
		WithLocaleSmallIcon("https://small").
		WithLocaleLargeIcon("https://large")

	// max is 3 example phrases
	sb.Locale("en-US").WithLocaleExamples([]string{"1", "2", "3", "4"})
	_, err := sb.Build()
	assert.Error(t, err)

	// max is 3 keywords
	sb.Locale("en-US").
		WithLocaleExamples([]string{"1", "2", "3"}).
		WithLocaleKeywords([]string{"1", "2", "3", "4"})
	_, err = sb.Build()
	assert.Error(t, err)

	// termsOfUse not allowed (yet)
	sb.Locale("en-US").
		WithLocaleKeywords([]string{"1", "2", "3"}).
		WithTermsURL("MyTermsOfUseURL").
		WithLocaleTermsURL("http://terms")
	_, err = sb.Build()
	assert.Error(t, err)

	// now it builds...
	sb.Locale("en-US").WithLocaleTermsURL("")
	_, err = sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
}

// SkillBuilder Incomplete locale errors are covered.
func TestSkillBuilder_ErrorsIfLocaleIsIncomplete(t *testing.T) {
	sb := skill.NewSkillBuilder().
		WithCategory(skill.CategoryCookingAndRecipe)

	// no locales
	_, err := sb.Build()
	assert.Error(t, err)

	// missing testing instructions translation
	sb.AddLocale("en-US")
	_, err = sb.Build()
	assert.Error(t, err)

	// missing translations
	sb.WithDefaultLocaleTestingInstructions("test it...")
	_, err = sb.Build()
	assert.Error(t, err)
	sb.Locale("en-US").WithLocaleName("my name")
	_, err = sb.Build()
	assert.Error(t, err)
	sb.Locale("en-US").WithLocaleSummary("summary")
	_, err = sb.Build()
	assert.Error(t, err)
	sb.Locale("en-US").WithLocaleDescription("description")
	_, err = sb.Build()
	assert.Error(t, err)
	sb.Locale("en-US").WithLocaleKeywords([]string{"key", "words"})
	_, err = sb.Build()
	assert.Error(t, err)
	sb.Locale("en-US").WithLocaleExamples([]string{"ex", "amples"})
	_, err = sb.Build()
	assert.Error(t, err)
	sb.Locale("en-US").WithLocaleSmallIcon("https://small")
	_, err = sb.Build()
	assert.Error(t, err)

	// last one missing
	sb.Locale("en-US").WithLocaleLargeIcon("https://large")
	_, err = sb.Build()
	assert.NoError(t, err)
	assert.NoError(t, testBuilderImmutability(sb))
}

// SkillBuilder Adding locale twice is covered.
func TestSkillBuilder_ErrorsAddingLocaleTwice(t *testing.T) {
	sb := skill.NewSkillBuilder().
		WithCategory(skill.CategoryKnowledgeAndTrivia).
		AddLocale("en-US").
		WithDefaultLocaleTestingInstructions("test it...").
		AddLocale("en-US")
	assert.IsType(t, &skill.SkillBuilder{}, sb)
	slb := sb.Locale("en-US")
	assert.IsType(t, &skill.SkillLocaleBuilder{}, slb)
	_, err := sb.Build()
	assert.Error(t, err)
}

// SkillBuilder No default locale is covered.
func TestSkillBuilder_ErrorsIfNoDefaultLocale(t *testing.T) {
	sb := skill.NewSkillBuilder().WithDefaultLocale("fr-FR")
	_, err := sb.Build()
	assert.Error(t, err)

	sb = skill.NewSkillBuilder().WithDefaultLocaleTestingInstructions("foo bar")
	_, err = sb.Build()
	assert.Error(t, err)

	sb = skill.NewSkillBuilder()
	slb := sb.Locale("fr-FR")
	_, err = sb.Build()
	assert.Error(t, err)
	assert.Equal(t, &skill.SkillLocaleBuilder{}, slb)
}

// SkillBuilder No model is covered.
func TestSkillBuilder_ErrorsIfNoModel(t *testing.T) {
	sb := skill.NewSkillBuilder()
	_, err := sb.BuildModels()
	assert.Error(t, err)

	sb.Model()
	_, err = sb.BuildModels()
	assert.Error(t, err)
}

// SkillBuilder Incomplete second locale is covered.
func TestSkillBuilder_ErrorsIfSecondLocaleIsIncomplete(t *testing.T) {
	// setup
	deDE := &l10n.Locale{
		Name: "de-DE",
		TextSnippets: map[string][]string{
			l10n.KeySkillTestingInstructions: {"Initial instructions"},
			l10n.KeySkillName:                {"SkillName"},
			l10n.KeySkillDescription:         {"SkillDescription"},
			l10n.KeySkillSummary:             {"SkillSummary"},
			l10n.KeySkillKeywords:            {"Keyword1", "Keyword2"},
			l10n.KeySkillExamplePhrases:      {"start me", "boot me up"},
			//l10n.KeySkillSmallIconURI:        {"https://small"},
			//l10n.KeySkillLargeIconURI:        {"https://large"},
			//l10n.KeySkillPrivacyPolicyURL:    {"https://policy"},
		},
	}
	err := registry.Register(deDE)
	assert.NoError(t, err)
	sb := skill.NewSkillBuilder().
		WithLocaleRegistry(registry).
		WithCategory(skill.CategoryNews)

	// some keys missing for 'de-DE'
	_, err = sb.Build()
	assert.Error(t, err)
}

// SkillLocaleBuilder No locale is covered.
func TestSkillLocaleBuilder_ErrorsIfNoLocale(t *testing.T) {
	// overwrite with empty registry
	l := skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(l10n.NewRegistry())
	_, err := l.BuildPublishingLocale()
	assert.Error(t, err)
	_, err2 := l.BuildPrivacyLocale()
	assert.Error(t, err2)
	assert.Equal(t, err, err2)

	// each setup fails (sets error on SkillLocaleBuilder)
	l = skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(l10n.NewRegistry()).
		WithLocaleName("foo")
	_, err = l.BuildPublishingLocale()
	assert.Error(t, err)

	l = skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(l10n.NewRegistry()).
		WithLocaleSummary("foo")
	_, err2 = l.BuildPublishingLocale()
	assert.Error(t, err2)
	assert.Equal(t, err, err2)

	l = skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(l10n.NewRegistry()).
		WithLocaleDescription("foo")
	_, err2 = l.BuildPublishingLocale()
	assert.Error(t, err2)
	assert.Equal(t, err, err2)

	l = skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(l10n.NewRegistry()).
		WithLocaleExamples([]string{"foo", "bar"})
	_, err2 = l.BuildPublishingLocale()
	assert.Error(t, err2)
	assert.Equal(t, err, err2)

	l = skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(l10n.NewRegistry()).
		WithLocaleKeywords([]string{"foo", "bar"})
	_, err2 = l.BuildPublishingLocale()
	assert.Error(t, err2)
	assert.Equal(t, err, err2)

	l = skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(l10n.NewRegistry()).
		WithLocaleSmallIcon("https://foo")
	_, err2 = l.BuildPublishingLocale()
	assert.Error(t, err2)
	assert.Equal(t, err, err2)

	l = skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(l10n.NewRegistry()).
		WithLocaleLargeIcon("https://bar")
	_, err2 = l.BuildPublishingLocale()
	assert.Error(t, err2)
	assert.Equal(t, err, err2)

	l = skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(l10n.NewRegistry()).
		WithLocalePrivacyURL("https://foo")
	_, err2 = l.BuildPublishingLocale()
	assert.Error(t, err2)
	assert.Equal(t, err, err2)

	l = skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(l10n.NewRegistry()).
		WithLocaleTermsURL("https://bar")
	pl, err := l.BuildPublishingLocale()
	assert.Error(t, err)
	assert.Empty(t, pl)
	assert.Equal(t, err, err2)

	// fails with the same error
	pl2, err2 := l.BuildPrivacyLocale()
	assert.Error(t, err2)
	assert.Empty(t, pl2)
	assert.Equal(t, err, err2)
}

// SkillLocaleBuilder Build with registry is covered.
func TestSkillLocaleBuilder_WithLocaleRegistry(t *testing.T) {
	// setup
	slb := skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(registry)

	// builds without errors
	_, err := slb.BuildPublishingLocale()
	assert.NoError(t, err)
	_, err = slb.BuildPrivacyLocale()
	assert.NoError(t, err)
	assert.NoError(t, testLocalBuilderImmutability(slb))
}

// SkillLocaleBuilder Lookup locale keys is covered.
func TestSkillLocaleBuilder_LocaleKeyLookups(t *testing.T) {
	// setup
	slb := skill.NewSkillLocaleBuilder("en-US").
		WithLocaleRegistry(registry)
	us := enUS.TextSnippets

	// takes default keys
	pbl, err := slb.BuildPublishingLocale()
	assert.NoError(t, err)
	assert.NoError(t, testLocalBuilderImmutability(slb))
	assert.Equal(t, us[l10n.KeySkillName][0], pbl.Name)
	assert.Equal(t, us[l10n.KeySkillSummary][0], pbl.Summary)
	assert.Equal(t, us[l10n.KeySkillDescription][0], pbl.Description)
	assert.Equal(t, us[l10n.KeySkillExamplePhrases][1], pbl.Examples[1])
	assert.Equal(t, us[l10n.KeySkillKeywords][0], pbl.Keywords[0])
	assert.Equal(t, us[l10n.KeySkillSmallIconURI][0], pbl.SmallIconURI)
	assert.Equal(t, us[l10n.KeySkillLargeIconURI][0], pbl.LargeIconURI)

	prl, err := slb.BuildPrivacyLocale()
	assert.NoError(t, err)
	assert.NoError(t, testLocalBuilderImmutability(slb))
	assert.Equal(t, us[l10n.KeySkillPrivacyPolicyURL][0], prl.PrivacyPolicyURL)

	// define our own keys
	slb.WithLocaleRegistry(registry).
		WithName("Name").
		WithSummary("Summary").
		WithDescription("Description").
		WithExamples("Examples").
		WithKeywords("Keywords").
		WithSmallIcon("SmallIcon").
		WithLargeIcon("LargeIcon").
		WithPrivacyURL("Privacy")

	pbl, err = slb.BuildPublishingLocale()
	assert.NoError(t, err)
	assert.NoError(t, testLocalBuilderImmutability(slb))
	assert.Equal(t, us["Name"][0], pbl.Name)
	assert.Equal(t, us["Summary"][0], pbl.Summary)
	assert.Equal(t, us["Description"][0], pbl.Description)
	assert.Equal(t, us["Examples"][1], pbl.Examples[1])
	assert.Equal(t, us["Keywords"][0], pbl.Keywords[0])
	assert.Equal(t, us["SmallIcon"][0], pbl.SmallIconURI)
	assert.Equal(t, us["LargeIcon"][0], pbl.LargeIconURI)

	prl, err = slb.BuildPrivacyLocale()
	assert.NoError(t, err)
	assert.NoError(t, testLocalBuilderImmutability(slb))
	assert.Equal(t, us["Privacy"][0], prl.PrivacyPolicyURL)
}

// SkillLocaleBuilder Set locale translations is covered.
func TestSkillLocaleBuilder_WithLocale(t *testing.T) {
	// setup
	us := enUS.TextSnippets
	// set translations of default keys
	slb := skill.NewSkillLocaleBuilder("en-US").
		WithLocaleName(us["Name"][0]).
		WithLocaleSummary(us["Summary"][0]).
		WithLocaleDescription(us["Description"][0]).
		WithLocaleSmallIcon(us["SmallIcon"][0]).
		WithLocaleLargeIcon(us["LargeIcon"][0]).
		WithLocaleExamples(us["Examples"]).
		WithLocaleKeywords(us["Keywords"]).
		WithLocalePrivacyURL(us["Privacy"][0])

	pbl, err := slb.BuildPublishingLocale()
	assert.NoError(t, err)
	assert.NoError(t, testLocalBuilderImmutability(slb))
	assert.Equal(t, us["Name"][0], pbl.Name)
	assert.Equal(t, us["Summary"][0], pbl.Summary)
	assert.Equal(t, us["Description"][0], pbl.Description)
	assert.Equal(t, us["Examples"][1], pbl.Examples[1])
	assert.Equal(t, us["Keywords"][0], pbl.Keywords[0])
	assert.Equal(t, us["SmallIcon"][0], pbl.SmallIconURI)
	assert.Equal(t, us["LargeIcon"][0], pbl.LargeIconURI)

	prl, err := slb.BuildPrivacyLocale()
	assert.NoError(t, err)
	assert.NoError(t, testLocalBuilderImmutability(slb))
	assert.Equal(t, us["Privacy"][0], prl.PrivacyPolicyURL)
}

// helper to compare two builds
func testBuilderImmutability(s *skill.SkillBuilder) error {
	s1, err := s.Build()
	if err != nil {
		return err
	}
	res1, err := json.MarshalIndent(s1, "", "  ")
	if err != nil {
		return err
	}

	s2, err := s.Build()
	if err != nil {
		return err
	}
	res2, err := json.MarshalIndent(s2, "", "  ")
	if err != nil {
		return err
	}
	if string(res1) != string(res2) {
		return fmt.Errorf("Building skill is not immutable!\n%+v\n%+v", string(res1), string(res2))
	}
	return nil
}

// helper to compare two builds
func testLocalBuilderImmutability(l *skill.SkillLocaleBuilder) error {
	pbl1, err := l.BuildPublishingLocale()
	if err != nil {
		return err
	}
	rpbl1, err := json.MarshalIndent(pbl1, "", "  ")
	if err != nil {
		return err
	}
	prl1, err := l.BuildPrivacyLocale()
	if err != nil {
		return err
	}
	rprl1, err := json.MarshalIndent(prl1, "", "  ")
	if err != nil {
		return err
	}

	pbl2, err := l.BuildPublishingLocale()
	if err != nil {
		return err
	}
	rpbl2, err := json.MarshalIndent(pbl2, "", "  ")
	if err != nil {
		return err
	}
	prl2, err := l.BuildPrivacyLocale()
	if err != nil {
		return err
	}
	rprl2, err := json.MarshalIndent(prl2, "", "  ")

	if string(rpbl1) != string(rpbl2) ||
		string(rprl1) != string(rprl2) {
		return fmt.Errorf("Building locale is not immutable!\n%+v\n%+v", string(rpbl1), string(rpbl2))
	}
	return nil
}
