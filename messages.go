package quasizero

import (
	"fmt"
	"sync"
)

// SetMeta sets a key/value metadata pair.
func (m *Request) SetMeta(key, value string) {
	if m.Metadata == nil {
		m.Metadata = make(map[string]string, 1)
	}
	m.Metadata[key] = value
}

func (m *Request) reuse() {
	*m = Request{Payload: m.Payload[:0]}
}

// --------------------------------------------------------------------

var responsePool sync.Pool

func fetchResponse() *Response {
	if v := responsePool.Get(); v != nil {
		m := v.(*Response)
		m.reuse()
		return m
	}
	return new(Response)
}

// Set sets the payload efficiently.
func (m *Response) Set(data []byte) {
	m.Payload = append(m.Payload[:0], data...)
}

// SetString sets the payload efficiently.
func (m *Response) SetString(data string) {
	m.Payload = append(m.Payload[:0], data...)
}

// SetErrorf sets a formatted error message.
func (m *Response) SetErrorf(msg string, args ...interface{}) {
	m.ErrorMessage = fmt.Sprintf(msg, args...)
}

// SetError sets an error.
func (m *Response) SetError(err error) {
	if err != nil {
		m.ErrorMessage = err.Error()
	}
}

// SetMeta sets a key/value metadata pair.
func (m *Response) SetMeta(key, value string) {
	if m.Metadata == nil {
		m.Metadata = make(map[string]string, 1)
	}
	m.Metadata[key] = value
}

// Release releases the message and returns it to the memory pool.
// You must not use it after calling this function.
func (m *Response) Release() {
	responsePool.Put(m)
}

func (m *Response) reuse() {
	*m = Response{Payload: m.Payload[:0]}
}

// --------------------------------------------------------------------

// ResponseBatch is a slice of individual responses.
type ResponseBatch []*Response

// Release releases the response batch and returns it to the memory pool.
// You must not use the batch or any of the included responses after calling
// this function.
func (b ResponseBatch) Release() {
	for _, r := range b {
		r.Release()
	}
}
