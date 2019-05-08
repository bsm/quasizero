package quasizero_test

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/bsm/quasizero"
)

func Example() {
	// start a TCP listener
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	// define command map
	cmds := map[int32]quasizero.Handler{
		// ECHO
		1: quasizero.HandlerFunc(func(req *quasizero.Request, res *quasizero.Response) error {
			res.Set(req.Payload)
			return nil
		}),

		// SHUTDOWN
		9: quasizero.HandlerFunc(func(req *quasizero.Request, res *quasizero.Response) error {
			go func() {
				time.Sleep(time.Second)
				_ = lis.Close()
			}()
			res.SetString("OK")
			return nil
		}),
	}

	// start serving (in background)
	go func() {
		srv := quasizero.NewServer(cmds, nil)
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// connect client
	clnt, err := quasizero.NewClient(context.TODO(), lis.Addr().String(), nil)
	if err != nil {
		panic(err)
	}
	defer clnt.Close()

	// send an echo request
	res1, err := clnt.Call(&quasizero.Request{Code: 1, Payload: []byte("hello")})
	if err != nil {
		panic(err)
	}
	fmt.Printf("server responded to ECHO with %q\n", res1.Payload)

	// send a shutdown request
	res2, err := clnt.Call(&quasizero.Request{Code: 9})
	if err != nil {
		panic(err)
	}
	fmt.Printf("server responded to SHUTDOWN with %q\n", res2.Payload)

	// Output:
	// server responded to ECHO with "hello"
	// server responded to SHUTDOWN with "OK"
}
