package alexa

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/hamba/logger/v2"
	jsoniter "github.com/json-iterator/go"
)

// Handler represents an alexa request handler.
type Handler interface {
	Serve(*ResponseBuilder, *RequestEnvelope)
}

// HandlerFunc is an adapter allowing a function to be used as a handler.
type HandlerFunc func(*ResponseBuilder, *RequestEnvelope)

// Serve serves the request.
func (fn HandlerFunc) Serve(b *ResponseBuilder, r *RequestEnvelope) {
	fn(b, r)
}

// A Server defines parameters for running an Alexa server.
type Server struct {
	Handler Handler
}

// Invoke calls the handler, and serializes the response.
func (s *Server) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
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
	// TODO: decide if we want a DefaultServeMux
	if s.Handler == nil {
		return errors.New("alexa: cannot serve empty handler")
	}

	lambda.StartHandler(s)
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
	logger      log.Logger
	types       map[RequestType]Handler
	intents     map[string]Handler
	intentSlots map[string]string
}

// NewServerMux creates a new server mux.
func NewServerMux(log log.Logger) *ServeMux {
	return &ServeMux{
		logger:      log,
		types:       map[RequestType]Handler{},
		intents:     map[string]Handler{},
		intentSlots: map[string]string{},
	}
}

// Logger returns the application logger.
func (m *ServeMux) Logger() log.Logger {
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
	return HandlerFunc(func(b *ResponseBuilder, r *RequestEnvelope) {
		b.WithSimpleCard("Fatal error", "error: "+err.Error()).
			WithShouldEndSession(true)
	})
}

// Serve serves the matched handler.
func (m *ServeMux) Serve(b *ResponseBuilder, r *RequestEnvelope) {
	json, _ := jsoniter.Marshal(r)
	m.logger.Debug(string(json))
	h, err := m.Handler(r)
	if err != nil {
		h = fallbackHandler(err)
	}

	h.Serve(b, r)
	json, _ = jsoniter.Marshal(b.Build())
	m.logger.Debug(string(json))
}

// DefaultServerMux is the default mux.
var DefaultServerMux = NewServerMux(*log.New(nil, log.ConsoleFormat(), log.Info))

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
