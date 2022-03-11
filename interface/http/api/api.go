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
	// 106e08fc-183d-4edc-93f9-391092ccfa97
	router.Route("/v1", func(r chi.Router) {
		r.Route(fmt.Sprintf("/uuid/{%s}", v1.UuidFieldName), func(r chi.Router) {
			r.Use(v1.UuidField)
			r.Get("/", apiV1.FindByUuid)
			r.Delete("/", apiV1.BanByUuid)
		})
		r.Route(fmt.Sprintf("/login/{%s}", v1.LoginFieldName), func(r chi.Router) {
			r.Use(v1.LoginField)
			r.Get("/", apiV1.FindByLogin)
			r.Put("/", apiV1.Add)
			r.Delete("/", apiV1.BanByLogin)
		})
		r.Route("/", func(r chi.Router) {
			r.Use(v1.PageFields)
			r.Get("/page", apiV1.Page)
		})
	})

	return router
}
