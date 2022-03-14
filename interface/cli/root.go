package cli

import (
	container2 "github.com/Diez37/logins/infrastructure/container"
	"github.com/Diez37/logins/interface/http"
	"github.com/diez37/go-packages/app"
	"github.com/diez37/go-packages/closer"
	"github.com/diez37/go-packages/configurator"
	bindFlags "github.com/diez37/go-packages/configurator/bind_flags"
	"github.com/diez37/go-packages/container"
	"github.com/diez37/go-packages/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

const (
	// AppName name of application
	AppName = "logins"
)

// NewRootCommand creating, configuration and return cobra.Command for root command
func NewRootCommand() (*cobra.Command, error) {
	container := container.GetContainer()

	if err := container2.AddProvide(container); err != nil {
		return nil, err
	}

	cmd := &cobra.Command{
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return container.Invoke(func(generalConfig *app.Config, configurator configurator.Configurator) {
				app.Configuration(generalConfig, configurator, app.WithAppName(AppName))
			})
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return container.Invoke(func(generalConfig *app.Config, logger log.Logger, closer closer.Closer, migrator *migrate.Migrate) error {
				logger.Infof("app: %s started", generalConfig.Name)
				logger.Infof("app: pid - %d", generalConfig.PID)

				if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
					return err
				}

				return http.Serve(closer.GetContext(), container, logger)
			})
		},
	}

	cmd, err := bindFlags.CobraCmd(container, cmd,
		bindFlags.HttpServer,
		bindFlags.Logger,
		bindFlags.Tracer,
		bindFlags.DataBase,
	)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
