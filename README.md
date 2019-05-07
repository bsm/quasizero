# QuasiZero

[![GoDoc](https://godoc.org/github.com/bsm/quasizero?status.svg)](https://godoc.org/github.com/bsm/quasizero)
[![Build Status](https://travis-ci.org/bsm/quasizero.svg)](https://travis-ci.org/bsm/quasizero)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A [Go](https://golang.org/) TCP server (and client) implementation, optimised for low latency and pipelined throughput.

## Usage

Server:

```go
// define a handler
echoHandler := quasizero.HandlerFunc(func(req *quasizero.Request) (*quasizero.Response, error) {
  return &quasizero.Response{Payload: req.Payload}, nil
})

// init a server
srv := quasizero.NewServer(map[int32]quasizero.Handler{
  1: echoHandler,
}, nil)

// listen and serve
lis, err := net.Listen("tcp", ":11111")
if err != nil {
  // handle error ...
}
defer lis.Close()

if err := srv.Serve(lis); err != nil {
  // handle error ...
}
```

Client:

```go
client, err := quasizero.NewClient(context.TODO(), "10.0.0.1:11111", nil)
if err != nil {
  // handle error ...
}
defer client.Close()

// send an echo request
res, err := client.Call(&quasizero.Request{Code: 1, Payload: []byte("hello")})
if err != nil {
  // handle error ...
}
fmt.Printf("server responded to ECHO with %q\n", res.Payload)
```

## Documentation

Please see the [API documentation](https://godoc.org/github.com/bsm/quasizero) for
package and API descriptions and examples.
