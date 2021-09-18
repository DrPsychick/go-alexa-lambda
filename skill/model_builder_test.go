package skill_test

import (
	"github.com/drpsychick/alexa-go-lambda/l10n"
	"github.com/drpsychick/alexa-go-lambda/skill"
	"github.com/stretchr/testify/assert"
	"testing"
)

// modelBuilder with invocation is covered.
func TestModelBuilder_WithInvocation(t *testing.T) {
	const InvocationKey string = "MyInvocationKey"
	assert.NotNil(t, registry)
	en, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	mb := skill.NewModelBuilder().
		WithLocaleRegistry(registry)

	// uses default invocation key
	ms1, err1 := mb.Build()

	// use our own key -> empty tranlation
	mb.WithInvocation(InvocationKey)
	ms2, err2 := mb.Build()

	// set value
	en.Set(InvocationKey, []string{"new invocation"})
	ms3, err3 := mb.Build()

	assert.NoError(t, err1)
	assert.Equal(t, en.Get(l10n.KeySkillInvocation), ms1["en-US"].Model.Language.Invocation)
	assert.NoError(t, err2)
	assert.Empty(t, ms2["en-US"].Model.Language.Invocation)
	assert.NoError(t, err3)
	assert.Equal(t, "new invocation", ms3["en-US"].Model.Language.Invocation)
}

// modelBuilder with delegation is covered.
func TestModelBuilder_WithDelegationStrategy(t *testing.T) {
	assert.NotNil(t, registry)
	mb := skill.NewModelBuilder().
		WithLocaleRegistry(registry)

	mb.WithDelegationStrategy(skill.DelegationSkillResponse)
	ms2, err2 := mb.Build()

	assert.NoError(t, err2)
	assert.Equal(t, skill.DelegationSkillResponse, ms2["en-US"].Model.Dialog.Delegation)

	mb.WithDelegationStrategy("foo")
	_, err1 := mb.Build()
	assert.Error(t, err1)

}

// modelBuilder with locale is covered.
func TestModelBuilder_WithLocale(t *testing.T) {
	r := l10n.NewRegistry()
	mb := skill.NewModelBuilder().
		WithLocaleRegistry(r)

	// add two locales
	mb.WithLocale("en-US", "invoke")
	ms1, err1 := mb.Build()
	mb.WithLocale("fr-FR", "lance")
	ms2, err2 := mb.Build()

	// cannot add locale twice
	mb.WithLocale("en-US", "invoke2")
	_, err3 := mb.Build()

	assert.NoError(t, err1)
	assert.Len(t, ms1, 1)
	assert.NoError(t, err2)
	assert.Len(t, ms2, 2)

	assert.Error(t, err3)
}

// modelBuilder with intent is covered.
func TestModelBuilder_WithIntent(t *testing.T) {
	const IntentSamplesKey string = "WithIntent_Samples"
	assert.NotNil(t, registry)
	en, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	mb := skill.NewModelBuilder().
		WithLocaleRegistry(registry)

	// add intent
	samples := en.GetAll(IntentSamplesKey)
	mb.WithIntent("WithIntent")
	i1 := mb.Intent("WithIntent")
	m1, err1 := mb.BuildLocale("en-US")

	// add second intent
	mb.WithIntent("SlotIntent")
	//mb.Intent("SlotIntent")
	m2, err2 := mb.BuildLocale("en-US")

	// overwrite first intent
	mb.WithIntent("WithIntent")
	i3 := mb.Intent("WithIntent").
		WithLocaleSamples("en-US", []string{"new", "samples"})
	m3, err3 := mb.BuildLocale("en-US")

	assert.NoError(t, err1)
	assert.Equal(t, samples, m1.Model.Language.Intents[0].Samples)
	assert.Len(t, m1.Model.Language.Intents, 1)

	assert.NoError(t, err2)
	assert.Len(t, m2.Model.Language.Intents, 2)

	assert.NoError(t, err3)
	assert.Equal(t, i1, i3)
	assert.Len(t, m3.Model.Language.Intents, 2)
}

// modelBuilder with type is covered.
func TestModelBuilder_WithType(t *testing.T) {
	const TypeValuesKey string = "MyType_Values"
	assert.NotNil(t, registry)
	en, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	mb := skill.NewModelBuilder().
		WithLocaleRegistry(registry)

	// add type
	mb.WithType("MyType")
	val := en.Get(TypeValuesKey)
	mt1 := mb.Type("MyType")
	m1, err1 := mb.BuildLocale("en-US")

	// add second type
	mb.WithType("bar")
	mb.Type("bar")
	m2, err2 := mb.BuildLocale("en-US")

	// overwrite type with values
	mb.WithType("MyType")
	mt3 := mb.Type("MyType").
		WithLocaleValues("en-US", []string{"foo"})
	m3, err3 := mb.BuildLocale("en-US")

	assert.NoError(t, err1)
	assert.Equal(t, val, m1.Model.Language.Types[0].Values[0].Name.Value)
	assert.Len(t, m1.Model.Language.Types, 1)

	assert.NoError(t, err2)
	assert.Len(t, m2.Model.Language.Types, 2)

	assert.NoError(t, err3)
	assert.Equal(t, mt1, mt3)
	assert.Len(t, m3.Model.Language.Types, 2)
	assert.Equal(t, []string{"foo"}, en.GetAll(TypeValuesKey))
}

// modelBuilder with slot prompts are covered.
func TestModelBuilder_WithSlotPrompt(t *testing.T) {
	// no matching intent
	mb := skill.NewModelBuilder()
	mb.WithConfirmationSlotPrompt("Intent", "Slot")
	mb.WithElicitationSlotPrompt("Intent", "Slot")
	_, err1 := mb.Build()
	_, err2 := mb.BuildLocale("en-US")

	// with matching intent
	mb = skill.NewModelBuilder()
	mb.WithIntent("Intent").Intent("Intent").
		WithSlot("Slot", "Foo")
	mb.WithConfirmationSlotPrompt("Intent", "Slot")
	mb.WithElicitationSlotPrompt("Intent", "Slot")
	ms1, err3 := mb.Build()

	// error locale not registered
	ms2, err4 := mb.BuildLocale("en-US")

	assert.Error(t, err1)
	assert.Error(t, err2)
	assert.Equal(t, err1, err2)

	assert.NoError(t, err3)
	assert.Empty(t, ms1)

	assert.Error(t, err4)
	assert.Empty(t, ms2)
}

// modelBuilder errors if no locale is covered.
func TestModelBuilder_ErrorsIfNoLocale(t *testing.T) {
	mb := skill.NewModelBuilder()

	ms1, err1 := mb.Build()
	m2, err2 := mb.BuildLocale("en-US")

	assert.NoError(t, err1)
	assert.Empty(t, ms1)

	assert.Error(t, err2)
	assert.Equal(t, &skill.Model{}, m2)
}

// modelIntentBuilder with samples is covered.
func TestModelIntentBuilder_WithSamples(t *testing.T) {
	// setup
	const IntentSamplesKey string = "MyIntentSamples"
	assert.NotNil(t, registry)
	en, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	mib := skill.NewModelIntentBuilder("MyIntent").
		WithLocaleRegistry(registry)

	// uses default key
	i1, err1 := mib.BuildLanguageIntent("en-US")
	samples := en.GetAll("MyIntent_Samples")

	// change key name -> no translations
	mib.WithSamples(IntentSamplesKey)
	i2, err2 := mib.BuildLanguageIntent("en-US")

	// set translations for new key
	en.Set(IntentSamplesKey, []string{"new", "samples"})
	i3, err3 := mib.BuildLanguageIntent("en-US")

	assert.NoError(t, err1)
	assert.Equal(t, samples, i1.Samples)

	assert.NoError(t, err2)
	assert.Empty(t, i2.Samples)

	assert.NoError(t, err3)
	assert.Equal(t, en.GetAll(IntentSamplesKey), i3.Samples)
}

// modelIntentBuilder with locale samples is covered.
func TestModelIntentBuilder_WithLocaleSamples(t *testing.T) {
	assert.NotNil(t, registry)
	en, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	mib := skill.NewModelIntentBuilder("MyIntent").
		WithLocaleRegistry(registry)

	// uses default key
	i1, err1 := mib.BuildLanguageIntent("en-US")
	samples := en.GetAll("MyIntent_Samples")

	// overwrites translations of default key
	mib.WithLocaleSamples("en-US", []string{"samples"})
	i2, err2 := mib.BuildLanguageIntent("en-US")

	assert.NoError(t, err1)
	assert.Equal(t, samples, i1.Samples)
	assert.NoError(t, err2)
	assert.Equal(t, []string{"samples"}, i2.Samples)
}

// modelIntentBuilder with slot is covered.
func TestModelIntentBuilder_WithSlot(t *testing.T) {
	assert.NotNil(t, registry)
	mib := skill.NewModelIntentBuilder("Intent").
		WithLocaleRegistry(registry)

	// same name -> overwrites
	mib.WithSlot("Foo", "FooType")
	sl1 := mib.Slot("Foo")
	mib.WithSlot("Foo", "BarType")
	sl2 := mib.Slot("Foo")
	i1, err1 := mib.BuildLanguageIntent("en-US")

	// different name -> adds
	mib.WithSlot("Bar", "BarType")
	i2, err2 := mib.BuildLanguageIntent("en-US")

	assert.NoError(t, err1)
	assert.NotEqual(t, sl1, sl2)
	assert.Len(t, i1.Slots, 1)
	assert.Equal(t, "BarType", i1.Slots[0].Type)

	assert.NoError(t, err2)
	assert.Len(t, i2.Slots, 2)
}

// modelIntentBuilder errors if no locale is covered.
func TestModelIntentBuilder_ErrorsIfNoLocale(t *testing.T) {
	mib := skill.NewModelIntentBuilder("Intent").
		WithSlot("Foo", "FooType")

	_, err1 := mib.BuildLanguageIntent("en-US")
	// fails when looping over slots
	_, err2 := mib.BuildDialogIntent("en-US")
	mib.WithLocaleSamples("en-US", []string{"samples"})
	_, err3 := mib.BuildLanguageIntent("en-US")

	assert.Error(t, err1)
	assert.Error(t, err2)
	assert.Error(t, err3)
}

// modelSlotBuilder with samples is covered.
func TestModelSlotBuilder_WithSamples(t *testing.T) {
	// setup
	const SlotSamplesKey string = "MySlotSamples"
	assert.NotNil(t, registry)
	en, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	msb := skill.NewModelSlotBuilder("SlotIntent", "SlotName", "MyType").
		WithLocaleRegistry(registry)

	// uses default key
	is1, err1 := msb.BuildIntentSlot("en-US")
	samples := en.GetAll("SlotIntent_SlotName_Samples")

	// change key name -> no translations
	msb.WithSamples(SlotSamplesKey)
	is2, err2 := msb.BuildIntentSlot("en-US")

	// set translations for new key
	en.Set(SlotSamplesKey, []string{"sample1", "sample2"})
	is3, err3 := msb.BuildIntentSlot("en-US")

	assert.NoError(t, err1)
	assert.Equal(t, samples, is1.Samples)

	assert.NoError(t, err2)
	assert.Empty(t, is2.Samples)

	assert.NoError(t, err3)
	assert.Equal(t, en.GetAll(SlotSamplesKey), is3.Samples)
}

// modelSlotBuilder with locale samples is covered.
func TestModelSlotBuilder_WithLocaleSamples(t *testing.T) {
	assert.NotNil(t, registry)
	en, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	msb := skill.NewModelSlotBuilder("SlotIntent", "SlotName", "SlotType").
		WithLocaleRegistry(registry)

	// uses default key
	is1, err1 := msb.BuildIntentSlot("en-US")
	samples := en.GetAll("SlotIntent_SlotName_Samples")

	// overwrites translations of default key
	msb.WithLocaleSamples("en-US", []string{"bar"})
	is2, err2 := msb.BuildIntentSlot("en-US")

	assert.NoError(t, err1)
	assert.Equal(t, samples, is1.Samples)
	assert.NoError(t, err2)
	assert.Equal(t, en.GetAll("SlotIntent_SlotName_Samples"), is2.Samples)
}

// modelSlotBuilder with elicitation/confirmation prompt is covered.
func TestModelSlotBuilder_WithPrompt(t *testing.T) {
	// setup
	en := l10n.NewLocale("en-US")
	registry := l10n.NewRegistry()
	err := registry.Register(en)
	assert.NoError(t, err)
	msb := skill.NewModelSlotBuilder("MyIntent", "MySlot", "SlotType").
		WithLocaleRegistry(registry).
		WithElicitationPrompt("elicitation_id").
		WithConfirmationPrompt("confirmation_id").
		WithElicitation(false).
		WithConfirmation(true)

	ds, err := msb.BuildDialogSlot("en-US")

	assert.NoError(t, err)
	assert.Equal(t, "elicitation_id", ds.Prompts.Elicitation)
	assert.False(t, ds.Elicitation)
	assert.Equal(t, "confirmation_id", ds.Prompts.Confirmation)
	assert.True(t, ds.Confirmation)
}

// modelSlotBuilder with intent confirmation prompt is covered.
func TestModelSlotBuilder_WithIntentConfirmationPrompt(t *testing.T) {
	msb := skill.NewModelSlotBuilder("MyIntent", "MySlot", "SlotType")

	msb2 := msb.WithIntentConfirmationPrompt("foo")

	assert.Equal(t, msb, msb2)
}

// TODO: test modelValidationRulesBuilder

// modelSlotBuilder errors if no locale is covered.
func TestModelSlotBuilder_ErrorsIfNoLocale(t *testing.T) {
	msb := skill.NewModelSlotBuilder("MyIntent", "MySlot", "SlotType")

	_, err1 := msb.BuildIntentSlot("en-US")
	_, err2 := msb.BuildDialogSlot("en-US")
	msb.WithLocaleSamples("en-US", []string{})
	_, err3 := msb.BuildIntentSlot("en-US")

	assert.Error(t, err1)
	assert.Error(t, err2)
	assert.Equal(t, err1, err2)
	assert.Error(t, err3)
}

// modelTypeBuilder with values is covered.
func TestModelTypeBuilder_WithValues(t *testing.T) {
	assert.NotNil(t, registry)
	mtb := skill.NewModelTypeBuilder("FooBar").
		WithLocaleRegistry(registry)

	// default lookup key does not exist
	tvs1, err := mtb.Build("en-US")
	// set lookup key
	mtb.WithValues("MyType_Values")
	tvs2, err2 := mtb.Build("en-US")

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.Empty(t, tvs1.Values)
	assert.NotEmpty(t, tvs2.Values)
}

// modelTypeBuilder with locale values is covered.
func TestModelTypeBuilder_WithLocaleValues(t *testing.T) {
	assert.NotNil(t, registry)
	en, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	mtb := skill.NewModelTypeBuilder("MyType").
		WithLocaleRegistry(registry)

	// uses default key
	tvs1, err := mtb.Build("en-US")
	vals := en.GetAll("MyType_Values")

	// overwrite value of default key
	mtb.WithLocaleValues("en-US", []string{"my", "values", "foo"})
	tvs2, err2 := mtb.Build("en-US")

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.Equal(t, len(vals), len(tvs1.Values))
	assert.Len(t, tvs2.Values, 3)
}

// modelTypeBuilder errors if no locale is covered.
func TestModelTypeBuilder_ErrorsIfNoLocale(t *testing.T) {
	mtb := skill.NewModelTypeBuilder("Type").
		WithLocaleRegistry(l10n.NewRegistry())

	_, err := mtb.Build("en-US")
	mtb.WithLocaleValues("en-US", []string{"foo"})
	_, err2 := mtb.Build("en-US")

	assert.Error(t, err)
	assert.Error(t, err2)
}

// ModelPromptBuilder new elicitation prompt is covered.
func TestNewElicitationPromptBuilder(t *testing.T) {
	assert.NotNil(t, registry)
	en, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	mpb := skill.NewElicitationPromptBuilder("SlotIntent", "SlotName").
		WithLocaleRegistry(registry).
		WithVariation("PlainText")

	pvs, err := mpb.Variation("PlainText").BuildLocale("en-US")

	assert.NoError(t, err)
	assert.NotEmpty(t, pvs)
	assert.Equal(t, en.Get("SlotIntent_SlotName_Elicit_Text"), pvs[0].Value)
}

// ModelPromptBuilder new confirmation prompt is covered.
func TestNewConfirmationPromptBuilder(t *testing.T) {
	assert.NotNil(t, registry)
	en, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	mpb := skill.NewConfirmationPromptBuilder("SlotIntent", "SlotName").
		WithLocaleRegistry(registry).
		WithVariation("SSML")

	pvs, err := mpb.Variation("SSML").BuildLocale("en-US")

	assert.NoError(t, err)
	assert.NotEmpty(t, pvs)
	assert.Equal(t, en.Get("SlotIntent_SlotName_Confirm_SSML"), pvs[0].Value)
}

// ModelPromptBuilder with variation is covered.
func TestModelPromptBuilder_WithVariation(t *testing.T) {
	assert.NotNil(t, registry)
	mpb := skill.NewElicitationPromptBuilder("SlotIntent", "SlotName").
		WithLocaleRegistry(registry)

	// add "PlainText" with 2 texts
	mpb.WithVariation("PlainText")
	pvbText := mpb.Variation("PlainText")
	pvs, err := mpb.BuildLocale("en-US")

	// add "SSML" with 1 text
	mpb.WithVariation("SSML")
	pvbSSML := mpb.Variation("SSML")
	pvs2, err2 := mpb.BuildLocale("en-US")

	// overwrite "PlainText" - no change
	mpb.WithVariation("PlainText")
	pvbText2 := mpb.Variation("PlainText")
	pvs3, err3 := mpb.BuildLocale("en-US")

	assert.NoError(t, err)
	assert.NotNil(t, pvbText)
	assert.Len(t, pvs.Variations, 2)

	assert.NoError(t, err2)
	assert.NotNil(t, pvbSSML)
	assert.Len(t, pvs2.Variations, 3)

	assert.NoError(t, err3)
	assert.Equal(t, pvbText, pvbText2)
	assert.Len(t, pvs3.Variations, 3)
}

// ModelPromptBuilder errors if no variations is covered.
func TestModelPromptBuilder_ErrorsIfNoVariations(t *testing.T) {
	assert.NotNil(t, registry)
	mpb := skill.NewConfirmationPromptBuilder("MyIntent", "MySlot").
		WithLocaleRegistry(registry)

	_, err := mpb.BuildLocale("en-US")

	assert.Error(t, err)
}

// ModelPromptBuilder errors if no locale is covered.
func TestModelPromptBuilder_ErrorsIfNoLocale(t *testing.T) {
	mpb := skill.NewElicitationPromptBuilder("MyIntent", "MySlot").
		WithVariation("PlainText")

	_, err := mpb.BuildLocale("de-DE")

	assert.Error(t, err)
}

// promptVariationsBuilder defining our own lookup key is covered.
func TestPromptVariationsBuilder_WithTypeValue(t *testing.T) {
	assert.NotNil(t, registry)
	pvb := skill.NewPromptVariations("SlotIntent", "SlotName", "Elicit", "SSML").
		WithLocaleRegistry(registry)
	l, err := registry.Resolve("en-US")
	assert.NoError(t, err)
	l.Set("MySSMLLookupKey", []string{"foo", "bar"})

	pvb.WithTypeValue("SSML", "MySSMLLookupKey")
	pvs, err := pvb.BuildLocale("en-US")

	assert.NoError(t, err)
	assert.Len(t, pvs, 2)
	assert.Equal(t, "bar", pvs[1].Value)
}

// promptVariationsBuilder set variations translations directly is covered.
func TestPromptVariationsBuilder_WithLocaleTypeValue(t *testing.T) {
	r := l10n.NewRegistry()
	err := r.Register(l10n.NewLocale("en-US"))
	assert.NoError(t, err)
	pvb := skill.NewPromptVariations("MyIntent", "MySLot", "Elicit", "PlainText").
		WithLocaleRegistry(r).
		WithLocaleTypeValue("en-US", "PlainText", []string{"bar", "foo"})

	pvs, err := pvb.BuildLocale("en-US")

	assert.NoError(t, err)
	assert.Equal(t, "bar", pvs[0].Value)
}

// promptVariationsBuilder with variations is covered.
func TestPromptVariationsBuilder_WithVariation(t *testing.T) {
	assert.NotNil(t, registry)
	pvb := skill.NewPromptVariations("SlotIntent", "SlotName", "Elicit", "SSML").
		WithLocaleRegistry(registry)

	_, err := pvb.BuildLocale("en-US")
	pvb.WithVariation("PlainText")
	pvs, err2 := pvb.BuildLocale("en-US")
	// overwrites original "SSML" variations -> no change
	pvb.WithVariation("SSML")

	assert.NoError(t, err)
	assert.NoError(t, err2)
	// 3 elements as Text has two translations
	assert.Len(t, pvs, 3)
}

// promptVariationsBuilder errors if no locale is registered is covered.
func TestPromptVariationsBuilder_ErrorsIfNoLocale(t *testing.T) {
	pvb := skill.NewPromptVariations("MyIntent", "MySlot", "Elicit", "PlainText")

	_, err := pvb.BuildLocale("en-US")
	pvb.WithLocaleTypeValue("en-US", "PlainText", []string{"foo", "bar"})
	_, err2 := pvb.BuildLocale("en-US")

	assert.Error(t, err)
	assert.Error(t, err2)
}

// promtVariationsBuilder errors if translations are missing is covered.
func TestPromptVariationsBuilder_ErrorsIfNoTranslations(t *testing.T) {
	pvb := skill.NewPromptVariations("MyIntent", "NoSlot", "Elicit", "PlainText").
		WithLocaleRegistry(registry)

	_, err := pvb.BuildLocale("en-US")

	assert.Error(t, err)
}
