package vault

import (
	"context"

	vaultapi "github.com/hashicorp/vault/api"
)

// logical is an interface over the subset of vault API logical client methods
// used by this package. It allows easy stubbing in tests.
type logical interface {
	ReadWithContext(ctx context.Context, path string) (*vaultapi.Secret, error)
}

// apiLogical wraps *vaultapi.Logical to satisfy the logical interface.
type apiLogical struct {
	inner *vaultapi.Logical
}

func (a *apiLogical) ReadWithContext(ctx context.Context, path string) (*vaultapi.Secret, error) {
	return a.inner.ReadWithContext(ctx, path)
}
