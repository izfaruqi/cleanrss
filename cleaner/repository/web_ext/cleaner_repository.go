package web_ext

import (
	"bytes"
	"cleanrss/domain"
	"github.com/valyala/fasthttp"
	"io"
)

type webExtCleanerRepository struct {
	client *fasthttp.Client
}

func NewWebExtCleanerRepository(client *fasthttp.Client) domain.WebExtCleanerRepository {
	return webExtCleanerRepository{client: client}
}

func (w webExtCleanerRepository) GetRawPage(url string, mobileUA bool) (io.Reader, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	req.SetRequestURI(url)
	if mobileUA {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.152 Mobile Safari/537.36")
	}

	err := w.client.DoRedirects(req, res, 20)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(res.Body()), nil
}
