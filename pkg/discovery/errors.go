package discovery

import "github.com/teakingwang/grpc-demo/pkg/errors"

var (
	ErrServiceNotFound = errors.FromError(errors.ErrNotFound)
)
