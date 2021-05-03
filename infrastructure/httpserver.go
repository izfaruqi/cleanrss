package infrastructure

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"sync"
	"time"
)

type HTTPServer struct {
	*chi.Mux
	srv *http.Server
}

func NewHTTPServer() HTTPServer {
	var server HTTPServer
	server.srv = &http.Server{}
	server.Mux = chi.NewRouter()
	server.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})
	return server
}

func (server HTTPServer) Listen(addr string, wg *sync.WaitGroup, ctx context.Context) (err error) {
	defer wg.Done()
	server.srv.Addr = addr
	server.srv.Handler = server
	go func() {
		if err = server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Failed to start server ")
			log.Println(err)
		}
	}()
	<-ctx.Done()
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() { cancel() }()
	if err = server.srv.Shutdown(ctxShutdown); err != nil && err != http.ErrServerClosed {
		log.Println("Server shutdown failed.")
		log.Println(err)
	}
	log.Println("Server shutdown successfully.")
	if err == http.ErrServerClosed {
		err = nil
	}
	return
}
