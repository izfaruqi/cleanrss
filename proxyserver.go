package main

import (
	"bytes"
	"cleanrss/utils"
	"log"
	neturl "net/url"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var (
	hostMutex sync.Mutex
	host      string
)

func ProxyServerInit(listenOn string, wg *sync.WaitGroup) {
	defer wg.Done()
	csp := []byte("Content-Security-Policy")
	server := fiber.New(fiber.Config{DisableStartupMessage: true})

	listenOnFull := "http://" + listenOn
	server.Use("/proxy", func(c *fiber.Ctx) error {
		url := c.Query("u")
		if url == "" {
			return c.SendStatus(201)
		}
		if c.Query("i") != listenOnFull {
			hostMutex.Lock()
			fullUrl, err := neturl.Parse(url)
			if err != nil {
				return err
			}
			host = fullUrl.Scheme + "://" + fullUrl.Host
			hostMutex.Unlock()
		}
		if strings.HasPrefix(url, listenOnFull) {
			url = host + url[21:]
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

		utils.FasthttpClient.Do(req, res)

		res.Header.VisitAll(func(key, value []byte) {
			if !bytes.Equal(csp, key) {
				c.Response().Header.SetBytesKV(key, value)
			}
		})
		c.Response().Header.Del("X-Frame-Options")
		c.Response().Header.Set("Access-Control-Allow-Origin", "*")

		return c.Status(res.StatusCode()).Send(res.Body())
	})

	log.Println("Proxy server will listen on " + listenOnFull)
	server.Listen(listenOn)
}
