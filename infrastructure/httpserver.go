package infrastructure

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"sync"
)

type HTTPServer struct {
	*chi.Mux
}

func NewHTTPServer() HTTPServer {
	var server HTTPServer
	server.Mux = chi.NewRouter()
	server.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})
	return server
}

func (server HTTPServer) Listen(addr string, wg *sync.WaitGroup) error {
	defer wg.Done()
	return http.ListenAndServe(addr, server)
}
