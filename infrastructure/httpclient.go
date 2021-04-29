package infrastructure

import (
	"errors"
	"net/http"
)

func NewHTTPClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 20 { // Max redirect: 20
				return errors.New("too many redirects")
			}
			return nil
		},
	}
}
