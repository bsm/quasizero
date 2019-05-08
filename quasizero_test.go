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

func pongHandler(_ *quasizero.Request, res *quasizero.Response) error {
	res.SetString("PONG")
	return nil
}

func echoHandler(req *quasizero.Request, res *quasizero.Response) error {
	res.Set(req.Payload)
	return nil
}

var commandMap = map[int32]quasizero.Handler{
	1: quasizero.HandlerFunc(pongHandler),
	2: quasizero.HandlerFunc(echoHandler),
}
