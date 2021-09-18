package ssml

import (
	"fmt"
	"github.com/drpsychick/alexa-go-lambda"
	"testing"
)

func TestBreak(t *testing.T) {
	type args struct {
		strength BreakStrength
		time     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Break", args{}, "<break/>"},
		{"BreakStrength", args{BreakStrengthWeak, ""}, `<break strength="weak"/>`},
		{"BreakTime", args{"", "1s"}, `<break time="1s"/>`},
		{"BreakText", args{BreakStrengthMedium, "5s"}, `<break strength="medium" time="5s"/>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Break(tt.args.strength, tt.args.time); got != tt.want {
				t.Errorf("Break() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestP(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Paragraph", args{"text"}, `<p>text</p>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := P(tt.args.text); got != tt.want {
				t.Errorf("P() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhoneme(t *testing.T) {
	type args struct {
		alphabet PhonemeAlphabet
		ph       string
		text     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Phoneme", args{PhonemeAlphabetXSampa, "n9f", "neuf"}, fmt.Sprintf(
			`<phoneme alphabet="%s" ph="%s">%s</phoneme>`, PhonemeAlphabetXSampa, "n9f", "neuf",
		)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Phoneme(tt.args.alphabet, tt.args.ph, tt.args.text); got != tt.want {
				t.Errorf("Phoneme() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProsody(t *testing.T) {
	type args struct {
		rate   ProsodyRate
		pitch  ProsodyPitch
		volume ProsodyVolume
		text   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"ProsodyNoArgs", args{}, `<prosody></prosody>`},
		{"ProsodyOnlyText", args{"", "", "", "text"}, `<prosody>text</prosody>`},
		{"ProsodyRate", args{ProsodyRateFast, "", "", "text"}, fmt.Sprintf(
			`<prosody rate="%s">%s</prosody>`, ProsodyRateFast, "text",
		)},
		{"ProsodyPitch", args{"", ProsodyPitchMedium, "", "text"}, fmt.Sprintf(
			`<prosody pitch="%s">%s</prosody>`, ProsodyPitchMedium, "text",
		)},
		{"ProsodyVolume", args{"", "", ProsodyVolumeLoud, "text"}, fmt.Sprintf(
			`<prosody volume="%s">%s</prosody>`, ProsodyVolumeLoud, "text",
		)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Prosody(tt.args.rate, tt.args.pitch, tt.args.volume, tt.args.text); got != tt.want {
				t.Errorf("Prosody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"SentenceNoArgs", args{}, `<s></s>`},
		{"Sentence", args{"foo bar"}, `<s>foo bar</s>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := S(tt.args.text); got != tt.want {
				t.Errorf("S() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSayAs(t *testing.T) {
	type args struct {
		interpretAs SayAsInterpretAs
		format      string
		text        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"SayAsNoArgs", args{}, `<say-as interpret-as=""></say-as>`},
		{"SayAsOrdinal", args{SayAsInterpretAsOrdinal, "", "123"},
			`<say-as interpret-as="ordinal">123</say-as>`,
		},
		{"SayAsDateFormat", args{SayAsInterpretAsDate, "dm", "5th September"},
			`<say-as interpret-as="date" format="dm">5th September</say-as>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SayAs(tt.args.interpretAs, tt.args.format, tt.args.text); got != tt.want {
				t.Errorf("SayAs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpeak(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"SpeakNoArgs", args{}, `<speak></speak>`},
		{"Speak", args{"foo"}, `<speak>foo</speak>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Speak(tt.args.text); got != tt.want {
				t.Errorf("Speak() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSub(t *testing.T) {
	type args struct {
		alias string
		text  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"SubNoArgs", args{}, `<sub alias=""></sub>`},
		{"Sub", args{"World of Warcraft", "WOW"}, `<sub alias="World of Warcraft">WOW</sub>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sub(tt.args.alias, tt.args.text); got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseAudio(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"AudioNoArgs", args{}, `<audio src=""/>`},
		{"Audio", args{"https://foo.bar"}, `<audio src="https://foo.bar"/>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UseAudio(tt.args.src); got != tt.want {
				t.Errorf("UseAudio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseDomain(t *testing.T) {
	type args struct {
		domain AmazonDomain
		text   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"DomainNoArgs", args{}, `<amazon:domain name=""></amazon:domain>`},
		{"DomainFun", args{AmazonDomainFun, "fun"},
			`<amazon:domain name="` + string(AmazonDomainFun) + `">fun</amazon:domain>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UseDomain(tt.args.domain, tt.args.text); got != tt.want {
				t.Errorf("UseDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseEffect(t *testing.T) {
	type args struct {
		effect AmazonEffect
		text   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"EffectNoArgs", args{}, `<amazon:effect name=""></amazon:effect>`},
		{"EffectWhisper", args{AmazonEffectWhispered, "don't tell anyone"},
			`<amazon:effect name="` + string(AmazonEffectWhispered) + `">don't tell anyone</amazon:effect>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UseEffect(tt.args.effect, tt.args.text); got != tt.want {
				t.Errorf("UseEffect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseEmotion(t *testing.T) {
	type args struct {
		name      AmazonEmotion
		intensity AmazonEmotionIntensity
		text      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"EmotionNoArgs", args{}, `<amazon:emotion name="" intensity=""></amazon:emotion>`},
		{"EmotionNoArgs", args{EmotionExcited, EmotionIntensityHigh, "Yes"},
			fmt.Sprintf(`<amazon:emotion name="%s" intensity="%s">%s</amazon:emotion>`,
				string(EmotionExcited), string(EmotionIntensityHigh), "Yes",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UseEmotion(tt.args.name, tt.args.intensity, tt.args.text); got != tt.want {
				t.Errorf("UseEmotion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseEmphasis(t *testing.T) {
	type args struct {
		level EmphasisLevel
		text  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"EmphasisNoArgs", args{}, `<emphasis></emphasis>`},
		{"EmphasisModerate", args{EmphasisLevelModerate, "hello"},
			`<emphasis level="` + string(EmphasisLevelModerate) + `">hello</emphasis>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UseEmphasis(tt.args.level, tt.args.text); got != tt.want {
				t.Errorf("UseEmphasis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseLang(t *testing.T) {
	type args struct {
		language string
		text     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"LangNoArgs", args{}, `<lang xml:lang=""></lang>`},
		{"Lang", args{"en-US", "Hello"}, `<lang xml:lang="en-US">Hello</lang>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UseLang(tt.args.language, tt.args.text); got != tt.want {
				t.Errorf("UseLang() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseVoice(t *testing.T) {
	type args struct {
		voice PollyVoice
		text  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"VoiceNoArgs", args{}, `<voice name=""></voice>`},
		{"Voice", args{AUVoiceNicole, "my name is nicole"},
			`<voice name="` + string(AUVoiceNicole) + `">my name is nicole</voice>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UseVoice(tt.args.voice, tt.args.text); got != tt.want {
				t.Errorf("UseVoice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseVoiceLang(t *testing.T) {
	type args struct {
		voice    PollyVoice
		language string
		text     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"VoiceLangNoArgs", args{}, `<voice name=""><lang xml:lang=""></lang></voice>`},
		{"VoiceLang", args{DEVoiceMarlene, string(alexa.LocaleGerman), "ich heisse Marlene"},
			fmt.Sprintf(`<voice name="%s"><lang xml:lang="%s">%s</lang></voice>`,
				string(DEVoiceMarlene), string(alexa.LocaleGerman), "ich heisse Marlene",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UseVoiceLang(tt.args.voice, tt.args.language, tt.args.text); got != tt.want {
				t.Errorf("UseVoiceLang() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestW(t *testing.T) {
	type args struct {
		role AmazonRole
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"WordNoArgs", args{}, `<w role=""></w>`},
		{"WordVB", args{AmazonRoleVB, "read"}, `<w role="` + string(AmazonRoleVB) + `">read</w>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := W(tt.args.role, tt.args.text); got != tt.want {
				t.Errorf("W() = %v, want %v", got, tt.want)
			}
		})
	}
}
