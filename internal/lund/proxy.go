package lund

import "time"

type ProxyOptions struct {
	Interval         time.Duration
	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	DNSCacheDuration time.Duration
	Concurrency      int
}

func CreateHTTPClient(opts *ProxyOptions) {}
