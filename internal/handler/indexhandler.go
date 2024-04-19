package handler

import (
	"github.com/run-bigpig/jb-active/internal/response"
	"github.com/valyala/fasthttp"
)

func IndexHandler(ctx *fasthttp.RequestCtx) {
	response.Html(ctx, "index.html")
}
