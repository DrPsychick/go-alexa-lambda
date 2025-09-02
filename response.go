package alexa

import (
	"fmt"
	"strings"
)

// Stream represents a response directive audio item stream.
type Stream struct {
	Token                string `json:"token,omitempty"`
	URL                  string `json:"url,omitempty"`
	OffsetInMilliseconds int    `json:"offsetInMilliseconds,omitempty"`
}

// AudioItem represents a response directive audio item.
type AudioItem struct {
	Stream Stream `json:"stream,omitempty"`
}

// DirectiveType represents various Directive Types.
type DirectiveType string

// Directive types.
const (
	DirectiveTypeDialogDelegate      DirectiveType = "Dialog.Delegate"
	DirectiveTypeDialogElicitSlot    DirectiveType = "Dialog.ElicitSlot"
	DirectiveTypeDialogConfirmSlot   DirectiveType = "Dialog.ConfirmSlot"
	DirectiveTypeDialogConfirmIntent DirectiveType = "Dialog.ConfirmIntent"
)

// Directive represents a response directive.
type Directive struct {
	Type          DirectiveType `json:"type,omitempty"`
	SlotToElicit  string        `json:"slotToElicit,omitempty"`
	UpdatedIntent *Intent       `json:"updatedIntent,omitempty"`
	PlayBehavior  string        `json:"playBehavior,omitempty"`
	AudioItem     *AudioItem    `json:"audioItem,omitempty"`
}

// OutputSpeech represents a speech response.
type OutputSpeech struct {
	Type         string `json:"type"`
	Text         string `json:"text,omitempty"`
	SSML         string `json:"ssml,omitempty"`
	PlayBehavior string `json:"playBehavior,omitempty"`
}

// Card presents a card response.
type Card struct {
	Type    string `json:"type"`
	Title   string `json:"title,omitempty"`
	Text    string `json:"text,omitempty"`
	Content string `json:"content,omitempty"`
	Image   *Image `json:"image,omitempty"`
}

// Image represents a card image.
type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

// Reprompt represents a reprompt response.
type Reprompt struct {
	OutputSpeech *OutputSpeech `json:"outputSpeech,omitempty"`
}

// CanFulfillSlot represents a slots fulfillment.
type CanFulfillSlot struct {
	CanUnderstand string `json:"canUnderstand"`
	CanFulfill    string `json:"canFulfill"`
}

// CanFulfillIntent represents a response indicating if an intent
// can be fulfilled.
type CanFulfillIntent struct {
	CanFulfill string                    `json:"canFulfill"`
	Slots      map[string]CanFulfillSlot `json:"slots"`
}

// ResponseEnvelope represents the wrapper for a response.
type ResponseEnvelope struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          response               `json:"response"`
}

// response represents the response.
type response struct {
	OutputSpeech     *OutputSpeech     `json:"outputSpeech,omitempty"`
	Card             *Card             `json:"card,omitempty"`
	Reprompt         *Reprompt         `json:"reprompt,omitempty"`
	Directives       []*Directive      `json:"directives,omitempty"`
	ShouldEndSession bool              `json:"shouldEndSession"`
	CanFulfillIntent *CanFulfillIntent `json:"canFulfillIntent,omitempty"`
}

// ResponseBuilder builds a response.
type ResponseBuilder struct {
	speech           *OutputSpeech
	card             *Card
	reprompt         *OutputSpeech
	directives       []*Directive
	shouldEndSession bool
	sessionAttr      map[string]interface{}
	canFulfillIntent *CanFulfillIntent
}

// With applies an Response.
func (b *ResponseBuilder) With(resp Response) {
	if resp.Image != "" {
		b.WithStandardCard(resp.Title, resp.Text, &Image{
			SmallImageURL: fmt.Sprintf(resp.Image, "small"),
			LargeImageURL: fmt.Sprintf(resp.Image, "large"),
		})
	}

	b.WithSimpleCard(resp.Title, resp.Text)

	if resp.Speech != "" {
		if resp.Reprompt {
			b.WithReprompt(resp.Speech)
		} else {
			b.WithSpeech(resp.Speech)
		}
	}

	b.WithShouldEndSession(resp.End)
}

// WithSpeech sets the output speech on the response.
//
// If the text contains SSML speak tags, it will be set as SSML speech,
// otherwise it will be set as plain text speech.
func (b *ResponseBuilder) WithSpeech(text string) *ResponseBuilder {
	if strings.HasPrefix(text, "<speak>") && strings.HasSuffix(text, "</speak>") {
		b.speech = &OutputSpeech{
			Type: "SSML",
			SSML: text,
		}

		return b
	}

	b.speech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}

	return b
}

// WithReprompt sets the reprompt output speech on the response.
func (b *ResponseBuilder) WithReprompt(text string) *ResponseBuilder {
	if strings.HasPrefix(text, "<speak>") && strings.HasSuffix(text, "</speak>") {
		b.reprompt = &OutputSpeech{
			Type: "SSML",
			SSML: text,
		}

		return b
	}

	b.reprompt = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}

	return b
}

// WithSimpleCard sets a simple card on the response.
func (b *ResponseBuilder) WithSimpleCard(title, text string) *ResponseBuilder {
	b.card = &Card{
		Type:    "Simple",
		Title:   title,
		Content: text,
	}

	return b
}

// WithStandardCard sets a standard card on the response.
func (b *ResponseBuilder) WithStandardCard(title, text string, image *Image) *ResponseBuilder {
	b.card = &Card{
		Type:  "Standard",
		Title: title,
		Text:  text,
		Image: image,
	}

	return b
}

// WithShouldEndSession determines if the session should end after the current response.
func (b *ResponseBuilder) WithShouldEndSession(end bool) *ResponseBuilder {
	b.shouldEndSession = end
	return b
}

// WithSessionAttributes sets the session attributes on the response.
func (b *ResponseBuilder) WithSessionAttributes(attr map[string]interface{}) *ResponseBuilder {
	b.sessionAttr = attr

	return b
}

// WithCanFulfillIntent sets the can fulfill intent response on the response.
func (b *ResponseBuilder) WithCanFulfillIntent(response *CanFulfillIntent) *ResponseBuilder {
	b.canFulfillIntent = response

	return b
}

// AddDirective adds a directive tp the response.
func (b *ResponseBuilder) AddDirective(directive *Directive) *ResponseBuilder {
	b.directives = append(b.directives, directive)
	return b
}

// Build builds the response from the given information.
func (b *ResponseBuilder) Build() *ResponseEnvelope {
	// TODO: empty response with directive(s), like Dialog:Delegate
	r := &ResponseEnvelope{
		Version:           "1.0",
		SessionAttributes: b.sessionAttr,
		Response: response{
			OutputSpeech:     b.speech,
			Card:             b.card,
			Directives:       b.directives,
			ShouldEndSession: b.shouldEndSession,
		},
	}

	if b.reprompt != nil {
		r.Response.Reprompt = &Reprompt{
			OutputSpeech: b.reprompt,
		}
	}

	if b.canFulfillIntent != nil {
		r.Response.CanFulfillIntent = b.canFulfillIntent
	}

	return r
}
