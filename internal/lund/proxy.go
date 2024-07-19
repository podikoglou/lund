package lund

import (
	"time"

	"github.com/valyala/fasthttp"
)

type ProxyOptions struct {
	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	DNSCacheDuration time.Duration
	Concurrency      int
}

func CreateHTTPClient(opt *ProxyOptions) *fasthttp.Client {
	return &fasthttp.Client{
		ReadTimeout:                   opt.ReadTimeout,
		WriteTimeout:                  opt.WriteTimeout,
		NoDefaultUserAgentHeader:      false,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      opt.Concurrency,
			DNSCacheDuration: opt.DNSCacheDuration,
		}).Dial,
	}
}
