// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/storage/pkg/command"
	svcconfig "github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
)

// StorageAuthRegistryCommand is the entrypoint for the reva-auth-basic command.
func StorageAuthRegistryCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-auth-registry",
		Usage:    "Start storage auth-registry service",
		Category: "Extensions",
		Flags:    flagset.AuthRegistryWithConfig(cfg.Storage),
		Action: func(c *cli.Context) error {
			origCmd := command.AuthRegistry(configureStorageAuthRegistry(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureStorageAuthRegistry(cfg *config.Config) *svcconfig.Config {
	cfg.Storage.Log.Level = cfg.Log.Level
	cfg.Storage.Log.Pretty = cfg.Log.Pretty
	cfg.Storage.Log.Color = cfg.Log.Color

	if cfg.Tracing.Enabled {
		cfg.Storage.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Storage.Tracing.Type = cfg.Tracing.Type
		cfg.Storage.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Storage.Tracing.Collector = cfg.Tracing.Collector
	}

	return cfg.Storage
}

func init() {
	register.AddCommand(StorageAuthRegistryCommand)
}
