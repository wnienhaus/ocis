package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// AuthBearerWithConfig applies cfg to the root flagset
func AuthBearerWithConfig(cfg *config.Config) []cli.Flag {
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
			Value:       "0.0.0.0:9149",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_AUTH_BEARER_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthBearer.DebugAddr,
		},

		// OIDC

		&cli.StringFlag{
			Name:        "oidc-issuer",
			Value:       "https://localhost:9200",
			Usage:       "OIDC issuer",
			EnvVars:     []string{"STORAGE_OIDC_ISSUER", "OCIS_URL"}, // STORAGE_OIDC_ISSUER takes precedence over OCIS_URL
			Destination: &cfg.Reva.OIDC.Issuer,
		},
		&cli.BoolFlag{
			Name:        "oidc-insecure",
			Value:       true,
			Usage:       "OIDC allow insecure communication",
			EnvVars:     []string{"STORAGE_OIDC_INSECURE"},
			Destination: &cfg.Reva.OIDC.Insecure,
		},
		&cli.StringFlag{
			Name: "oidc-id-claim",
			// preferred_username is a workaround
			// the user manager needs to take care of the sub to user metadata lookup, which ldap cannot do
			// TODO sub is stable and defined as unique.
			// AFAICT we want to use the account id from ocis-accounts
			// TODO add an ocis middleware to storage that changes the users opaqueid?
			// TODO add an ocis-accounts backed user manager
			Value:       "preferred_username",
			Usage:       "OIDC id claim",
			EnvVars:     []string{"STORAGE_OIDC_ID_CLAIM"},
			Destination: &cfg.Reva.OIDC.IDClaim,
		},
		&cli.StringFlag{
			Name:        "oidc-uid-claim",
			Value:       "",
			Usage:       "OIDC uid claim",
			EnvVars:     []string{"STORAGE_OIDC_UID_CLAIM"},
			Destination: &cfg.Reva.OIDC.UIDClaim,
		},
		&cli.StringFlag{
			Name:        "oidc-gid-claim",
			Value:       "",
			Usage:       "OIDC gid claim",
			EnvVars:     []string{"STORAGE_OIDC_GID_CLAIM"},
			Destination: &cfg.Reva.OIDC.GIDClaim,
		},

		// Services

		// AuthBearer

	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
