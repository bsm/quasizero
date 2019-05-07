package quasizero_test

import (
	"testing"

	"github.com/bsm/quasizero"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "quasizero")
}

func pongHandler(_ *quasizero.Request) (*quasizero.Response, error) {
	return &quasizero.Response{Payload: []byte("PONG")}, nil
}

func echoHandler(req *quasizero.Request) (*quasizero.Response, error) {
	return &quasizero.Response{Payload: req.Payload}, nil
}

var commandMap = map[int32]quasizero.Handler{
	1: quasizero.HandlerFunc(pongHandler),
	2: quasizero.HandlerFunc(echoHandler),
}
