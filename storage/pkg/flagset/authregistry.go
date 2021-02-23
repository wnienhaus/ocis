package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// AuthRegistryWithConfig applies cfg to the root flagset
func AuthRegistryWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "name",
			Value:       "authregistry",
			Usage:       "TODO",
			EnvVars:     []string{"STORAGE_AUTH_REGISTRY_NAME"},
			Destination: &cfg.Reva.AuthRegistry.Name,
		}, &cli.StringFlag{
			Name:        "namespace",
			Value:       "com.owncloud",
			Usage:       "TODO",
			EnvVars:     []string{"STORAGE_AUTH_REGISTRY_NAMESPACE"},
			Destination: &cfg.Reva.AuthRegistry.Namespace,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9800",
			Usage:       "Address to bind auth registry service",
			EnvVars:     []string{"STORAGE_AUTH_REGISTRY_ADDR"},
			Destination: &cfg.Reva.AuthRegistry.GRPCAddr,
		},
		// no flags about ports or anything, the gateway service will discover the auth registry using micro
		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9801",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_AUTH_REGISTRY_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthRegistry.DebugAddr,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
