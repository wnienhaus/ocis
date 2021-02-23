package grpc

import (
	"context"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/proto/v0"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Name                string
	Address             string
	Metadata            map[string]string
	Logger              log.Logger
	Context             context.Context
	Config              *config.Config
	Namespace           string
	Flags               []cli.Flag
	AuthRegistryHandler proto.RegistryAPIHandler
	AuthProviderHandler proto.ProviderAPIHandler
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Logger provides a function to set the logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// Name provides a name for the service.
func Name(val string) Option {
	return func(o *Options) {
		o.Name = val
	}
}

// Metadata associated with the server
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

// Address provides an address for the service.
func Address(val string) Option {
	return func(o *Options) {
		o.Address = val
	}
}

// Context provides a function to set the context option.
func Context(val context.Context) Option {
	return func(o *Options) {
		o.Context = val
	}
}

// Config provides a function to set the config option.
func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// Namespace provides a function to set the namespace option.
func Namespace(val string) Option {
	return func(o *Options) {
		o.Namespace = val
	}
}

// Flags provides a function to set the flags option.
func Flags(flags []cli.Flag) Option {
	return func(o *Options) {
		o.Flags = append(o.Flags, flags...)
	}
}

// AuthRegistryHandler provides a function to set the AuthRegistryHandler option.
func AuthRegistryHandler(h proto.RegistryAPIHandler) Option {
	return func(o *Options) {
		o.AuthRegistryHandler = h
	}
}

// AuthProviderHandler provides a function to set the AuthProviderHandler option.
func AuthProviderHandler(h proto.ProviderAPIHandler) Option {
	return func(o *Options) {
		o.AuthProviderHandler = h
	}
}
