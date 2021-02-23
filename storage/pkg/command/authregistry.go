package command

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/server/grpc"
	svc "github.com/owncloud/ocis/storage/pkg/service/v0"
)

// AuthRegistry is the entrypoint for the auth-registry command.
func AuthRegistry(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth-registry",
		Usage: "Start cs3 authregistry",
		Flags: flagset.AuthRegistryWithConfig(cfg),
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			// TODO add tracing

			var (
				gr          = run.Group{}
				ctx, cancel = context.WithCancel(context.Background())
				//mtrcs       = metrics.New()
			)

			defer cancel()

			// first initialize a service implementation
			handler, err := svc.NewAuthRegistry(svc.Logger(logger), svc.Config(cfg))
			if err != nil {
				logger.Fatal().Err(err).Msg("could not initialize service handler")
			}

			// configure and run the grpc server
			{

				service := grpc.NewAuthRegistry(
					grpc.Logger(logger),
					grpc.Context(ctx),
					grpc.Config(cfg),
					grpc.Namespace(cfg.Reva.AuthRegistry.Namespace),
					grpc.Name(cfg.Reva.AuthRegistry.Name),
					grpc.Address(cfg.Reva.AuthRegistry.GRPCAddr),
					// grpc.Metrics(metrics), // TODO metrics are part of the ocis-pkg grpc service, right?
					//grpc.Flags(flagset.RootWithConfig(config.New())),
					grpc.AuthRegistryHandler(handler),
				)

				gr.Add(func() error {
					return service.Run()
				}, func(err error) {
					if err != nil {
						logger.Error().Err(err).Msg("interrupted service with error")
					}
					cancel()
				})
			}

			// configure and run a http debug server for /readyz, /healthz, /metrics and /debug
			{
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Config(cfg),
					debug.Name(c.Command.Name+"-debug"),
					debug.Addr(cfg.Reva.AuthRegistry.DebugAddr),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("transport", "debug").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(_ error) {
					ctx, timeout := context.WithTimeout(ctx, 5*time.Second)

					defer timeout()
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						logger.Info().
							Err(err).
							Str("transport", "debug").
							Msg("Failed to shutdown server")
					} else {
						logger.Info().
							Str("transport", "debug").
							Msg("Shutting down server")
					}
				})
			}

			// capture cli interrupts
			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
					cancel()
				})
			}

			// finally, bring it all up
			return gr.Run()
		},
	}
}
