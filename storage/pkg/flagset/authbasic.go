package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// AuthBasicWithConfig applies cfg to the root flagset
func AuthBasicWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "name",
			Value:       "authprovider",
			Usage:       "TODO",
			EnvVars:     []string{"STORAGE_AUTH_PROVIDER_NAME"},
			Destination: &cfg.Reva.AuthProvider.Name,
		}, &cli.StringFlag{
			Name:        "namespace",
			Value:       "com.owncloud",
			Usage:       "TODO",
			EnvVars:     []string{"STORAGE_AUTH_PROVIDER_NAMESPACE"},
			Destination: &cfg.Reva.AuthProvider.Namespace,
		},

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9147",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_AUTH_BASIC_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthBasic.DebugAddr,
		},

		// Auth

		&cli.StringFlag{
			Name:        "auth-driver",
			Value:       "ldap",
			Usage:       "auth driver: 'demo', 'json' or 'ldap'",
			EnvVars:     []string{"STORAGE_AUTH_DRIVER"},
			Destination: &cfg.Reva.AuthProvider.Driver,
		},
		&cli.StringFlag{
			Name:        "auth-json",
			Value:       "",
			Usage:       "Path to users.json file",
			EnvVars:     []string{"STORAGE_AUTH_JSON"},
			Destination: &cfg.Reva.AuthProvider.JSON,
		},
		// Gateway

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the storage gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_ENDPOINT"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, LDAPWithConfig(cfg)...)

	return flags
}
