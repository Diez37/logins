package api

import (
	"fmt"
	"github.com/Diez37/logins/infrastructure/repository"
	v1 "github.com/Diez37/logins/interface/http/api/v1"
	"github.com/diez37/go-packages/log"
	"github.com/diez37/go-packages/server/http/helpers"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/trace"
)

func Router(repository repository.Repository, tracer trace.Tracer, errorHelper *helpers.Error, logger log.Logger) chi.Router {
	apiV1 := v1.NewAPI(repository, tracer, errorHelper, logger)

	router := chi.NewRouter()

	router.Route("/v1", func(r chi.Router) {
		r.Put("/login", apiV1.Add)

		r.Route(fmt.Sprintf("/uuid/{%s}", v1.UuidFieldName), func(r chi.Router) {
			r.Use(v1.NewUuid(logger).Middleware)
			r.Get("/", apiV1.FindByUuid)
			r.Delete("/", apiV1.BanByUuid)
			r.Post("/", apiV1.UpdateByUuid)
		})

		r.Route(fmt.Sprintf("/login/{%s}", v1.LoginFieldName), func(r chi.Router) {
			r.Use(v1.LoginField)
			r.Get("/", apiV1.FindByLogin)
		})

		r.Get("/count", apiV1.Count)
		r.Route("/logins", func(r chi.Router) {
			r.Use(v1.PageField)
			r.Use(v1.LimitField)
			r.Get("/", apiV1.Page)
		})
	})

	return router
}
