// Package skill serves generating the skill and model.
package skill

// https://developer.amazon.com/en-US/docs/alexa/smapi/interaction-model-schema.html#dialog_intents

// Model is the root of an interactionModel.
type Model struct {
	Model InteractionModel `json:"interactionModel"`
}

// InteractionModel defines the base model structure.
type InteractionModel struct {
	Language LanguageModel `json:"languageModel"`
	Dialog   *Dialog       `json:"dialog,omitempty"`
	Prompts  []ModelPrompt `json:"prompts,omitempty"`
}

// LanguageModel defines conversational primitives for the skill.
type LanguageModel struct {
	Invocation string        `json:"invocationName"`
	Intents    []ModelIntent `json:"intents"`
	Types      []ModelType   `json:"types,omitempty"`
}

// ModelIntent defines intents and their slots.
type ModelIntent struct {
	Name    string      `json:"name"`
	Samples []string    `json:"samples,omitempty"`
	Slots   []ModelSlot `json:"slots,omitempty"`
}

// ModelSlot defines slots within the intent.
type ModelSlot struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Samples []string `json:"samples,omitempty"`
}

// ModelType defines custom slot types.
type ModelType struct {
	Name   string      `json:"name"`
	Values []TypeValue `json:"values"`
}

// TypeValue defines a representative value for a custom slot type.
type TypeValue struct {
	ID   string    `json:"id,omitempty"`
	Name NameValue `json:"name"`
}

// NameValue defines a value of a custom slot type.
type NameValue struct {
	Value    string   `json:"value"`
	Synonyms []string `json:"synonyms,omitempty"`
}

// Dialog defines rules for conducting a multi-turn dialog with the user.
type Dialog struct {
	Delegation string         `json:"delegationStrategy"`
	Intents    []DialogIntent `json:"intents,omitempty"`
}

const (
	// DelegationSkillResponse delegates dialogs to lambda.
	DelegationSkillResponse string = "SKILL_RESPONSE"
	// DelegationAlways delegates dialogs to Alexa.
	DelegationAlways string = "ALWAYS"
)

// DialogIntent defines an intent that has dialog rules associated with it.
type DialogIntent struct {
	Name         string             `json:"name"`
	Confirmation bool               `json:"confirmationRequired"`
	Delegation   string             `json:"delegationStrategy,omitempty"`
	Prompts      IntentPrompt       `json:"prompts,omitempty"`
	Slots        []DialogIntentSlot `json:"slots,omitempty"`
}

// IntentPrompt defines a prompt for an intent.
type IntentPrompt struct {
	Confirmation string `json:"confirmation,omitempty"`
}

// DialogIntentSlot defines a slot in this intent that has dialog rules.
type DialogIntentSlot struct {
	Name         string           `json:"name"`
	Type         string           `json:"type"`
	Confirmation bool             `json:"confirmationRequired"`
	Elicitation  bool             `json:"elicitationRequired"`
	Prompts      SlotPrompts      `json:"prompts,omitempty"`
	Validations  []SlotValidation `json:"validations,omitempty"`
}

// SlotPrompts defines a collection of prompts for a slot.
type SlotPrompts struct {
	Elicitation  string `json:"elicitation,omitempty"`
	Confirmation string `json:"confirmation,omitempty"`
}

// see https://developer.amazon.com/docs/custom-skills/validate-slot-values.html#validation-rules
const (
	ValidationTypeHasMatch      string = "hasEntityResolutionMatch"
	ValidationTypeInSet         string = "isInSet"
	ValidationTypeNotInSet      string = "isNotInSet"
	ValidationTypeGreaterThan   string = "isGreaterThan"
	ValidationTypeGreaterEqal   string = "isGreaterThanOrEqualTo"
	ValidationTypeLessThan      string = "isLessThan"
	ValidationTypeLessEqual     string = "isLessThanOrEqualTo"
	ValidationTypeInDuration    string = "isInDuration"
	ValidationTypeNotInDuration string = "isNotInDuration"
)

// SlotValidation defines a validation rule for a prompt.
//
// see https://developer.amazon.com/docs/custom-skills/validate-slot-values.html#validation-rules
type SlotValidation struct {
	Type   string   `json:"type"`
	Prompt string   `json:"prompt"` //
	Values []string `json:"values,omitempty"`
}

// ModelPrompt defines cues to the user on behalf of the skill for eliciting data or providing feedback.
type ModelPrompt struct {
	ID         string            `json:"id"`
	Variations []PromptVariation `json:"variations"`
}

// PromptVariation defines a variation of the prompt.
type PromptVariation struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
