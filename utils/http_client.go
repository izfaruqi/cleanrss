package utils

import "github.com/valyala/fasthttp"

var FasthttpClient *fasthttp.Client

func HttpClientInit() {
	FasthttpClient = &fasthttp.Client{}
}