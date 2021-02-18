package main

import "net/http"

// Middleware is a chain of http handlers
type Middleware []http.Handler

// MiddlewareResponseWriter marks if it has written to the http ResponseWriter
type MiddlewareResponseWriter struct {
	http.ResponseWriter
	written bool
}

// NewMiddlewareResponseWriter returns a new MiddlewareResponseWriter
func NewMiddlewareResponseWriter(w http.ResponseWriter) *MiddlewareResponseWriter {
	return &MiddlewareResponseWriter{
		ResponseWriter: w,
	}
}

func (w *MiddlewareResponseWriter) Write(bytes []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(bytes)
}

// WriteHeader writes the http header to the ResponseWriter
func (w *MiddlewareResponseWriter) WriteHeader(code int) {
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

// Add adds a handler to the middleware list
func (m *Middleware) Add(handler http.Handler) {
	*m = append(*m, handler)
}

// ServeHTTP processes all handlers in the chain and stops as soon as a
// handler has written to the ResponseWriter or returns a 404 if no handlers has
// written anything
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Process the middleware
	mw := NewMiddlewareResponseWriter(w)

	// Loop through all registered handlers
	for _, handler := range m {
		// call the handler with our MiddlewareResponseWriter
		handler.ServeHTTP(mw, r)

		// stop processing the chain when there has been an output
		if mw.written {
			return
		}
	}
	// If no handlers in the chain wrote a response, we return a 404
	http.NotFound(w, r)
}
