package static

import (
	"embed"
	"github.com/go-chi/chi/v5"
	fs2 "io/fs"
	"net/http"
)

//go:embed assets
var static embed.FS

func NewServeStaticHTTPHandler() http.Handler {
	router := chi.NewRouter()

	root, _ := fs2.Sub(static, "assets")
	fs := http.FileServer(http.FS(root))

	router.Handle("/*", fs)
	return router
}
