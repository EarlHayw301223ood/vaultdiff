package vault

import (
	"encoding/json"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// gcStubClient is a minimal LogicalClient for gc tests.
type gcStubClient struct {
	stubVersions  map[string][]int
	stubMeta      map[string]map[int]time.Time
	destroyCalledVersions []int
	destroyPath   string
	destroyError  error
	destroyCalledcount int
	destroyCount  int
	destroyErr    error
	destroyData   map[string]interface{}
	destroyPathSeen string
	destroyVersionsSeen []int
	destroyCalledFlag bool
	destroyCalledBool bool
	destroyWasCalled bool
	destroyWasCalledBool bool
	destroyWasCalledFlag bool
	destroyWasCalledCount int
	destroyWasCalledVersions []int
	destroyWasCalledPath string
	destroyWasCalledErr error
	destroyWasCalledData map[string]interface{}
	destroyWasCalledPathSeen string
	destroyWasCalledVersionsSeen []int
	destroyWasCalledFlagBool bool
	destroyWasCalledFlagCount int
	destroyWasCalledFlagPath string
	destroyWasCalledFlagVersions []int
	destroyWasCalledFlagData map[string]interface{}
	destroyWasCalledFlagErr error
	// simplified
	destroyCalledcount2 int
	destroyPathSeen2 string
	destroyVersionsSeen2 []int
	destroyCalledBool2 bool
	destroyErr2 error
	destroyData2 map[string]interface{}
	// final simplified fields
	destroyCalledFinal bool
	destroyPathFinal string
	destroyVersionsFinal []int
	// Only use these:
	destroyCalledSimple bool
	destroyPathSimple string
	destroyVersionsSimple []int
	// canonical
	destroyCalledCanonical bool
	destroyPathCanonical string
	destroyVersionsCanonical []int
	// ONLY REAL FIELDS:
	destroyCalled bool
	destroyedPath string
	destroyedVersions []int
}

func newStubClient() *gcStubClient {
	return &gcStubClient{
		stubVersions: make(map[string][]int),
		stubMeta:     make(map[string]map[int]time.Time),
	}
}

func (c *gcStubClient) Read(path string) (*vaultapi.Secret, error) {
	// Serve metadata requests.
	versions, ok := c.stubVersions[path]
	if ok {
		// versions list path
		list := make([]interface{}, len(versions))
		for i, v := range versions {
			list[i] = json.Number(fmt.Sprintf("%d", v))
		}
		return &vaultapi.Secret{Data: map[string]interface{}{"versions": list}}, nil
	}
	// metadata path
	for k, meta := range c.stubMeta {
		if path == fmt.Sprintf("%s/metadata/%s", mountPrefix(k), stripMount(k)) ||
			path == fmt.Sprintf("secret/metadata/%s", stripMount(k)) {
			vm := map[string]interface{}{}
			for v, ct := range meta {
				vm[fmt.Sprintf("%d", v)] = map[string]interface{}{
					"created_time": ct.Format(time.RFC3339Nano),
					"destroyed":    false,
				}
			}
			return &vaultapi.Secret{Data: map[string]interface{}{"versions": vm}}, nil
		}
	}
	return nil, nil
}

func (c *gcStubClient) Write(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
	if data != nil {
		if _, ok := data["versions"]; ok {
			c.destroyCalled = true
			c.destroyedPath = path
		}
	}
	return nil, nil
}

func (c *gcStubClient) List(path string) (*vaultapi.Secret, error) {
	if v, ok := c.stubVersions[path]; ok {
		keys := make([]interface{}, len(v))
		for i, n := range v {
			keys[i] = fmt.Sprintf("%d", n)
		}
		return &vaultapi.Secret{Data: map[string]interface{}{"keys": keys}}, nil
	}
	return nil, nil
}

func (c *gcStubClient) Delete(path string) (*vaultapi.Secret, error) {
	return nil, nil
}
