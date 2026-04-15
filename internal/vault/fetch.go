package vault

import (
	"context"
	"fmt"
)

// FetchAtRef retrieves a KV v2 secret at the version described by ref.
// If ref.IsLatest is true the latest version is fetched; otherwise the
// specific version number is requested via the "version" query parameter.
func FetchAtRef(ctx context.Context, c *Client, mount string, ref VersionRef) (map[string]string, error) {
	if c == nil {
		return nil, fmt.Errorf("vault client must not be nil")
	}

	var kvPath string
	if ref.IsLatest {
		kvPath = fmt.Sprintf("%s/data/%s", mount, ref.Path)
	} else {
		kvPath = fmt.Sprintf("%s/data/%s", mount, ref.Path)
	}

	var secret secretReader
	var err error

	if ref.IsLatest {
		secret, err = c.Logical().ReadWithContext(ctx, kvPath)
	} else {
		secret, err = c.Logical().ReadWithDataWithContext(ctx, kvPath, map[string][]string{
			"version": {fmt.Sprintf("%d", ref.Version)},
		})
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", ref, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found: %s", ref)
	}

	return extractStringMap(secret)
}

// FetchPairAtRefs fetches two secret versions described by refA and refB
// and returns their data maps ready for diffing.
func FetchPairAtRefs(
	ctx context.Context,
	c *Client,
	mount string,
	refA, refB VersionRef,
) (map[string]string, map[string]string, error) {
	a, err := FetchAtRef(ctx, c, mount, refA)
	if err != nil {
		return nil, nil, fmt.Errorf("fetching ref A (%s): %w", refA, err)
	}

	b, err := FetchAtRef(ctx, c, mount, refB)
	if err != nil {
		return nil, nil, fmt.Errorf("fetching ref B (%s): %w", refB, err)
	}

	return a, b, nil
}
