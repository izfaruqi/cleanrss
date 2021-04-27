package infrastructure

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"sync"
)

type HTTPChiServer struct {
	*chi.Mux
}

func NewHTTPChiServer() HTTPChiServer {
	var server HTTPChiServer
	server.Mux = chi.NewRouter()
	server.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})
	return server
}

func (server HTTPChiServer) Listen(addr string, wg *sync.WaitGroup) error {
	defer wg.Done()
	return http.ListenAndServe(addr, server)
}
