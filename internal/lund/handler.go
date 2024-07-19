package lund

import "github.com/valyala/fasthttp"

func MakeRequestHandler(state *State) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		server, err := state.GetNextServer()

		if err != nil {
			ctx.Response.SetBodyString(err.Error())
			return
		}

		ctx.Response.SetBodyString(server.URL)
	}
}
