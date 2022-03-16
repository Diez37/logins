package api

import (
	"fmt"
	"github.com/Diez37/logins/infrastructure/repository"
	v1 "github.com/Diez37/logins/interface/http/api/v1"
	"github.com/diez37/go-packages/log"
	"github.com/diez37/go-packages/router/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/trace"
)

func Router(repository repository.Repository, tracer trace.Tracer, logger log.Logger, validator *validator.Validate) chi.Router {
	apiV1 := v1.NewAPI(repository, tracer, logger, validator)

	router := chi.NewRouter()

	router.Route("/v1", func(r chi.Router) {
		r.Put("/login", apiV1.Add)

		r.Route(fmt.Sprintf("/uuid/{%s}", v1.UuidFieldName), func(r chi.Router) {
			r.Use(middlewares.NewUUID(logger, middlewares.WithName(v1.UuidFieldName), middlewares.WithUri(v1.UuidFieldName)).Middleware)
			r.Get("/", apiV1.FindByUuid)
			r.Delete("/", apiV1.BanByUuid)
			r.Post("/", apiV1.UpdateByUuid)
		})

		r.Route(fmt.Sprintf("/login/{%s}", v1.LoginFieldName), func(r chi.Router) {
			r.Use(middlewares.NewString(logger, middlewares.WithName(v1.LoginFieldName), middlewares.WithUri(v1.LoginFieldName)).Middleware)
			r.Get("/", apiV1.FindByLogin)
		})

		r.Get("/count", apiV1.Count)
		r.Route("/logins", func(r chi.Router) {
			r.Use(middlewares.NewUint64(
				logger,
				middlewares.WithName(v1.PageFieldName),
				middlewares.WithQuery(v1.PageFieldName),
				middlewares.WithHeader(v1.PageHeaderName),
				middlewares.WithDefault(v1.PageDefault),
			).Middleware)

			r.Use(middlewares.NewUint64(
				logger,
				middlewares.WithName(v1.LimitFieldName),
				middlewares.WithQuery(v1.LimitFieldName),
				middlewares.WithHeader(v1.PageHeaderName),
				middlewares.WithDefault(v1.LimitDefault),
			).Middleware)
			r.Get("/", apiV1.Page)
		})
	})

	return router
}
