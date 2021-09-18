package alexa

import (
	"github.com/drpsychick/alexa-go-lambda/ssml"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestWith_Functions(t *testing.T) {
	b := &ResponseBuilder{}

	b.WithSpeech("speech")
	b.WithSimpleCard("title", "text")
	b.WithShouldEndSession(true)
	b.WithReprompt("reprompt")

	res := b.Build()
	assert.Equal(t, "title", res.Response.Card.Title)
	assert.Equal(t, "text", res.Response.Card.Content)
	assert.Equal(t, "speech", res.Response.OutputSpeech.Text)
	assert.Equal(t, "reprompt", res.Response.Reprompt.OutputSpeech.Text)

	b.WithSpeech(ssml.Speak("speech"))
	b.WithReprompt(ssml.Speak("reprompt"))
	res = b.Build()
	assert.Equal(t, ssml.Speak("speech"), res.Response.OutputSpeech.SSML)
	assert.Equal(t, ssml.Speak("reprompt"), res.Response.Reprompt.OutputSpeech.SSML)

	b.WithStandardCard("title", "text", &Image{})
	res = b.Build()
	assert.Equal(t, "title", res.Response.Card.Title)
	assert.Equal(t, "text", res.Response.Card.Text)
	assert.Equal(t, &Image{}, res.Response.Card.Image)
	assert.Empty(t, res.Response.Card.Content)
}

func TestCanFulfillIntent(t *testing.T) {
	b := ResponseBuilder{}
	b.WithCanFulfillIntent(&CanFulfillIntent{
		CanFulfill: string(TypeLaunchRequest),
		Slots:      map[string]CanFulfillSlot{},
	})

	res := b.Build()
	assert.Equal(t, string(TypeLaunchRequest), res.Response.CanFulfillIntent.CanFulfill)
}

func TestAddDirective(t *testing.T) {
	b := &ResponseBuilder{}

	b.AddDirective(&Directive{
		Type:          DirectiveTypeDialogDelegate,
		SlotToElicit:  "",
		UpdatedIntent: nil,
		PlayBehavior:  "",
		AudioItem:     nil,
	})
	res := b.Build()

	assert.Equal(t, DirectiveTypeDialogDelegate, res.Response.Directives[0].Type)
}

func TestSessionAttributes(t *testing.T) {
	b := &ResponseBuilder{}

	b.WithSessionAttributes(map[string]interface{}{
		"foo": "Bar",
	})
	res := b.Build()
	assert.Equal(t, "Bar", res.SessionAttributes["foo"])
}

func TestResponseBuilder_With(t *testing.T) {
	type fields struct {
		speech           *OutputSpeech
		card             *Card
		reprompt         *OutputSpeech
		directives       []*Directive
		shouldEndSession bool
		sessionAttr      map[string]interface{}
		canFulfillIntent *CanFulfillIntent
	}
	type args struct {
		resp Response
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"WithEmptyResponse", fields{}, args{}},
		{"WithTextResponse", fields{}, args{Response{Title: "title", Text: "text"}}},
		{"WithSpeechResponse", fields{}, args{Response{Text: "text", Speech: "<speak>hello</speak>"}}},
		{"WithRepromptResponse", fields{}, args{Response{Text: "text", Speech: "<speak>hello</speak>", Reprompt: true}}},
		{"WithTextSpeechResponse", fields{}, args{Response{Text: "text", Speech: "hello", Reprompt: true}}},
		{"WithEndResponse", fields{}, args{Response{End: true}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ResponseBuilder{}
			b.With(tt.args.resp)

			if tt.args.resp.Title == "" {
				assert.Equal(t, "", b.card.Title)
			}
			if tt.args.resp.Title != "" {
				assert.Equal(t, tt.args.resp.Title, b.card.Title)
			}
			if tt.args.resp.Text != "" {
				assert.Equal(t, tt.args.resp.Text, b.card.Content)
			}
			if tt.args.resp.End {
				assert.True(t, b.shouldEndSession)
			}
			if tt.args.resp.Speech != "" {
				if strings.HasPrefix(tt.args.resp.Speech, `<speak>`) {
					if tt.args.resp.Reprompt {
						assert.Equal(t, tt.args.resp.Speech, b.reprompt.SSML)
					} else {
						assert.Equal(t, tt.args.resp.Speech, b.speech.SSML)
					}
				} else {
					if tt.args.resp.Reprompt {
						assert.Equal(t, tt.args.resp.Speech, b.reprompt.Text)
					} else {
						assert.Equal(t, tt.args.resp.Speech, b.speech.Text)
					}
				}
			}
		})
	}
}
