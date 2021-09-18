// Package ssml provides functions to simplify working with SSML speech.
// https://developer.amazon.com/en-US/docs/alexa/custom-skills/speech-synthesis-markup-language-ssml-reference.html#incompatible-tags
package ssml

import (
	"fmt"
	"strings"
)

// AmazonDomain is the domain of speech (news, music, ...).
// https://developer.amazon.com/en-US/docs/alexa/custom-skills/speech-synthesis-markup-language-ssml-reference.html#amazon-domain
type AmazonDomain string

const (
	AmazonDomainConversational AmazonDomain = "conversational" //nolint:revive
	AmazonDomainLong           AmazonDomain = "long-form"
	AmazonDomainMusic          AmazonDomain = "music"
	AmazonDomainNews           AmazonDomain = "news"
	AmazonDomainFun            AmazonDomain = "fun"
)

// UseDomain uses a specific domain of speech.
func UseDomain(domain AmazonDomain, text string) string {
	return `<amazon:domain name="` + string(domain) + `">` + text + `</amazon:domain>`
}

// AmazonEffect is a speech effect.
type AmazonEffect string

const (
	// AmazonEffectWhispered lowers the voice to be whispering.
	AmazonEffectWhispered AmazonEffect = "whispered"
)

// UseEffect wraps text in an effect.
func UseEffect(effect AmazonEffect, text string) string {
	return `<amazon:effect name="` + string(effect) + `">` + text + `</amazon:effect>`
}

// AmazonEmotion adds emotion to speech.
type AmazonEmotion string

const (
	// EmotionExcited is an exiting voice.
	EmotionExcited AmazonEmotion = "excited"
	// EmotionDisappointed is a disappointed voice.
	EmotionDisappointed AmazonEmotion = "disappointed"
)

// AmazonEmotionIntensity is the intensity of the emotion.
type AmazonEmotionIntensity string

const (
	EmotionIntensityLow    AmazonEmotionIntensity = "low" //nolint:revive
	EmotionIntensityMedium AmazonEmotionIntensity = "medium"
	EmotionIntensityHigh   AmazonEmotionIntensity = "high"
)

// UseEmotion wraps the text in an emotion tag.
func UseEmotion(name AmazonEmotion, intensity AmazonEmotionIntensity, text string) string {
	// <amazon:emotion name="excited" intensity="medium">
	return fmt.Sprintf(`<amazon:emotion name="%s" intensity="%s">%s</amazon:emotion>`,
		string(name), string(intensity), text,
	)
}

// UseAudio uses an URL for an MP3 file to play.
// https://developer.amazon.com/en-US/docs/alexa/custom-skills/speech-synthesis-markup-language-ssml-reference.html#audio
// <audio src="soundbank://soundlibrary/transportation/amzn_sfx_car_accelerate_01" />.
func UseAudio(src string) string {
	return `<audio src="` + src + `"/>`
}

// BreakStrength is one way to define the length of a break.
type BreakStrength string

const (
	BreakStrengthNone  BreakStrength = "none" //nolint:revive
	BreakStrengthXWeak BreakStrength = "x-weak"
	BreakStrengthWeak  BreakStrength = "weak"
	// BreakStrengthMedium is the default.
	BreakStrengthMedium  BreakStrength = "medium"
	BreakStrengthStrong  BreakStrength = "strong"
	BreakStrengthXStrong BreakStrength = "x-strong"
)

// Break adds a break to speech.
// only use `strength` or `time`, not both.
// time in `ms` or `s` - may not exceed 10s.
func Break(strength BreakStrength, time string) string {
	params := []string{""}
	if strength != "" {
		params = append(params, `strength="`+string(strength)+`"`)
	}
	if time != "" {
		params = append(params, `time="`+time+`"`)
	}
	return fmt.Sprintf(`<break%s/>`, strings.Join(params, " "))
}

// EmphasisLevel is the level of emphasis.
type EmphasisLevel string

const (
	// EmphasisLevelStrong is louder and slower.
	EmphasisLevelStrong EmphasisLevel = "strong"
	// EmphasisLevelModerate is the default emphasis if no level is given.
	EmphasisLevelModerate EmphasisLevel = "moderate"
	// EmphasisLevelReduced is softer and faster.
	EmphasisLevelReduced EmphasisLevel = "reduced"
)

// UseEmphasis adds emphasis to the give text.
func UseEmphasis(level EmphasisLevel, text string) string {
	params := []string{""}
	if level != "" {
		params = append(params, `level="`+string(level)+`"`)
	}
	return fmt.Sprintf(`<emphasis%s>%s</emphasis>`, strings.Join(params, " "), text)
}

// UseLang speaks given text in the specified language.
func UseLang(language, text string) string {
	return `<lang xml:lang="` + language + `">` + text + `</lang>`
}

// P wraps text in a paragraph.
func P(text string) string {
	return `<p>` + text + `</p>`
}

// PhonemeAlphabet is the alphabet to interpret the `ph` parameter with.
// https://developer.amazon.com/en-US/docs/alexa/custom-skills/speech-synthesis-markup-language-ssml-reference.html#phoneme
type PhonemeAlphabet string

const (
	// PhonemeAlphabetIPA is the International Phonetic Alphabet (IPA).
	PhonemeAlphabetIPA PhonemeAlphabet = "ipa"
	// PhonemeAlphabetXSampa the Extended Speech Assessment Methods Phonetic Alphabet (X-SAMPA).
	PhonemeAlphabetXSampa PhonemeAlphabet = "x-sampa"
)

// Phoneme pronounces the given text based on the provided alphabet and characters.
// <phoneme alphabet="ipa" ph="pɪˈkɑːn">pecan</phoneme>.
func Phoneme(alphabet PhonemeAlphabet, ph, text string) string {
	return `<phoneme alphabet="` + string(alphabet) + `" ph="` + ph + `">` + text + `</phoneme>`
}

// ProsodyRate defines the speed of the voice. Can be provided in %: 100% is normal speed.
type ProsodyRate string

const (
	ProsodyRateXSlow  ProsodyRate = "x-slow" //nolint:revive
	ProsodyRateSlow   ProsodyRate = "slow"
	ProsodyRateMedium ProsodyRate = "medium"
	ProsodyRateFast   ProsodyRate = "fast"
	ProsodyRateXFast  ProsodyRate = "x-fast"
)

// ProsodyPitch defines the pitch of the voice. Can be provided in %: positiv is higher, negative lower.
type ProsodyPitch string

const (
	ProsodyPitchXLow   ProsodyPitch = "x-low" //nolint:revive
	ProsodyPitchLow    ProsodyPitch = "low"
	ProsodyPitchMedium ProsodyPitch = "medium"
	ProsodyPitchHigh   ProsodyPitch = "high"
	ProsodyPitchXHigh  ProsodyPitch = "x-high"
)

// ProsodyVolume defines the volume of the voice. Can be provided in +/-dB, e.g. "+2dB".
type ProsodyVolume string

const (
	ProsodyVolumeSilent ProsodyVolume = "silent" //nolint:revive
	ProsodyVolumeXSoft  ProsodyVolume = "x-soft"
	ProsodyVolumeSoft   ProsodyVolume = "soft"
	ProsodyVolumeMedium ProsodyVolume = "medium"
	ProsodyVolumeLoud   ProsodyVolume = "loud"
	ProsodyVolumeXLoud  ProsodyVolume = "x-loud"
)

// Prosody modifies the volume, pitch, and rate of the tagged speech.
func Prosody(rate ProsodyRate, pitch ProsodyPitch, volume ProsodyVolume, text string) string {
	params := []string{""}
	if rate != "" {
		params = append(params, `rate="`+string(rate)+`"`)
	}
	if pitch != "" {
		params = append(params, `pitch="`+string(pitch)+`"`)
	}
	if volume != "" {
		params = append(params, `volume="`+string(volume)+`"`)
	}
	return fmt.Sprintf(`<prosody%s>%s</prosody>`, strings.Join(params, " "), text)
}

// S wraps text in a sentence.
func S(text string) string {
	return `<s>` + text + `</s>`
}

// SayAsInterpretAs is the type of pronunciation.
type SayAsInterpretAs string

const (
	// SayAsInterpretAsCharacters spells out each letter.
	SayAsInterpretAsCharacters SayAsInterpretAs = "characters"
	// SayAsInterpretAsCardinal interprets the value as a cardinal number.
	SayAsInterpretAsCardinal SayAsInterpretAs = "cardinal"
	// SayAsInterpretAsOrdinal interprets the value as an ordinal number.
	SayAsInterpretAsOrdinal SayAsInterpretAs = "ordinal"
	// SayAsInterpretAsDigits spells each digit separately.
	SayAsInterpretAsDigits SayAsInterpretAs = "digits"
	// SayAsInterpretAsFraction interprets the value as a fraction. This works for 3/20 as well as 1+1/2.
	SayAsInterpretAsFraction SayAsInterpretAs = "fraction"
	// SayAsInterpretAsUnit interprets a value as a measurement. Either a number or fraction including unit or just a unit.
	SayAsInterpretAsUnit SayAsInterpretAs = "unit"
	// SayAsInterpretAsDate interprets the value as a date. Specify the format with the format attribute.
	SayAsInterpretAsDate SayAsInterpretAs = "date"
	// SayAsInterpretAsTime interprets a value such as 1'21" as duration in minutes and seconds.
	SayAsInterpretAsTime SayAsInterpretAs = "time"
	// SayAsInterpretAsTelephone interprets a value as a 7-digit or 10-digit telephone number.
	// This can also handle extensions (for example, 2025551212x345).
	SayAsInterpretAsTelephone SayAsInterpretAs = "telephone"
	// SayAsInterpretAsAddress interprets a value as part of street address.
	SayAsInterpretAsAddress SayAsInterpretAs = "address"
	// SayAsInterpretAsInterjection interprets the value as an interjection.
	// Alexa speaks the text in a more expressive voice. For optimal results, only use the supported interjections
	// and surround each speechcon with a pause. For example: <say-as interpret-as="interjection">Wow.</say-as>.
	// Speechcons are supported for the languages listed below.
	// https://developer.amazon.com/en-US/docs/alexa/custom-skills/speech-synthesis-markup-language-ssml-reference.html#supported-speechcons
	SayAsInterpretAsInterjection SayAsInterpretAs = "interjection"
	// SayAsInterpretAsExpletive "Bleeps" out the content inside the tag.
	SayAsInterpretAsExpletive SayAsInterpretAs = "expletive"
)

// SayAs instructs the voice to "read" the text in a specific way.
// <say-as interpret-as="cardinal">12345</say-as>.
func SayAs(interpretAs SayAsInterpretAs, format, text string) string {
	if interpretAs == SayAsInterpretAsDate && format != "" {
		return `<say-as interpret-as="` + string(interpretAs) + `" format="` + format + `">` + text + `</say-as>`
	}
	return `<say-as interpret-as="` + string(interpretAs) + `">` + text + `</say-as>`
}

// Speak wraps text in <speak> tags.
func Speak(text string) string {
	return "<speak>" + text + "</speak>"
}

// Sub provides an alias for the voice (e.g. how to pronounce abbreviations) for the text.
// <sub alias="aluminum">Al</sub>
// <sub alias="if I remember correctly">IIRC</sub>.
func Sub(alias, text string) string {
	return `<sub alias="` + alias + `">` + text + `</sub>`
}

// PollyVoice defines the voice name for speech.
type PollyVoice string

// https://developer.amazon.com/en-US/docs/alexa/custom-skills/speech-synthesis-markup-language-ssml-reference.html#supported-voices
const (
	// US : Ivy, Joanna, Joey, Justin, Kendra, Kimberly, Matthew, Salli.
	USVoiceIvy      PollyVoice = "Ivy"
	USVoiceJoanna   PollyVoice = "Joanna"
	USVoiceJustin   PollyVoice = "Justin"
	USVoiceKendra   PollyVoice = "Kendra"
	USVoiceKimberly PollyVoice = "Kimberly"
	USVoiceMatthew  PollyVoice = "Matthew"
	USVoiceSalli    PollyVoice = "Salli"
	// AU : Nicole, Russell.
	AUVoiceNicole PollyVoice = "Nicole"
	AUVoiceRussel PollyVoice = "Russel"
	// GB : Amy, Brian, Emma.
	GBVoiceAmy   PollyVoice = "Amy"
	GBVoiceBrian PollyVoice = "Brian"
	GBVoiceEmma  PollyVoice = "Emma"
	// IN : Aditi, Raveena.
	INVoiceAditi   PollyVoice = "Aditi"
	INVoiceRaveena PollyVoice = "Raveena"
	// CA : Chantal.
	CAVoiceChantal PollyVoice = "Chantal"
	// FR : Celine, Lea, Mathieu.
	FRVoiceCeline  PollyVoice = "Celine"
	FRVoiceLea     PollyVoice = "Lea"
	FRVoiceMathieu PollyVoice = "Mathieu"
	// DE : Hans, Marlene, Vicki.
	DEVoiceHans    PollyVoice = "Hans"
	DEVoiceMarlene PollyVoice = "Marlene"
	DEVoiceVicki   PollyVoice = "Vicki"
	// HI : Aditi.
	HIVoiceAditi PollyVoice = "Aditi"
	// IT : Carla, Giorgio, Bianca.
	ITVoiceCarla   PollyVoice = "Carla"
	ITVoiceGiorgio PollyVoice = "Giorgio"
	ITVoiceBianca  PollyVoice = "Bianca"
	// JP : Mizuki, Takumi.
	JPVoiceMizuki PollyVoice = "Mitzuki"
	JPVoiceTakumi PollyVoice = "Takumi"
	// BR : Vitoria, Camila, Ricardo.
	BRVoiceVitoria PollyVoice = "Vitoria"
	BRVoiceCamila  PollyVoice = "Camila"
	BRVoiceRicardo PollyVoice = "Ricardo"
	// es-US : Penelope, Lupe, Miguel.
	EsUSVoicePenelope PollyVoice = "Penelope"
	EsUSVoiceLupe     PollyVoice = "Lupe"
	EsUSVoiceMiguel   PollyVoice = "Miguel"
	// ES : Conchita, Enrique, Lucia.
	ESVoiceConchita PollyVoice = "Conchita"
	ESVoiceEnrique  PollyVoice = "Enrique"
	ESVoiceLucia    PollyVoice = "Lucia"
	// MX : Mia.
	MXVoiceMia PollyVoice = "Mia"
)

// UseVoice wraps text in tags using a specific voice.
func UseVoice(voice PollyVoice, text string) string {
	return `<voice name="` + string(voice) + `">` + text + `</voice>`
}

// UseVoiceLang wraps text in tags using a specific voice and language.
func UseVoiceLang(voice PollyVoice, language, text string) string {
	return `<voice name="` + string(voice) + `"><lang xml:lang="` + language + `">` + text + `</lang></voice>`
}

// AmazonRole is a customized pronunciation of words.
type AmazonRole string

const (
	// AmazonRoleVB interprets the word as a verb (present simple).
	AmazonRoleVB AmazonRole = "amazon:VB"
	// AmazonRoleVBD interprets the word as a past participle.
	AmazonRoleVBD AmazonRole = "amazon:VBD"
	// AmazonRoleNN interprets the word as a noun.
	AmazonRoleNN AmazonRole = "amazon:NN"
	// AmazonRoleSense1 uses the non-default sense of the word.
	// For example, the noun "bass" is pronounced differently depending on meaning.
	// The "default" meaning is the lowest part of the musical range. The alternate sense
	// (which is still a noun) is a freshwater fish. Specifying
	// <speak><w role="amazon:SENSE_1">bass</w>"</speak> renders the non-default
	// pronunciation (freshwater fish).
	AmazonRoleSense1 AmazonRole = "amazon:SENSE-1"
)

// W wraps text for a customized pronunciation.
// <w role="amazon:VB">read</w>
// Similar to say-as, this tag customizes the pronunciation of words by specifying the word's part of speech.
func W(role AmazonRole, text string) string {
	return `<w role="` + string(role) + `">` + text + `</w>`
}
