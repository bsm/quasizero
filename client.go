package quasizero

import (
	"context"
	"net"

	"github.com/bsm/pool"
)

// Client holds a pool of connections to a quasizero server instance.
type Client struct {
	cns *pool.Pool
}

// NewClient connects a client.
func NewClient(ctx context.Context, addr string, opt *pool.Options) (*Client, error) {
	return NewClientDialer(ctx, new(net.Dialer), addr, opt)
}

// NewClientDialer connects a client through a custom dialer.
func NewClientDialer(ctx context.Context, d *net.Dialer, addr string, opt *pool.Options) (*Client, error) {
	pool, err := pool.New(opt, func() (net.Conn, error) {
		cn, err := d.DialContext(ctx, "tcp", addr)
		if err != nil {
			return nil, err
		}
		return wrapConn(cn), nil
	})
	if err != nil {
		return nil, err
	}
	return &Client{cns: pool}, nil
}

// Close closes all connections.
func (c *Client) Close() error {
	return c.cns.Close()
}

// Pipeline starts a pipeline.
func (c *Client) Pipeline() *Pipeline {
	return &Pipeline{c: c}
}

// Call executes a single command and returns a response.
func (c *Client) Call(req *Request) (*Response, error) {
	cn, err := c.cns.Get()
	if err != nil {
		return nil, err
	}
	pc := cn.(*protoConn)
	defer c.cns.Put(pc)

	if err := pc.w.WriteMsg(req); err != nil {
		return nil, err
	}
	if err := pc.w.Flush(); err != nil {
		return nil, err
	}

	res := new(Response)
	if err := pc.r.ReadMsg(res); err != nil {
		return nil, err
	}
	return res, nil
}

// Pipeline can execute commands.
type Pipeline struct {
	c *Client
	r []*Request
}

// Call adds a call to the pipeline.
func (p *Pipeline) Call(req *Request) {
	p.r = append(p.r, req)
}

// Reset resets the pipeline.
func (p *Pipeline) Reset() {
	p.r = p.r[:0]
}

// Exec executes the pipeline and returns responses.
func (p *Pipeline) Exec() ([]*Response, error) {
	cn, err := p.c.cns.Get()
	if err != nil {
		return nil, err
	}
	pc := cn.(*protoConn)
	defer p.c.cns.Put(pc)

	for _, req := range p.r {
		if err := pc.w.WriteMsg(req); err != nil {
			return nil, err
		}
	}

	if err := pc.w.Flush(); err != nil {
		return nil, err
	}

	rs := make([]*Response, 0, len(p.r))
	for range p.r {
		res := new(Response)
		if err := pc.r.ReadMsg(res); err != nil {
			return nil, err
		}
		rs = append(rs, res)
	}
	return rs, nil
}
