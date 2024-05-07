package alexa

import (
	"context"
	"fmt"
	"io"
	l2 "log"
	"net/http"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/hamba/logger/v2"
	lctx "github.com/hamba/logger/v2/ctx"
	jsoniter "github.com/json-iterator/go"
)

// Handler represents an alexa request handler.
type Handler interface {
	Serve(builder *ResponseBuilder, req *RequestEnvelope)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// HandlerFunc is an adapter allowing a function to be used as a handler.
type HandlerFunc func(*ResponseBuilder, *RequestEnvelope)

// Serve serves the request.
func (fn HandlerFunc) Serve(b *ResponseBuilder, r *RequestEnvelope) {
	fn(b, r)
}

// ServeHTTP serves a HTTP request.
func (fn HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL != nil && (strings.HasSuffix(r.URL.Path, "/livez") || strings.HasSuffix(r.URL.Path, "/readyz")) {
		if _, err := rw.Write([]byte("ok")); err != nil {
			l2.Fatal("error: could not write response")
		}
		return
	}

	req, err := parseRequest(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(`{"error": "failed to parse request"}`))
		return
	}
	defer func() { _ = r.Body.Close() }()

	builder := &ResponseBuilder{}
	fn(builder, req)

	resp, err := jsoniter.Marshal(builder.Build())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		resp = []byte(`{"error": "failed to marshal response"}`)
	}
	_, _ = rw.Write(resp)
}

// A Server defines parameters for running an Alexa server.
type Server struct {
	Handler Handler
}

// Invoke calls the handler, and serializes the response.
func (s *Server) Invoke(_ context.Context, payload []byte) ([]byte, error) {
	req := &RequestEnvelope{}
	if err := jsoniter.Unmarshal(payload, req); err != nil {
		return nil, err
	}

	builder := &ResponseBuilder{}
	s.Handler.Serve(builder, req)

	// Idea: BuildJson -> then the `build()` can be private
	return jsoniter.Marshal(builder.Build())
}

// Serve serves the handler.
func (s *Server) Serve() error {
	if s.Handler == nil {
		s.Handler = DefaultServerMux
	}

	lambda.Start(s)
	return nil
}

// Serve serves the given handler.
func Serve(h Handler) error {
	srv := &Server{Handler: h}

	return srv.Serve()
}

// ServeMux is an Alexa request multiplexer.
type ServeMux struct {
	mu          sync.RWMutex
	logger      *log.Logger
	types       map[RequestType]Handler
	intents     map[string]Handler
	intentSlots map[string]string
}

// NewServerMux creates a new server mux.
func NewServerMux(log *log.Logger) *ServeMux {
	return &ServeMux{
		logger:      log,
		types:       map[RequestType]Handler{},
		intents:     map[string]Handler{},
		intentSlots: map[string]string{},
	}
}

// Logger returns the application logger.
func (m *ServeMux) Logger() *log.Logger {
	return m.logger
}

// Handler returns the matched handler for a request, or an error.
func (m *ServeMux) Handler(r *RequestEnvelope) (Handler, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if h, ok := m.types[r.RequestType()]; ok {
		return h, nil
	}

	if r.RequestType() != TypeIntentRequest {
		return nil, fmt.Errorf("server: unknown intent type %s", r.RequestType())
	}

	h, ok := m.intents[r.IntentName()]
	if !ok {
		return nil, fmt.Errorf("server: unknown intent %s", r.IntentName())
	}

	return h, nil
}

// HandleRequestType registers the handler for the given request type.
//
// Any attempt to handle the IntentRequest type will be ignored, use Intent instead.
func (m *ServeMux) HandleRequestType(requestType RequestType, handler Handler) {
	if requestType == TypeIntentRequest {
		return
	}
	if handler == nil {
		panic("alexa: nil handler")
	}

	m.mu.Lock()

	m.types[requestType] = handler

	m.mu.Unlock()
}

// HandleRequestTypeFunc registers the handler function for the given request type.
//
// Any attempt to handle the IntentRequest type will be ignored, use Intent instead.
func (m *ServeMux) HandleRequestTypeFunc(requestType RequestType, handler HandlerFunc) {
	m.HandleRequestType(requestType, handler)
}

// HandleIntent registers the handler for the given intent.
func (m *ServeMux) HandleIntent(intent string, handler Handler) {
	if handler == nil {
		panic("alexa: nil handler")
	}

	m.mu.Lock()

	m.intents[intent] = handler

	m.mu.Unlock()
}

// HandleIntentFunc registers the handler function for the given intent.
func (m *ServeMux) HandleIntentFunc(intent string, handler HandlerFunc) {
	m.HandleIntent(intent, handler)
}

// fallbackHandler returns a fatal error card.
func fallbackHandler(err error) HandlerFunc {
	return HandlerFunc(func(b *ResponseBuilder, _ *RequestEnvelope) {
		b.WithSimpleCard("Fatal error", "error: "+err.Error()).
			WithShouldEndSession(true)
	})
}

// Serve serves the matched handler.
func (m *ServeMux) Serve(b *ResponseBuilder, r *RequestEnvelope) {
	json, _ := jsoniter.Marshal(r)
	m.logger.Debug("request", lctx.Str("json", string(json)))
	h, err := m.Handler(r)
	if err != nil {
		h = fallbackHandler(err)
	}

	h.Serve(b, r)
	json, _ = jsoniter.Marshal(b.Build())
	m.logger.Debug("response", lctx.Str("json", string(json)))
}

// ServeHTTP dispatches the request to the handler whose
// alexa intent matches the request URL.
func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL != nil && (strings.HasSuffix(r.URL.Path, "/livez") || strings.HasSuffix(r.URL.Path, "/readyz")) {
		if _, err := w.Write([]byte("ok")); err != nil {
			m.logger.Debug("failed to write response")
		}
		return
	}

	var h Handler
	req, err := parseRequest(r.Body)
	if err != nil {
		h = fallbackHandler(err)
	} else {
		h, err = m.Handler(req)
		if err != nil {
			h = fallbackHandler(err)
		}
	}
	defer func() { _ = r.Body.Close() }()

	builder := &ResponseBuilder{}
	h.Serve(builder, req)

	resp, err := jsoniter.Marshal(builder.Build())
	if err != nil {
		m.logger.Error("failed to marshal response", lctx.Error("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		resp = []byte(`{"error": "failed to marshal response"}`)
	}
	if _, err := w.Write(resp); err != nil {
		m.logger.Debug("failed to write response")
	}
}

func parseRequest(b io.Reader) (*RequestEnvelope, error) {
	payload, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}

	req := &RequestEnvelope{}
	if err := jsoniter.Unmarshal(payload, req); err != nil {
		return nil, err
	}
	return req, nil
}

// DefaultServerMux is the default mux.
var DefaultServerMux = NewServerMux(log.New(nil, log.ConsoleFormat(), log.Info))

// HandleRequestType registers the handler for the given request type on the DefaultServeMux.
//
// Any attempt to handle the IntentRequest type will be ignored, use Intent instead.
func HandleRequestType(requestType RequestType, handler Handler) {
	DefaultServerMux.HandleRequestType(requestType, handler)
}

// HandleRequestTypeFunc registers the handler function for the given request type on the DefaultServeMux.
//
// Any attempt to handle the IntentRequest type will be ignored, use Intent instead.
func HandleRequestTypeFunc(requestType RequestType, handler HandlerFunc) {
	DefaultServerMux.HandleRequestTypeFunc(requestType, handler)
}

// HandleIntent registers the handler for the given intent on the DefaultServeMux.
func HandleIntent(intent string, handler Handler) {
	DefaultServerMux.HandleIntent(intent, handler)
}

// HandleIntentFunc registers the handler function for the given intent on the DefaultServeMux.
func HandleIntentFunc(intent string, handler HandlerFunc) {
	DefaultServerMux.HandleIntentFunc(intent, handler)
}
