package command

import (
	"context"
	"os"
	"os/signal"
	"time"

	roidcm "github.com/cs3org/reva/pkg/auth/manager/oidc"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/server/grpc"
	svc "github.com/owncloud/ocis/storage/pkg/service/v0"
)

// AuthBearer is the entrypoint for the auth-bearer command.
func AuthBearer(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "auth-bearer",
		Usage: "Start authprovider for basic auth",
		Flags: flagset.AuthBearerWithConfig(cfg),
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
			config := map[string]interface{}{
				"issuer":    cfg.Reva.OIDC.Issuer,
				"insecure":  cfg.Reva.OIDC.Insecure,
				"id_claim":  cfg.Reva.OIDC.IDClaim,
				"uid_claim": cfg.Reva.OIDC.UIDClaim,
				"gid_claim": cfg.Reva.OIDC.GIDClaim,
			}
			authmgr, err := roidcm.New(config)
			if err != nil {
				logger.Fatal().Err(err).Msg("could not initialize oidc handler")
			}

			handler, err := svc.NewAuthProvider(
				svc.Logger(logger),
				svc.Config(cfg),
				svc.AuthManager(authmgr),
			)
			if err != nil {
				logger.Fatal().Err(err).Msg("could not initialize service handler")
			}

			// configure and run the grpc server
			{

				service := grpc.NewAuthProvider(
					grpc.Logger(logger),
					grpc.Context(ctx),
					grpc.Config(cfg),
					grpc.Namespace(cfg.Reva.AuthProvider.Namespace),
					grpc.Name(cfg.Reva.AuthProvider.Name),
					grpc.Metadata(map[string]string{"type": "bearer"}),
					// grpc.Metrics(metrics), // TODO metrics are part of the ocis-pkg grpc service, right?
					//grpc.Flags(flagset.RootWithConfig(config.New())),
					grpc.AuthProviderHandler(handler),
				)

				gr.Add(func() error {
					return service.Run()
				}, func(_ error) {
					cancel()
				})
			}

			// configure and run a http debug server for /readyz, /healthz, /metrics and /debug
			{
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Config(cfg),
					debug.Name(c.Command.Name+"-debug"),
					debug.Addr(cfg.Reva.AuthBearer.DebugAddr),
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
