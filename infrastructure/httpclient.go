package infrastructure

import "github.com/valyala/fasthttp"

func NewHTTPClient() *fasthttp.Client {
	return &fasthttp.Client{}
}
