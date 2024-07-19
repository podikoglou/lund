package lund

import (
	"github.com/valyala/fasthttp"
)

func MakeRequestHandler(state *State) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		server, err := state.GetNextServer()

		if err != nil {
			ctx.Response.SetBodyString(err.Error())
			return
		}

		// this is so smart
		req := &ctx.Request
		req.SetHost(server.GetHost())

		resp := &ctx.Response

		// perform request
		err = server.Client.Do(req, resp)

		// do we need to do this? or does the server do it?
		// probably the server, right?
		fasthttp.ReleaseRequest(req)

		if err != nil {
			ctx.Response.SetBodyString(err.Error())
		}
	}
}
