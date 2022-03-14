package http

import (
	"context"
	"github.com/Diez37/logins/infrastructure/repository"
	"github.com/Diez37/logins/interface/http/api"
	"github.com/diez37/go-packages/container"
	"github.com/diez37/go-packages/log"
	httpServer "github.com/diez37/go-packages/server/http"
	"github.com/diez37/go-packages/server/http/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"net/http"
)

// Serve configuration and running http server
func Serve(ctx context.Context, container container.Container, logger log.Logger) error {
	return container.Invoke(func(
		server *http.Server,
		config *httpServer.Config,
		repository repository.Repository,
		tracer trace.Tracer,
		errorHelper *helpers.Error,
		router chi.Router,
		validator *validator.Validate,
	) error {

		router.Mount("/api", api.Router(
			repository,
			tracer,
			errorHelper,
			logger,
			validator,
		))

		errGroup := &errgroup.Group{}

		errGroup.Go(func() error {
			logger.Infof("http server: started")
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				return err
			}

			return nil
		})

		errGroup.Go(func() error {
			<-ctx.Done()

			logger.Infof("http server: shutdown")

			ctxTimeout, cancelFnc := context.WithTimeout(context.Background(), config.ShutdownTimeout)
			defer cancelFnc()

			return server.Shutdown(ctxTimeout)
		})

		return errGroup.Wait()
	})
}
