# Purpose
Provide support in localizing an Alexa skill.
* clear and easy structure of translations (one file per locale encouraged)
* simple "key" lookup that allows placeholders (using `fmt.Sprintf`)
* register locales (translations), ~define fallback locales~
* separating logic from translations (logic/flow is in the code, e.g. which Intent uses which Slots)

What does `l10n` NOT provide or aim to support:
* it does not aim to be feature complete (yet)

# Usage
## Default keys
see [l10n.go](l10n.go)

Default "postfix" for keys:
`_Samples` for samples of an intent or slot.
`_Values` for a type

## Example:
see [skill_test.go](../skill/skill_test.go)

if you use multiple locales, it's easiest to define your own keys:
```go
// keys of the project
const(
	ByeBye       string = "byebye"
	StopTitle    string = "stop_title"
	Stop         string = "stop"
	GenericTitle string = "Alexa"

	// Intents
	DemoIntent                string = "DemoIntent"
	DemoIntentSamples         string = "DemoIntent_Samples"
	DemoIntentTitle           string = "DemoIntent_Title"
	DemoIntentText            string = "DemoIntent_Text"
	DemoIntentSSML            string = "DemoIntent_SSML"
)
```