package quasizero_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/bsm/quasizero"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var subject *quasizero.Server
	var lis net.Listener
	var ctx = context.Background()

	BeforeEach(func() {
		var err error
		lis, err = net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())

		config := &quasizero.ServerConfig{Timeout: 100 * time.Millisecond}
		subject = quasizero.NewServer(commandMap, config)
		go func() {
			defer GinkgoRecover()
			Expect(subject.Serve(lis)).NotTo(HaveOccurred())
		}()
	})

	AfterEach(func() {
		Expect(lis.Close()).To(Succeed())
	})

	It("should handle commands", func() {
		c, err := quasizero.NewClient(ctx, lis.Addr().String(), nil)
		Expect(err).NotTo(HaveOccurred())
		defer c.Close()

		Expect(c.Call(&quasizero.Request{
			Code: 1,
		})).To(Equal(&quasizero.Response{Payload: []byte("PONG")}))

		Expect(c.Call(&quasizero.Request{
			Code:    2,
			Payload: []byte("HeLLo"),
		})).To(Equal(&quasizero.Response{Payload: []byte("HeLLo")}))
	})

	It("should handle pipelines", func() {
		c, err := quasizero.NewClient(ctx, lis.Addr().String(), nil)
		Expect(err).NotTo(HaveOccurred())
		defer c.Close()

		p := c.Pipeline()
		for i := 0; i < 147; i++ {
			p.Call(&quasizero.Request{Code: 1})
		}
		res, err := p.Exec()
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(HaveLen(147))

		for _, r := range res {
			Expect(r).To(Equal(&quasizero.Response{Payload: []byte("PONG")}))
		}
	})

	It("should handle multiple clients", func() {
		c1, err := quasizero.NewClient(ctx, lis.Addr().String(), nil)
		Expect(err).NotTo(HaveOccurred())
		defer c1.Close()

		c2, err := quasizero.NewClient(ctx, lis.Addr().String(), nil)
		Expect(err).NotTo(HaveOccurred())
		defer c2.Close()

		Expect(c1.Call(&quasizero.Request{
			Code: 1,
		})).To(Equal(&quasizero.Response{Payload: []byte("PONG")}))
		Expect(c1.Close()).To(Succeed())
		Expect(c2.Call(&quasizero.Request{
			Code: 1,
		})).To(Equal(&quasizero.Response{Payload: []byte("PONG")}))
	})

	It("should handle invalid commands", func() {
		c, err := quasizero.NewClient(ctx, lis.Addr().String(), nil)
		Expect(err).NotTo(HaveOccurred())
		defer c.Close()

		Expect(c.Call(&quasizero.Request{
			Code: 99,
		})).To(Equal(&quasizero.Response{ClientError: "unknown command code 99"}))
	})
})

// --------------------------------------------------------------------

func BenchmarkServer(b *testing.B) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	defer lis.Close()

	srv := quasizero.NewServer(commandMap, nil)
	go func() {
		if err := srv.Serve(lis); err != nil {
			b.Fatal(err)
		}
	}()

	clnt, err := quasizero.NewClient(context.Background(), lis.Addr().String(), nil)
	if err != nil {
		b.Fatal(err)
	}
	defer clnt.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := clnt.Call(&quasizero.Request{Code: 1})
		if err != nil {
			b.Fatal(err)
		} else if len(res.Payload) != 4 {
			b.Fatalf("expected PONG but got %s", res.Payload)
		}
		res.Release()
	}
}

func BenchmarkServer_Pipeline(b *testing.B) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	defer lis.Close()

	srv := quasizero.NewServer(commandMap, nil)
	go func() {
		if err := srv.Serve(lis); err != nil {
			b.Fatal(err)
		}
	}()

	clnt, err := quasizero.NewClient(context.Background(), lis.Addr().String(), nil)
	if err != nil {
		b.Fatal(err)
	}
	defer clnt.Close()

	p := clnt.Pipeline()
	b.ResetTimer()
	for i := 0; i < b.N; i += 10 {
		p.Reset()
		for j := 0; j < 10; j++ {
			p.Call(&quasizero.Request{Code: 1})
		}

		res, err := p.Exec()
		if err != nil {
			b.Fatal(err)
		} else if len(res) != 10 {
			b.Fatalf("expected 10xPONG but got %v", res)
		}
		res.Release()
	}
}
