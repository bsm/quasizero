package quasizero_test

import (
	"fmt"
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

func failingHandler(req *quasizero.Request, res *quasizero.Response) error {
	return fmt.Errorf("something went wrong")
}

var commandMap = map[int32]quasizero.Handler{
	1: quasizero.HandlerFunc(pongHandler),
	2: quasizero.HandlerFunc(echoHandler),
	3: quasizero.HandlerFunc(failingHandler),
}
