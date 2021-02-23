package grpc

import (
	"context"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/metrics"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Name                   string
	Logger                 log.Logger
	Context                context.Context
	Config                 *config.Config
	Metrics                *metrics.Metrics
	Flags                  []cli.Flag
	AccountsServiceHandler proto.AccountsServiceHandler
	GroupsServiceHandler   proto.GroupsServiceHandler
	IndexServiceHandler    proto.IndexServiceHandler
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Name provides a name for the service.
func Name(val string) Option {
	return func(o *Options) {
		o.Name = val
	}
}

// Logger provides a function to set the logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
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

// Metrics provides a function to set the metrics option.
func Metrics(val *metrics.Metrics) Option {
	return func(o *Options) {
		o.Metrics = val
	}
}

// Flags provides a function to set the flags option.
func Flags(val []cli.Flag) Option {
	return func(o *Options) {
		o.Flags = append(o.Flags, val...)
	}
}

// AccountsServiceHandler provides a function to set the AccountsServiceHandler option.
func AccountsServiceHandler(val proto.AccountsServiceHandler) Option {
	return func(o *Options) {
		o.AccountsServiceHandler = val
	}
}

// GroupsServiceHandler provides a function to set the GroupsServiceHandler option.
func GroupsServiceHandler(val proto.GroupsServiceHandler) Option {
	return func(o *Options) {
		o.GroupsServiceHandler = val
	}
}

// IndexServiceHandler provides a function to set the handler option.
func IndexServiceHandler(val proto.IndexServiceHandler) Option {
	return func(o *Options) {
		o.IndexServiceHandler = val
	}
}
