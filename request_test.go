package alexa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntent(t *testing.T) {
	r := &RequestEnvelope{Request: &Request{}}
	_, err := r.Intent()

	assert.Error(t, err)
	assert.Empty(t, r.IntentName())
	assert.False(t, r.IsIntentConfirmed())

	r.Request.Intent = Intent{Name: "Intent"}
	i, err := r.Intent()

	assert.NoError(t, err)
	assert.Equal(t, "Intent", i.Name)

	assert.False(t, r.IsIntentRequest())
	assert.False(t, r.IsIntentConfirmed())

	r.Request.Type = TypeIntentRequest
	assert.True(t, r.IsIntentRequest())

	r.Request.Intent.ConfirmationStatus = ConfirmationStatusConfirmed
	assert.True(t, r.IsIntentConfirmed())
}

func TestSlot(t *testing.T) {
	r := &RequestEnvelope{Request: &Request{}}
	_, err := r.Slot("Foo")

	assert.Error(t, err)
	assert.Empty(t, r.Slots())
	assert.Empty(t, r.SlotValue("foo"))

	r.Request.Intent = Intent{Name: "Intent"}
	assert.Empty(t, r.Slots())
	assert.Empty(t, r.SlotValue("bar"))

	r.Request.Intent.Slots = map[string]*Slot{
		"Foo": {Name: "Foo"},
	}
	assert.Len(t, r.Slots(), 1)
	assert.Equal(t, "Foo", r.Slots()["Foo"].Name)
	assert.Empty(t, r.SlotValue("Foo"))

	_, err = r.Slot("Bar")
	assert.Error(t, err)

	s, err := r.Slot("Foo")
	assert.NoError(t, err)
	assert.Equal(t, "", s.Value)

	r.Request.Intent.Slots["Foo"].Value = "Bar"
	assert.Equal(t, "Bar", r.SlotValue("Foo"))
}

func TestSlotResolutions(t *testing.T) {
	r := &RequestEnvelope{
		Request: &Request{
			Intent: Intent{
				Name: "Foo",
				Slots: map[string]*Slot{
					"Slot": &Slot{
						Value: "Value",
					},
				},
			},
		}}

	s, _ := r.Slot("Slot")
	_, err := s.SlotResolutionsPerAuthority()
	_, err2 := s.FirstAuthorityWithMatch()
	assert.Error(t, err)
	assert.Error(t, err2)

	r.Request.Intent.Slots["Slot"].Resolutions = &Resolutions{}
	s, _ = r.Slot("Slot")
	auths, err := s.SlotResolutionsPerAuthority()
	match, err2 := s.FirstAuthorityWithMatch()
	assert.NoError(t, err)
	assert.Error(t, err2)
	assert.Empty(t, auths)

	r.Request.Intent.Slots["Slot"].Resolutions.ResolutionsPerAuthority = []*PerAuthority{
		{
			Authority: "",
			Status: &ResolutionStatus{
				Code: ResolutionStatusNoMatch,
			},
			Values: []*AuthorityValue{
				{
					Value: &AuthorityValueValue{
						Name: "Foo",
						ID:   "123",
					},
				},
			},
		},
	}

	s, _ = r.Slot("Slot")
	auths, err = s.SlotResolutionsPerAuthority()
	match, err2 = s.FirstAuthorityWithMatch()
	assert.NoError(t, err)
	assert.Error(t, err2)
	assert.Empty(t, match)
	assert.Len(t, auths, 1)

	r.Request.Intent.Slots["Slot"].Resolutions.ResolutionsPerAuthority[0].Status.Code = ResolutionStatusMatch
	s, _ = r.Slot("Slot")
	auths, err = s.SlotResolutionsPerAuthority()
	match, err2 = s.FirstAuthorityWithMatch()
	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.Equal(t, auths[0], match)
	assert.Len(t, auths, 1)
}

func TestLocale(t *testing.T) {
	r := &RequestEnvelope{}

	assert.Empty(t, r.RequestLocale())

	r.Request = &Request{Locale: LocaleGerman}

	assert.Equal(t, string(LocaleGerman), r.RequestLocale())
}

func TestDialogState(t *testing.T) {
	r := &RequestEnvelope{}

	assert.Empty(t, r.RequestDialogState())

	r.Request = &Request{DialogState: DialogStateStarted}

	assert.Equal(t, DialogStateStarted, r.RequestDialogState())
}

func TestApplicationID(t *testing.T) {
	r := &RequestEnvelope{}

	_, err := r.ApplicationID()
	assert.Error(t, err)

	r.Context = &Context{
		System: &ContextSystem{
			Application: &ContextApplication{
				ApplicationID: "foo",
			},
		},
	}
	ID, err := r.ApplicationID()

	assert.Equal(t, "foo", ID)

	r.Session = &Session{}
	ID, err = r.ApplicationID()

	assert.Error(t, err)
	assert.Equal(t, "", ID)

	r.Session.Application = &ContextApplication{
		ApplicationID: "bar",
	}
	ID, err = r.ApplicationID()

	assert.Equal(t, "bar", ID)
}

func TestSessionID(t *testing.T) {
	r := &RequestEnvelope{}

	ID := r.SessionID()
	assert.Empty(t, ID)

	r.Session = &Session{}
	ID = r.SessionID()
	assert.Empty(t, ID)

	r.Session.SessionID = "foo"
	ID = r.SessionID()
	assert.Equal(t, "foo", ID)
}

func TestPersonAndUser(t *testing.T) {
	r := &RequestEnvelope{Context: &Context{}}

	_, err := r.ContextPerson()
	_, err2 := r.ContextUser()
	assert.Error(t, err)
	assert.Error(t, err2)

	r.Context.System = &ContextSystem{}
	_, err = r.ContextPerson()
	_, err2 = r.ContextUser()
	assert.Error(t, err)
	assert.Error(t, err2)

	s, _ := r.System()
	s.Person = &ContextSystemPerson{PersonID: "John"}
	s.User = &ContextUser{UserID: "dd9657c7-2246-4d09-b137-671d8de4b56f"}
	p, err := r.ContextPerson()
	u, err2 := r.ContextUser()
	assert.Equal(t, "John", p.PersonID)
	assert.Equal(t, "dd9657c7-2246-4d09-b137-671d8de4b56f", u.UserID)
}
