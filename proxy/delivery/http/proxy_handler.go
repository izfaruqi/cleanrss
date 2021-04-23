package http

import (
	"bytes"
	"cleanrss/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	neturl "net/url"
	"strings"
	"sync"
)

type proxyHandler struct {
	listeningAddress string

	hostMutex sync.Mutex
	host      string
}

func NewProxyHandler(httpServer *fiber.App, endpoint string, listeningAddress string) {
	proxyHandler := &proxyHandler{listeningAddress: listeningAddress}
	httpServer.Use(endpoint, proxyHandler.HandleRequests)
}

func (h *proxyHandler) HandleRequests(c *fiber.Ctx) error {
	csp := []byte("Content-Security-Policy")
	url := c.Query("u")
	if url == "" {
		return c.SendStatus(201)
	}
	if c.Query("i") != h.listeningAddress { // TODO: Make so that this ignores "http"/"https"
		h.hostMutex.Lock()
		fullUrl, err := neturl.Parse(url)
		if err != nil {
			return err
		}
		h.host = fullUrl.Scheme + "://" + fullUrl.Host
		h.hostMutex.Unlock()
	}
	if strings.HasPrefix(url, h.listeningAddress) {
		url = h.host + url[21:]
	}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(url)
	req.Header.SetMethod(c.Method())

	c.Request().Header.VisitAll(func(key, value []byte) {
		if !bytes.Equal(csp, key) {
			req.Header.SetBytesKV(key, value)
		}
	})

	// TODO: Make this configurable.
	req.Header.SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36")

	if c.Method() != "GET" || c.Method() != "HEAD" {
		req.SetBody(c.Body())
	}

	err := utils.FasthttpClient.Do(req, res)
	if err != nil {
		return err
	}

	res.Header.VisitAll(func(key, value []byte) {
		if !bytes.Equal(csp, key) {
			c.Response().Header.SetBytesKV(key, value)
		}
	})
	c.Response().Header.Del("X-Frame-Options")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	return c.Status(res.StatusCode()).Send(res.Body())
}
