package service

import (
	"github.com/go-chi/chi"
	"github.com/shportix/blob-svc/internal/data/pg"
	"github.com/shportix/blob-svc/internal/service/handlers"
	"github.com/shportix/blob-svc/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
			helpers.CtxBlobsQ(pg.NewBlobsQ(s.db)),
		),
	)
	r.Route("/blob", func(r chi.Router) {
		r.Post("/", handlers.CreateBlob)
		r.Get("/", handlers.GetBlobsList)
		r.Get("/{id}/", handlers.GetBlob)
		r.Delete("/{id}/", handlers.DeleteBlob)

	})

	return r
}
