// Package quasizero implements a general purpose, ultra-low latency TCP server.
package quasizero

// Handler instances process commands.
type Handler interface {
	// ServeQZ serves a request.
	ServeQZ(*Request) (*Response, error)
}

// HandlerFunc is a Handler short-cut.
type HandlerFunc func(*Request) (*Response, error)

// ServeQZ implements the Handler interface.
func (f HandlerFunc) ServeQZ(req *Request) (*Response, error) { return f(req) }
