package web_ext

import (
	"cleanrss/domain"
	"io"
	"net/http"
)

type webExtCleanerRepository struct {
	client *http.Client
}

func NewWebExtCleanerRepository(client *http.Client) domain.WebExtCleanerRepository {
	return webExtCleanerRepository{client: client}
}

func (w webExtCleanerRepository) GetRawPage(url string, mobileUA bool) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if mobileUA {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.152 Mobile Safari/537.36")
	}

	res, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
