package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

// cloneTestClient is a minimal LogicalClient for Clone tests.
type cloneTestClient struct {
	readFunc  func(path string) (*vaultapi.Secret, error)
	writeFunc func(path string, data map[string]interface{}) (*vaultapi.Secret, error)
	listFunc  func(path string) (*vaultapi.Secret, error)
}

func (c *cloneTestClient) Read(path string) (*vaultapi.Secret, error) {
	return c.readFunc(path)
}
func (c *cloneTestClient) Write(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
	return c.writeFunc(path, data)
}
func (c *cloneTestClient) List(path string) (*vaultapi.Secret, error) {
	return c.listFunc(path)
}

func TestClone_EmptySourcePath(t *testing.T) {
	_, err := Clone(nil, "", "secret/dest", "latest", CloneOptions{})
	if err == nil {
		t.Fatal("expected error for empty sourcePath")
	}
}

func TestClone_EmptyDestPath(t *testing.T) {
	_, err := Clone(nil, "secret/src", "", "latest", CloneOptions{})
	if err == nil {
		t.Fatal("expected error for empty destPath")
	}
}

func TestClone_SkipsWhenDestExists(t *testing.T) {
	client := &cloneTestClient{
		readFunc: func(path string) (*vaultapi.Secret, error) {
			return makeVersionsSecret(3), nil
		},
		writeFunc: func(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
			t.Fatal("write should not be called when dest exists and overwrite is false")
			return nil, nil
		},
		listFunc: func(path string) (*vaultapi.Secret, error) { return nil, nil },
	}

	result, err := Clone(client, "secret/src", "secret/dest", "latest", CloneOptions{OverwriteExisting: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Skipped {
		t.Error("expected result to be skipped")
	}
	if result.SkipReason == "" {
		t.Error("expected a skip reason")
	}
}

func TestClone_WritesWhenOverwriteEnabled(t *testing.T) {
	written := false
	client := &cloneTestClient{
		readFunc: func(path string) (*vaultapi.Secret, error) {
			return makeVersionsSecret(2), nil
		},
		writeFunc: func(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
			written = true
			return nil, nil
		},
		listFunc: func(path string) (*vaultapi.Secret, error) { return nil, nil },
	}

	_, err := Clone(client, "secret/src", "secret/dest", "latest", CloneOptions{OverwriteExisting: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !written {
		t.Error("expected write to be called")
	}
}

func TestClone_SameSourceAndDest(t *testing.T) {
	_, err := Clone(nil, "secret/same", "secret/same", "latest", CloneOptions{})
	if err == nil {
		t.Fatal("expected error when sourcePath and destPath are identical")
	}
}
