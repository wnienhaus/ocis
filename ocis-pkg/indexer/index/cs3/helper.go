package cs3

import (
	"context"
	"fmt"
	"path"

	"github.com/owncloud/ocis/accounts/pkg/proto/v0"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
)

func deleteIndexRoot(ctx context.Context, storageProvider proto.ProviderAPIService, indexRootDir string) error {
	res, err := storageProvider.Delete(ctx, &proto.DeleteRequest{
		Ref: &proto.Reference{
			Spec: &proto.Reference_Path{Path: path.Join("/meta", indexRootDir)},
		},
	})
	if err != nil {
		return err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return fmt.Errorf("error deleting index root dir: %v", indexRootDir)
	}

	return nil
}
