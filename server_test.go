package alexa

import (
	"bytes"
	ctx "context"
	log "github.com/hamba/logger/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"
)

func TestServer(t *testing.T) {
	s := Server{
		Handler: HandlerFunc(
			func(b *ResponseBuilder, r *RequestEnvelope) {},
		),
	}
	context := ctx.Background()

	res, err := s.Invoke(context, []byte("{}"))
	resp := &ResponseEnvelope{}
	err2 := jsoniter.Unmarshal(res, resp)

	assert.NoError(t, err)
	assert.NoError(t, err2)
	assert.NotEmpty(t, resp)
}

func TestMuxServeHTTP(t *testing.T) {
	mux := NewServerMux(log.New(nil, log.ConsoleFormat(), log.Info))
	rw := httptest.NewRecorder()
	b := ioutil.NopCloser(bytes.NewReader([]byte(`{}`)))
	r := &http.Request{Method: http.MethodGet, Body: b}

	mux.ServeHTTP(rw, r)

	res, _ := ioutil.ReadAll(rw.Result().Body)
	resp := &ResponseEnvelope{}
	err := jsoniter.Unmarshal(res, resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
	assert.Contains(t, string(res), "error")

	rw = httptest.NewRecorder()
	b = ioutil.NopCloser(bytes.NewReader([]byte(`foo`)))
	r = &http.Request{Method: http.MethodGet, Body: b}

	mux.ServeHTTP(rw, r)

	res, _ = ioutil.ReadAll(rw.Result().Body)
	resp = &ResponseEnvelope{}
	err = jsoniter.Unmarshal(res, resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
	assert.Contains(t, string(res), "error")

	rw = httptest.NewRecorder()
	req := &RequestEnvelope{Request: &Request{Type: TypeIntentRequest, Intent: Intent{Name: HelpIntent}}}
	content, err := jsoniter.Marshal(req)
	assert.NoError(t, err)
	b = ioutil.NopCloser(bytes.NewReader(content))
	r = &http.Request{Method: http.MethodGet, Body: b}

	mux.ServeHTTP(rw, r)

	res, _ = ioutil.ReadAll(rw.Result().Body)
	resp = &ResponseEnvelope{}
	err = jsoniter.Unmarshal(res, resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
	assert.Contains(t, string(res), "error")
}

func TestHandler(t *testing.T) {
	mux := NewServerMux(log.New(nil, log.ConsoleFormat(), log.Info))
	h := HandlerFunc(func(b *ResponseBuilder, r *RequestEnvelope) { b.WithSimpleCard("title", "") })

	mux.HandleRequestType(TypeLaunchRequest, h)
	r := &RequestEnvelope{
		Request: &Request{
			Type: TypeLaunchRequest,
		},
	}
	_, err := mux.Handler(r)
	assert.NoError(t, err)

	mux.HandleIntent("Intent", h)
	r = &RequestEnvelope{
		Request: &Request{
			Type: TypeIntentRequest,
			Intent: Intent{
				Name: "Intent",
			},
		},
	}
	b := &ResponseBuilder{}

	_, err = mux.Handler(r)
	mux.Serve(b, r)
	out := b.Build()

	assert.NoError(t, err)
	assert.Equal(t, "title", out.Response.Card.Title)
}

func TestHandler_Errors(t *testing.T) {
	mux := NewServerMux(log.New(nil, log.ConsoleFormat(), log.Info))

	r := &RequestEnvelope{
		Request: &Request{
			Locale: "de-DE",
		},
	}
	_, err := mux.Handler(r)

	assert.Error(t, err)

	r.Request.Type = TypeIntentRequest
	r.Request.Intent = Intent{Name: "Test"}

	_, err = mux.Handler(r)

	assert.Error(t, err)
}

func TestHandleRequestType_IntentRequest(t *testing.T) {
	mux := NewServerMux(log.New(nil, log.ConsoleFormat(), log.Info))
	h := HandlerFunc(func(b *ResponseBuilder, r *RequestEnvelope) {})

	mux.HandleRequestTypeFunc(TypeIntentRequest, h)

	assert.Empty(t, mux.types)
}

func TestHandleIntentFunc(t *testing.T) {
	mux := NewServerMux(log.New(nil, log.ConsoleFormat(), log.Info))
	h := HandlerFunc(func(b *ResponseBuilder, r *RequestEnvelope) {})

	mux.HandleIntentFunc("Intent", h)

	funcName1 := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(mux.intents["Intent"]).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)
}

func TestServe(t *testing.T) {
	mux := NewServerMux(log.New(nil, log.ConsoleFormat(), log.Info))
	r := &RequestEnvelope{}
	b := &ResponseBuilder{}

	mux.Serve(b, r)

	assert.Equal(t, "Fatal error", b.card.Title)
}
