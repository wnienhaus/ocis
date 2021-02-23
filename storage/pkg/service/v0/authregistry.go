package service

import (
	"context"

	"github.com/asim/go-micro/v3/registry"
	"github.com/cs3org/reva/pkg/rgrpc/status"
	"github.com/owncloud/ocis/ocis-pkg/log"
	oregistry "github.com/owncloud/ocis/ocis-pkg/registry"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/proto/v0"
)

// NewAuthRegistry returns a new instance of an authregistry handler
func NewAuthRegistry(opts ...Option) (s proto.RegistryAPIHandler, err error) {
	options := newOptions(opts...)

	return &authRegistry{
		log:    options.Logger,
		Config: options.Config,
		r:      oregistry.GetRegistry(),
	}, nil
}

// Service implements the AuthRegistryService interface
type authRegistry struct {
	log    log.Logger
	Config *config.Config
	r      registry.Registry
}

func (s *authRegistry) ListAuthProviders(ctx context.Context, req *proto.ListAuthProvidersRequest, res *proto.ListAuthProvidersResponse) (err error) {
	services, err := s.r.GetService("com.owncloud.authprovider") // this will list all auth providers

	if err != nil {
		res.Status = status.NewInternal(ctx, err, "error getting list of auth providers")
		return nil
	}

	// TODO filter nodes
	if res.Providers == nil {
		res.Providers = []*proto.ProviderInfo{}
	}
	for i := range services {
		for j := range services[i].Nodes {

			if _, ok := services[i].Nodes[j].Metadata["type"]; !ok {
				continue
			}
			pi := proto.ProviderInfo{
				Address:      services[i].Nodes[j].Address, // TODO make cs3 api aware of multiple nodes, for now we always return the first one
				ProviderType: services[i].Nodes[j].Metadata["type"],
			}
			if d, ok := services[i].Nodes[j].Metadata["description"]; ok {
				pi.Description = d
			}

			res.Providers = append(res.Providers, &pi)
		}
	}

	res.Status = status.NewOK(ctx)
	return
}

func (s *authRegistry) GetAuthProvider(ctx context.Context, req *proto.GetAuthProviderRequest, res *proto.GetAuthProviderResponse) (err error) {
	services, err := s.r.GetService("com.owncloud.authprovider") // this will list all auth providers

	if err != nil {
		res.Status = status.NewInternal(ctx, err, "error getting list of auth providers")
		return nil
	}

	for i := range services {
		for j := range services[i].Nodes {

			if t, ok := services[i].Nodes[j].Metadata["type"]; !ok || t != req.Type {
				continue
			}

			res.Provider = &proto.ProviderInfo{
				ProviderType: services[i].Nodes[j].Metadata["type"],
				Address:      services[i].Nodes[j].Address, // TODO make cs3 api aware of multiple nodes, for now we always return the first one
			}
			if d, ok := services[i].Nodes[j].Metadata["description"]; ok {
				res.Provider.Description = d
			}
			res.Status = status.NewOK(ctx)
			return nil
		}
	}

	res.Status = status.NewNotFound(ctx, "provider not found")
	return nil
}
