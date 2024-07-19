package lund

import (
	"net/url"
	"sync/atomic"

	"github.com/valyala/fasthttp"
)

type Server struct {
	URL   string
	Alive atomic.Bool

	// wait, what if we use PipelineClient?
	Client *fasthttp.Client
}

func (s *Server) GetHost() string {
	// we can safely ignore the error, because we validate the URL
	// when we're inserting it
	url, _ := url.Parse(s.URL)

	return url.Host
}
