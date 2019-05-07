package quasizero

import (
	"fmt"
	"net"
	"time"
)

// ServerConfig holds the server configuration
type ServerConfig struct {
	// Timeout represents the per-request socket read/write timeout.
	// Default: 0 (disabled)
	Timeout time.Duration

	// IdleTimeout forces servers to close idle connection once timeout is reached.
	// Default: 0 (disabled)
	IdleTimeout time.Duration

	// If non-zero, use SO_KEEPALIVE to send TCP ACKs to clients in absence
	// of communication. This is useful for two reasons:
	// 1) Detect dead peers.
	// 2) Take the connection alive from the point of view of network
	//    equipment in the middle.
	// On Linux, the specified value (in seconds) is the period used to send ACKs.
	// Note that to close the connection the double of the time is needed.
	// On other kernels the period depends on the kernel configuration.
	// Default: 0 (disabled)
	TCPKeepAlive time.Duration

	// OnError is called on client errors. Use for verbose logging.
	OnError func(error)
}

func (c *ServerConfig) norm() *ServerConfig {
	if c != nil {
		return c
	}
	return new(ServerConfig)
}

// --------------------------------------------------------------------

// Server instances can handle client requests.
type Server struct {
	hs map[int32]Handler
	cf *ServerConfig
}

// NewServer creates a new server instance.
func NewServer(commands map[int32]Handler, cfg *ServerConfig) *Server {
	return &Server{hs: commands, cf: cfg.norm()}
}

// Serve accepts incoming connections on a listener, creating a
// new service goroutine for each.
func (s *Server) Serve(lis net.Listener) error {
	for {
		cn, err := lis.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() || ne.Timeout() {
				continue
			}
			return nil
		}

		if ka := s.cf.TCPKeepAlive; ka > 0 {
			if tc, ok := cn.(*net.TCPConn); ok {
				tc.SetKeepAlive(true)
				tc.SetKeepAlivePeriod(ka)
			}
		}
		go s.serveClient(wrapConn(cn))
	}
}

// Starts a new session, serving client
func (s *Server) serveClient(c *protoConn) {
	// close client on exit
	defer c.Close()

	for {
		// set deadline
		if d := s.cf.Timeout; d > 0 {
			c.SetDeadline(time.Now().Add(d))
		}

		// perform pipeline
		if err := s.pipeline(c); err != nil {
			if s.cf.OnError != nil {
				s.cf.OnError(err)
			}
			return
		}
	}
}

func (s *Server) pipeline(c *protoConn) error {
	for more := true; more; more = c.r.Buffered() > 0 {
		req := new(Request)
		if err := c.r.ReadMsg(req); err != nil {
			return err
		}

		res, err := s.process(req)
		if err != nil {
			return err
		}

		if err := c.w.WriteMsg(res); err != nil {
			return err
		}
	}
	return c.w.Flush()
}

func (s *Server) process(req *Request) (*Response, error) {
	if handler, ok := s.hs[req.Code]; ok {
		return handler.ServeQZ(req)
	}

	return &Response{
		ClientError: fmt.Sprintf("unknown command code %d", req.Code),
	}, nil
}
