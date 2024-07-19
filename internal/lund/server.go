package lund

import (
	"sync/atomic"

	"github.com/valyala/fasthttp"
)

type Server struct {
	URL   string
	Alive atomic.Bool

	// wait, what if we use PipelineClient?
	Client *fasthttp.Client
}
