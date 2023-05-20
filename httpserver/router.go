package httpserver

import (
	"monitoring/service/httpserver/middleware"
	"monitoring/service/httpserver/private"
	"monitoring/service/httpserver/public"
	"time"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func init() {

}

func newAppRouter() (*chi.Mux, error) {
	root := chi.NewRouter()
	root.Use(chimiddleware.Recoverer)
	root.Use(chimiddleware.RequestID)
	root.Use(chimiddleware.RealIP)
	root.Use(chimiddleware.Timeout(time.Minute * 5))
	root.Use(corsConfig().Handler)

	publicRouter := chi.NewRouter()
	publicRouter.Use(render.SetContentType(render.ContentTypeJSON))
	publicRouter.Group(func(r chi.Router) {
		r.Use(chimiddleware.DefaultCompress)
		r.Use(chimiddleware.RedirectSlashes)
		r.Use(middleware.Metric)
		r.Use(middleware.Logger)
		r.Mount("/", public.GetPublicRouter())
	})

	adminRouter := chi.NewRouter()
	adminRouter.Use(render.SetContentType(render.ContentTypeJSON))
	adminRouter.Group(func(r chi.Router) {
		r.Use(chimiddleware.DefaultCompress)
		r.Use(chimiddleware.RedirectSlashes)
		r.Use(middleware.Logger)
		r.Mount("/", private.GetPrivateRouter())
	})

	root.Mount(public.PublicRouterPathPrefix, publicRouter)
	root.Mount(private.PrivateRouterPathPrefix, adminRouter)
	return root, nil
}

func corsConfig() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400,
	})
}
