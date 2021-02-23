package service

import (
	"context"

	"github.com/cs3org/reva/pkg/appctx"
	"github.com/cs3org/reva/pkg/auth"
	"github.com/cs3org/reva/pkg/errtypes"
	"github.com/cs3org/reva/pkg/rgrpc/status"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/proto/v0"
	"github.com/pkg/errors"
)

// NewAuthProvider returns a new instance of an authprovider handler
func NewAuthProvider(opts ...Option) (s proto.ProviderAPIHandler, err error) {
	options := newOptions(opts...)

	return &authProvider{
		log:     options.Logger,
		Config:  options.Config,
		authmgr: options.AuthManager,
	}, nil
}

// Service implements the AuthProviderService interface
type authProvider struct {
	log     log.Logger
	Config  *config.Config
	authmgr auth.Manager
}

func (s *authProvider) Authenticate(ctx context.Context, req *proto.AuthenticateRequest, res *proto.AuthenticateResponse) (err error) {
	log := appctx.GetLogger(ctx)
	username := req.ClientId
	password := req.ClientSecret

	u, err := s.authmgr.Authenticate(ctx, username, password)
	switch v := err.(type) {
	case nil:
		log.Info().Msgf("user %s authenticated", u.String())
		res.Status = status.NewOK(ctx)
		res.User = u
	case errtypes.InvalidCredentials:
		res.Status = status.NewPermissionDenied(ctx, v, "wrong password")
	case errtypes.NotFound:
		res.Status = status.NewNotFound(ctx, "unknown client id")
	default:
		err = errors.Wrap(err, "authsvc: error in Authenticate")
		res.Status = status.NewUnauthenticated(ctx, err, "error authenticating user")
	}
	return
}
