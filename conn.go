package quasizero

import (
	"bufio"
	"io"
	"net"

	pio "github.com/gogo/protobuf/io"
)

type protoConn struct {
	net.Conn
	r protoReader
	w protoWriter
}

func wrapConn(cn net.Conn) *protoConn {
	pc := new(protoConn)
	pc.Reset(cn)
	return pc
}

// Reset rewraps a native conn.
func (c *protoConn) Reset(cn net.Conn) {
	c.Conn = cn
	c.r.Reset(cn)
	c.w.Reset(cn)
}

// Close closes the conn.
func (c *protoConn) Close() (err error) {
	if e2 := c.r.Close(); e2 != nil {
		err = e2
	}
	if e2 := c.w.Close(); e2 != nil {
		err = e2
	}
	if e2 := c.Conn.Close(); e2 != nil {
		err = e2
	}
	return
}

type protoReader struct {
	buf *bufio.Reader
	pio.ReadCloser
}

// Buffered exposes number of bytes in the buffer.
func (pr *protoReader) Buffered() int { return pr.buf.Buffered() }

// Reset resets.
func (pr *protoReader) Reset(r io.Reader) {
	if pr.buf == nil {
		pr.buf = bufio.NewReader(r)
	} else {
		pr.buf.Reset(r)
	}
	pr.ReadCloser = pio.NewDelimitedReader(pr.buf, 1<<24)
}

type protoWriter struct {
	buf *bufio.Writer
	pio.WriteCloser
}

// Reset resets.
func (pw *protoWriter) Reset(w io.Writer) {
	if pw.buf == nil {
		pw.buf = bufio.NewWriter(w)
	} else {
		pw.buf.Reset(w)
	}
	pw.WriteCloser = pio.NewDelimitedWriter(pw.buf)
}

// Flush flushes the output buffer.
func (pw *protoWriter) Flush() error { return pw.buf.Flush() }
