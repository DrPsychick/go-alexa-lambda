// Package skill contains the builders for the skill manifest and interaction models.
package skill

import (
	"fmt"

	"github.com/drpsychick/go-alexa-lambda/l10n"
)

// modelBuilder builds an alexa.Model instance for a locale.
type modelBuilder struct {
	registry   l10n.LocaleRegistry
	invocation string
	delegation string
	intents    map[string]*modelIntentBuilder
	types      map[string]*modelTypeBuilder
	prompts    map[string]*ModelPromptBuilder
	error      error
}

// NewModelBuilder returns an initialized modelBuilder.
func NewModelBuilder() *modelBuilder { //nolint:revive
	return &modelBuilder{
		registry:   l10n.NewRegistry(),
		invocation: l10n.KeySkillInvocation,
		delegation: DelegationAlways,
		intents:    map[string]*modelIntentBuilder{},
		types:      map[string]*modelTypeBuilder{},
		prompts:    map[string]*ModelPromptBuilder{},
	}
}

// WithLocaleRegistry passes a locale registry.
func (m *modelBuilder) WithLocaleRegistry(r l10n.LocaleRegistry) *modelBuilder {
	m.registry = r
	return m
}

// WithInvocation sets the lookup key for the invocation.
func (m *modelBuilder) WithInvocation(invocation string) *modelBuilder {
	m.invocation = invocation
	return m
}

// WithDelegationStrategy sets the model delegation strategy.
func (m *modelBuilder) WithDelegationStrategy(strategy string) *modelBuilder {
	if strategy != DelegationAlways && strategy != DelegationSkillResponse {
		m.error = fmt.Errorf("unsupported 'delegation': %s", strategy)
		return m
	}
	m.delegation = strategy
	return m
}

// WithLocale creates and sets a new locale.
func (m *modelBuilder) WithLocale(locale, invocation string) *modelBuilder {
	loc := l10n.NewLocale(locale)
	if err := m.registry.Register(loc); err != nil {
		m.error = err
		return m
	}
	loc.Set(m.invocation, []string{invocation})
	return m
}

// WithIntent creates and sets a new named intent.
func (m *modelBuilder) WithIntent(name string) *modelBuilder {
	i := NewModelIntentBuilder(name).
		WithLocaleRegistry(m.registry)
	m.intents[name] = i
	return m
}

// WithType creates and sets a new named type.
func (m *modelBuilder) WithType(name string) *modelBuilder {
	t := NewModelTypeBuilder(name).
		WithLocaleRegistry(m.registry)
	m.types[name] = t
	return m
}

// WithElicitationSlotPrompt creates and sets an elicitation prompt for the intent-slot.
func (m *modelBuilder) WithElicitationSlotPrompt(intent, slot string) *modelBuilder {
	// intent and slot must exist!
	var sl *modelSlotBuilder
	for _, i := range m.intents {
		for _, s := range i.slots {
			if i.name == intent && s.name == slot {
				sl = s
				break
			}
		}
	}
	if sl == nil {
		m.error = fmt.Errorf("no matching intent slot: %s-%s", intent, slot)
		return m
	}

	p := NewElicitationPromptBuilder(intent, slot).
		WithLocaleRegistry(m.registry)
	m.prompts[p.id] = p

	// link slot to prompt
	sl.WithElicitationPrompt(p.id)
	return m
}

// WithConfirmationSlotPrompt creates and sets a confirmation prompt for the intent-slot.
func (m *modelBuilder) WithConfirmationSlotPrompt(intent, slot string) *modelBuilder {
	// intent and slot must exist!
	var sl *modelSlotBuilder
	for _, i := range m.intents {
		for _, s := range i.slots {
			if i.name == intent && s.name == slot {
				sl = s
				break
			}
		}
	}
	if sl == nil {
		m.error = fmt.Errorf("no matching intent slot: %s-%s", intent, slot)
		return nil
	}

	p := NewConfirmationPromptBuilder(intent, slot).
		WithLocaleRegistry(m.registry)
	m.prompts[p.id] = p

	// link slot to prompt
	sl.WithConfirmationPrompt(p.id)
	return m
}

// WithValidationSlotPrompt creates and sets a validation prompt for a slot dialog.
func (m *modelBuilder) WithValidationSlotPrompt(slot, t string, valuesKey ...string) *modelBuilder {
	// slot must exist!
	var sl *modelSlotBuilder
	for _, i := range m.intents {
		for _, s := range i.slots {
			if s.name == slot {
				sl = s
				break
			}
		}
	}
	if sl == nil {
		m.error = fmt.Errorf("no matching intent slot: %s", slot)
		return nil
	}

	p := NewValidationPromptBuilder(slot, t).
		WithLocaleRegistry(m.registry)
	m.prompts[p.id] = p

	// link slot to prompt
	sl.WithValidationRule(t, p.id, valuesKey...)
	return m
}

// Intent returns the named intent.
func (m *modelBuilder) Intent(name string) *modelIntentBuilder {
	return m.intents[name]
}

// Type returns the named type.
func (m *modelBuilder) Type(name string) *modelTypeBuilder {
	return m.types[name]
}

// ElicitationPrompt returns the elicitation prompt for the intent-slot.
func (m *modelBuilder) ElicitationPrompt(intent, slot string) *ModelPromptBuilder {
	pb := NewElicitationPromptBuilder(intent, slot)
	return m.prompts[pb.id]
}

// ConfirmationPrompt returns the confirmation prompt for the intent-slot.
func (m *modelBuilder) ConfirmationPrompt(intent, slot string) *ModelPromptBuilder {
	pb := NewConfirmationPromptBuilder(intent, slot)
	return m.prompts[pb.id]
}

// ValidationPrompt returns the validation prompt for the intent-slot
// TODO: a slot can have multiple validation prompts! ID is not unique!
func (m *modelBuilder) ValidationPrompt(intent, slot string) *ModelPromptBuilder {
	pb := NewValidationPromptBuilder(intent, slot)
	return m.prompts[pb.id]
}

// Build generates a Model for each locale.
func (m *modelBuilder) Build() (map[string]*Model, error) {
	if m.error != nil {
		return nil, m.error
	}
	ams := make(map[string]*Model)

	// build model for each locale registered
	for _, l := range m.registry.GetLocales() {
		m, err := m.BuildLocale(l.GetName())
		if err != nil {
			return nil, err
		}
		ams[l.GetName()] = m
	}
	return ams, nil
}

// BuildLocale generates a Model for the locale.
func (m *modelBuilder) BuildLocale(locale string) (*Model, error) {
	if m.error != nil {
		return nil, m.error
	}
	loc, err := m.registry.Resolve(locale)
	if err != nil {
		return &Model{}, err
	}
	// create basic model
	am := &Model{
		Model: InteractionModel{
			Language: LanguageModel{
				Invocation: loc.Get(m.invocation),
			},
		},
	}

	mts := []ModelType{}
	for _, t := range m.types {
		mt, err := t.Build(locale)
		if err != nil {
			return &Model{}, err
		}
		mts = append(mts, mt)
	}
	am.Model.Language.Types = mts

	// add prompts - only if we have intents with slots
	// TODO: "Add...Prompt" should not fail, it should fail during build()!
	am.Model.Prompts = []ModelPrompt{}
	for _, p := range m.prompts {
		mp, err := p.BuildLocale(locale)
		if err != nil {
			return &Model{}, err
		}
		am.Model.Prompts = append(am.Model.Prompts, mp)
	}

	// add intents
	// TODO: ensure that slot types are defined, if not: fail
	am.Model.Dialog = &Dialog{}
	if m.delegation != "" {
		am.Model.Dialog.Delegation = m.delegation
	}
	for _, i := range m.intents {
		li, err := i.BuildLanguageIntent(locale)
		if err != nil {
			return &Model{}, err
		}
		am.Model.Language.Intents = append(am.Model.Language.Intents, li)

		// only needed for intents with slots
		if len(i.slots) > 0 {
			di, err := i.BuildDialogIntent(locale)
			if err != nil {
				return &Model{}, err
			}
			am.Model.Dialog.Intents = append(am.Model.Dialog.Intents, di)
		}
	}
	return am, nil
}

type modelIntentBuilder struct {
	registry     l10n.LocaleRegistry
	name         string
	samplesName  string
	delegation   string
	confirmation bool
	slots        map[string]*modelSlotBuilder
	error        error
}

// NewModelIntentBuilder returns an initialized modelIntentBuilder.
func NewModelIntentBuilder(name string) *modelIntentBuilder { //nolint:revive
	return &modelIntentBuilder{
		registry:     l10n.NewRegistry(),
		name:         name,
		samplesName:  name + l10n.KeyPostfixSamples,
		delegation:   DelegationAlways, // lets alexa or lambda handle the dialog for intent slots
		confirmation: false,            // should alexa ask to confirm the intent?
		slots:        map[string]*modelSlotBuilder{},
	}
}

// WithLocaleRegistry passes a locale registry.
func (i *modelIntentBuilder) WithLocaleRegistry(registry l10n.LocaleRegistry) *modelIntentBuilder {
	i.registry = registry
	return i
}

// WithSamples overwrites the locale lookup key.
func (i *modelIntentBuilder) WithSamples(samplesName string) *modelIntentBuilder {
	i.samplesName = samplesName
	return i
}

// WithLocaleSamples sets the lookup key translations for a specific locale.
func (i *modelIntentBuilder) WithLocaleSamples(locale string, samples []string) *modelIntentBuilder {
	loc, err := i.registry.Resolve(locale)
	if err != nil {
		i.error = err
		return i
	}
	loc.Set(i.samplesName, samples)
	return i
}

// WithSlot creates and sets a named slot for the intent.
func (i *modelIntentBuilder) WithSlot(name, typeName string) *modelIntentBuilder {
	sb := NewModelSlotBuilder(i.name, name, typeName).
		WithLocaleRegistry(i.registry)
	i.slots[name] = sb
	return i
}

// Slot returns a named slot of the intent.
func (i *modelIntentBuilder) Slot(name string) *modelSlotBuilder {
	return i.slots[name]
}

// WithDelegation sets the dialog delegation for the intent.
func (i *modelIntentBuilder) WithDelegation(d string) *modelIntentBuilder {
	if d != DelegationAlways && d != DelegationSkillResponse {
		i.error = fmt.Errorf("unsupported 'delegation': %s", d)
		return i
	}
	i.delegation = d
	return i
}

// WithConfirmation sets the dialog confirmation for the intent.
func (i *modelIntentBuilder) WithConfirmation(c bool) *modelIntentBuilder {
	i.confirmation = c
	return i
}

// BuildLanguageIntent generates a ModelIntent for the locale.
func (i *modelIntentBuilder) BuildLanguageIntent(locale string) (ModelIntent, error) {
	loc, err := i.registry.Resolve(locale)
	if err != nil {
		return ModelIntent{}, err
	}

	mi := ModelIntent{
		Name:    i.name,
		Samples: loc.GetAll(i.samplesName),
	}

	mss := []ModelSlot{}
	for _, s := range i.slots {
		is, err := s.BuildIntentSlot(locale)
		if err != nil {
			return ModelIntent{}, err
		}
		mss = append(mss, is)
	}
	mi.Slots = mss

	return mi, nil
}

// BuildDialogIntent generates a DialogIntent for the locale.
func (i *modelIntentBuilder) BuildDialogIntent(locale string) (DialogIntent, error) {
	di := DialogIntent{
		Name:         i.name,
		Delegation:   i.delegation,
		Confirmation: i.confirmation,
	}
	dis := []DialogIntentSlot{}
	for _, s := range i.slots {
		ds, err := s.BuildDialogSlot(locale)
		if err != nil {
			return DialogIntent{}, err
		}
		dis = append(dis, ds)
	}
	di.Slots = dis
	return di, nil
}

type modelSlotBuilder struct {
	registry           l10n.LocaleRegistry
	intent             string
	name               string
	typeName           string
	samplesName        string
	withConfirmation   bool
	withElicitation    bool
	elicitationPrompt  string
	confirmationPrompt string
	validationRules    *modelValidationRulesBuilder
}

// NewModelSlotBuilder returns an initialized modelSlotBuilder.
func NewModelSlotBuilder(intent, name, typeName string) *modelSlotBuilder { //nolint:revive
	return &modelSlotBuilder{
		registry:    l10n.NewRegistry(),
		intent:      intent,
		name:        name,
		typeName:    typeName,
		samplesName: intent + "_" + name + l10n.KeyPostfixSamples,
	}
}

// WithLocaleRegistry passes a locale registry.
func (s *modelSlotBuilder) WithLocaleRegistry(registry l10n.LocaleRegistry) *modelSlotBuilder {
	s.registry = registry
	return s
}

// WithSamples set the lookup key for the slot samples.
func (s *modelSlotBuilder) WithSamples(samplesName string) *modelSlotBuilder {
	s.samplesName = samplesName
	return s
}

// WithLocaleSamples sets the translated slot samples for the locale.
func (s *modelSlotBuilder) WithLocaleSamples(locale string, samples []string) *modelSlotBuilder {
	loc, err := s.registry.Resolve(locale)
	if err != nil {
		return s
	}
	loc.Set(s.samplesName, samples)
	return s
}

// WithConfirmation sets confirmationRequired for the slot.
func (s *modelSlotBuilder) WithConfirmation(c bool) *modelSlotBuilder {
	s.withConfirmation = c
	return s
}

// WithConfirmationPrompt requires confirmation and links to the prompt id.
func (s *modelSlotBuilder) WithConfirmationPrompt(id string) *modelSlotBuilder {
	s.withConfirmation = true
	s.confirmationPrompt = id
	return s
}

// WithElicitation sets elicitationRequired for the slot.
func (s *modelSlotBuilder) WithElicitation(e bool) *modelSlotBuilder {
	s.withElicitation = e
	return s
}

// WithElicitationPrompt requires elicitation and links to the prompt id.
func (s *modelSlotBuilder) WithElicitationPrompt(id string) *modelSlotBuilder {
	s.withElicitation = true
	s.elicitationPrompt = id
	return s
}

// WithIntentConfirmationPrompt does nothing.
func (s *modelSlotBuilder) WithIntentConfirmationPrompt(prompt string) *modelSlotBuilder {
	// TODO: WithIntentConfirmationPrompt
	// https://developer.amazon.com/docs/custom-skills/
	// -> define-the-dialog-to-collect-and-confirm-required-information.html#intent-confirmation
	// https://developer.amazon.com/en-US/docs/alexa/custom-skills/dialog-interface-reference.html#confirmintent
	return s
}

// WithValidationRule adds a validation rule to the slot.
func (s *modelSlotBuilder) WithValidationRule(t, prompt string, valuesKey ...string) *modelSlotBuilder {
	if nil == s.validationRules {
		s.validationRules = NewModelValidationRulesBuilder().
			WithLocaleRegistry(s.registry)
	}
	s.validationRules.WithRule(t, prompt, valuesKey...)
	return s
}

// BuildIntentSlot generates a ModelSlot for the locale.
func (s *modelSlotBuilder) BuildIntentSlot(locale string) (ModelSlot, error) {
	l, err := s.registry.Resolve(locale)
	if err != nil {
		return ModelSlot{}, err
	}
	ms := ModelSlot{
		Name: s.name,
		Type: s.typeName,
	}
	ms.Samples = l.GetAll(s.samplesName)
	return ms, nil
}

// BuildDialogSlot generates a DialogIntentSlot for the locale.
func (s *modelSlotBuilder) BuildDialogSlot(locale string) (DialogIntentSlot, error) {
	if _, err := s.registry.Resolve(locale); err != nil {
		return DialogIntentSlot{}, err
	}
	ds := DialogIntentSlot{
		Name:         s.name,
		Type:         s.typeName,
		Confirmation: s.withConfirmation,
		Elicitation:  s.withElicitation,
	}
	if s.confirmationPrompt != "" {
		ds.Prompts.Confirmation = s.confirmationPrompt
	}
	if s.elicitationPrompt != "" {
		ds.Prompts.Elicitation = s.elicitationPrompt
	}
	if s.validationRules != nil {
		vs, err := s.validationRules.BuildRules(locale)
		if err != nil {
			return ds, err
		}
		ds.Validations = vs
	}
	return ds, nil
}

// modelTypeBuilder.
type modelTypeBuilder struct {
	registry   l10n.LocaleRegistry
	name       string
	valuesName string
}

// NewModelTypeBuilder returns an initialized modelTypeBuilder.
func NewModelTypeBuilder(name string) *modelTypeBuilder { //nolint:revive
	return &modelTypeBuilder{
		registry:   l10n.NewRegistry(),
		name:       name,
		valuesName: name + l10n.KeyPostfixValues,
	}
}

// WithLocaleRegistry passes a locale registry.
func (t *modelTypeBuilder) WithLocaleRegistry(registry l10n.LocaleRegistry) *modelTypeBuilder {
	t.registry = registry
	return t
}

// WithValues sets the lookup key for the type values.
func (t *modelTypeBuilder) WithValues(valuesName string) *modelTypeBuilder {
	t.valuesName = valuesName
	return t
}

// WithLocaleValues sets the translated values for the type.
func (t *modelTypeBuilder) WithLocaleValues(locale string, values []string) *modelTypeBuilder {
	loc, err := t.registry.Resolve(locale)
	if err != nil {
		return t
	}
	loc.Set(t.valuesName, values)
	return t
}

// Build generates a ModelType.
func (t *modelTypeBuilder) Build(locale string) (ModelType, error) {
	loc, err := t.registry.Resolve(locale)
	if err != nil {
		return ModelType{}, err
	}
	tvs := []TypeValue{}
	for _, v := range loc.GetAll(t.valuesName) {
		tvs = append(tvs, TypeValue{Name: NameValue{Value: v}})
	}
	return ModelType{Name: t.name, Values: tvs}, nil
}

type modelValidationRulesBuilder struct {
	registry l10n.LocaleRegistry
	rules    []modelValidationRule
}

type modelValidationRule struct {
	validationType string
	prompt         string
	valuesKey      string
}

// NewModelValidationRulesBuilder returns an initialized modelValidationBuilder.
func NewModelValidationRulesBuilder() *modelValidationRulesBuilder { //nolint:revive
	return &modelValidationRulesBuilder{
		registry: l10n.NewRegistry(),
		rules:    []modelValidationRule{},
	}
}

// WithLocaleRegistry passes a locale registry.
func (v *modelValidationRulesBuilder) WithLocaleRegistry(registry l10n.LocaleRegistry) *modelValidationRulesBuilder {
	v.registry = registry
	return v
}

// WithRule adds a validation rule.
func (v *modelValidationRulesBuilder) WithRule(t, p string, valuesKey ...string) *modelValidationRulesBuilder {
	vr := modelValidationRule{
		validationType: t,
		prompt:         p,
	}
	if len(valuesKey) > 0 {
		vr.valuesKey = valuesKey[0]
	}
	v.rules = append(v.rules, vr)
	return v
}

func (v *modelValidationRulesBuilder) BuildRules(locale string) ([]SlotValidation, error) {
	sv := []SlotValidation{}
	loc, err := v.registry.Resolve(locale)
	if err != nil {
		return sv, err
	}

	// create and append SlotValidations
	for _, r := range v.rules {
		val := SlotValidation{
			Type:   r.validationType,
			Prompt: r.prompt,
		}

		if values := loc.GetAll(r.valuesKey); r.valuesKey != "" && len(values) > 0 {
			val.Values = values
		}

		// TODO: implement value:
		// https://developer.amazon.com/docs/smapi/interaction-model-schema.html#dialog_slot_validations
		// types isInSet/isNotInSet/... require values
		if (val.Type == ValidationTypeInSet || val.Type == ValidationTypeNotInSet) &&
			(val.Values == nil || len(val.Values) == 0) {
			return sv, fmt.Errorf("validation type requires values (%s: %s)", locale, val.Prompt)
		}

		sv = append(sv, val)
	}
	return sv, nil
}

// ModelPromptBuilder will build a ModelPrompt object.
type ModelPromptBuilder struct {
	registry   l10n.LocaleRegistry
	intent     string
	slot       string
	promptType string
	id         string
	variations map[string]*promptVariationsBuilder
}

// NewElicitationPromptBuilder returns an initialized ModelPromptBuilder for Elicitation.
func NewElicitationPromptBuilder(intent, slot string) *ModelPromptBuilder {
	return &ModelPromptBuilder{
		registry:   l10n.NewRegistry(),
		intent:     intent,
		slot:       slot,
		promptType: "Elicit",
		id:         fmt.Sprintf("Elicit.Intent-%s.IntentSlot-%s", intent, slot),
		variations: map[string]*promptVariationsBuilder{},
	}
}

// NewConfirmationPromptBuilder returns an initialized ModelPromptBuilder for Confirmation.
func NewConfirmationPromptBuilder(intent, slot string) *ModelPromptBuilder {
	return &ModelPromptBuilder{
		registry:   l10n.NewRegistry(),
		intent:     intent,
		slot:       slot,
		promptType: "Confirm",
		id:         fmt.Sprintf("Confirm.Intent-%s.IntentSlot-%s", intent, slot),
		variations: map[string]*promptVariationsBuilder{},
	}
}

// NewValidationPromptBuilder returns an initialized ModelPromptBuilder for Validation.
func NewValidationPromptBuilder(slot, t string) *ModelPromptBuilder {
	return &ModelPromptBuilder{
		registry:   l10n.NewRegistry(),
		slot:       slot,
		promptType: "Validate",
		id:         fmt.Sprintf("Validate.Slot-%s.Type-%s", slot, t),
		variations: map[string]*promptVariationsBuilder{},
	}
}

// WithLocaleRegistry passes a locale registry.
func (p *ModelPromptBuilder) WithLocaleRegistry(registry l10n.LocaleRegistry) *ModelPromptBuilder {
	p.registry = registry
	return p
}

// WithVariation creates and sets variations for the varType.
func (p *ModelPromptBuilder) WithVariation(varType string) *ModelPromptBuilder {
	v := NewPromptVariations(p.intent, p.slot, p.promptType, varType).
		WithLocaleRegistry(p.registry)
	p.variations[varType] = v
	return p
}

// Variation returns the variations for the varType.
func (p *ModelPromptBuilder) Variation(varType string) *promptVariationsBuilder { //nolint:revive
	return p.variations[varType]
}

// BuildLocale generates a ModelPrompt for the locale.
func (p *ModelPromptBuilder) BuildLocale(locale string) (ModelPrompt, error) {
	if len(p.variations) == 0 {
		return ModelPrompt{}, fmt.Errorf(
			"prompt '%s' requires variations (%s)",
			p.id, locale)
	}
	mp := ModelPrompt{
		ID:         p.id,
		Variations: []PromptVariation{},
	}
	for _, v := range p.variations {
		pv, err := v.BuildLocale(locale)
		if err != nil {
			return ModelPrompt{}, err
		}
		mp.Variations = append(mp.Variations, pv...)
	}
	return mp, nil
}

type promptVariationsBuilder struct {
	registry   l10n.LocaleRegistry
	intent     string
	slot       string
	promptType string
	vars       map[string]string
	error      error
}

// NewPromptVariations returns an initialized builder with lookup key "$intent_$slot_$promptType_(Text|SSML)".
func NewPromptVariations(intent, slot, promptType, varType string) *promptVariationsBuilder { //nolint:revive
	t := l10n.KeyPostfixSSML
	if varType == "PlainText" {
		t = l10n.KeyPostfixText
	}
	return &promptVariationsBuilder{
		registry:   l10n.NewRegistry(),
		intent:     intent,
		slot:       slot,
		promptType: promptType,
		// TODO: l10n key structure should depend on prompt type (without intent for validation prompts)
		vars: map[string]string{varType: fmt.Sprintf("%s_%s_%s%s", intent, slot, promptType, t)},
	}
}

// WithLocaleRegistry passes a locale registry.
func (v *promptVariationsBuilder) WithLocaleRegistry(registry l10n.LocaleRegistry) *promptVariationsBuilder {
	v.registry = registry
	return v
}

// WithVariation sets the lookup key for the varType.
func (v *promptVariationsBuilder) WithVariation(varType string) *promptVariationsBuilder {
	t := l10n.KeyPostfixSSML
	if varType == "PlainText" {
		t = l10n.KeyPostfixText
	}
	v.vars[varType] = fmt.Sprintf("%s_%s_%s%s", v.intent, v.slot, v.promptType, t)
	return v
}

// WithTypeValue sets valueName as the lookup key for the varType.
func (v *promptVariationsBuilder) WithTypeValue(varType, valueName string) *promptVariationsBuilder {
	v.vars[varType] = valueName
	return v
}

// WithLocaleTypeValue sets the values for the type of the locale.
func (v *promptVariationsBuilder) WithLocaleTypeValue(locale, varType string, values []string) *promptVariationsBuilder { //nolint:lll
	loc, err := v.registry.Resolve(locale)
	if err != nil {
		v.error = err
		return v
	}
	loc.Set(v.vars[varType], values)
	return v
}

// BuildLocale generates a PromptVariation for the locale.
func (v *promptVariationsBuilder) BuildLocale(locale string) ([]PromptVariation, error) {
	var vs []PromptVariation
	if v.error != nil {
		return vs, v.error
	}
	loc, err := v.registry.Resolve(locale)
	if err != nil {
		return vs, err
	}
	// only useful with content, can never happen as you must use NewPromptVariationsBuilder.
	if len(v.vars) == 0 {
		return []PromptVariation{}, fmt.Errorf(
			"prompt requires variations (%s: %s_%s_%s)",
			locale, v.intent, v.slot, v.promptType)
	}
	// loop over variation types
	for t, n := range v.vars {
		for _, val := range loc.GetAll(n) {
			vs = append(vs, PromptVariation{
				Type:  t,
				Value: val,
			})
		}
	}
	if len(vs) == 0 {
		return []PromptVariation{}, fmt.Errorf(
			"prompt requires variations with values (%s: %s_%s_%s)",
			locale, v.intent, v.slot, v.promptType)
	}
	return vs, nil
}
