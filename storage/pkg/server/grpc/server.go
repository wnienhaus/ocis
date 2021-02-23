package grpc

import (
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/storage/pkg/proto/v0"
	"github.com/owncloud/ocis/storage/pkg/version"
)

// NewAuthRegistry initializes a grpc service and server for a cs3 auth registry
func NewAuthRegistry(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Logger(options.Logger),
		grpc.Namespace(options.Namespace),
		grpc.Name(options.Name),
		grpc.Address(options.Address),
		grpc.Metadata(options.Metadata),
		grpc.Version(version.String),
		grpc.Context(options.Context),
		grpc.Flags(options.Flags...),
		grpc.Version(options.Config.Reva.AuthProvider.Version),
	)

	// add handlers

	if err := proto.RegisterRegistryAPIHandler(service.Server(), options.AuthRegistryHandler); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register service handler")
	}

	return service
}

// NewAuthProvider initializes a grpc service and server for a cs3 auth registry
func NewAuthProvider(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Logger(options.Logger),
		grpc.Namespace(options.Namespace),
		grpc.Metadata(options.Metadata),
		grpc.Name(options.Name),
		grpc.Version(version.String),
		grpc.Context(options.Context),
		grpc.Flags(options.Flags...),
		grpc.Version(options.Config.Reva.AuthProvider.Version),
	)

	// add handlers

	if err := proto.RegisterProviderAPIHandler(service.Server(), options.AuthProviderHandler); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register service handler")
	}

	return service
}
