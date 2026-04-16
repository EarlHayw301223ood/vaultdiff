package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestArchive_EmptyPath(t *testing.T) {
	client := &stubLogical{}
	_, err := Archive(client, "", 1)
	if err == nil || err.Error() != "archive: path must not be empty" {
		t.Fatalf("expected empty path error, got %v", err)
	}
}

func TestArchive_InvalidVersion(t *testing.T) {
	client := &stubLogical{}
	_, err := Archive(client, "secret/foo", 0)
	if err == nil || err.Error() != "archive: version must be >= 1" {
		t.Fatalf("expected version error, got %v", err)
	}
}

func TestArchive_FetchError(t *testing.T) {
	client := &stubLogical{
		readFn: func(path string) (*vaultapi.Secret, error) {
			return nil, errStub
		},
	}
	_, err := Archive(client, "secret/foo", 1)
	if err == nil {
		t.Fatal("expected fetch error")
	}
}

func TestArchive_NilSecret(t *testing.T) {
	client := &stubLogical{
		readFn: func(path string) (*vaultapi.Secret, error) {
			return nil, nil
		},
	}
	_, err := Archive(client, "secret/foo", 1)
	if err == nil {
		t.Fatal("expected nil secret error")
	}
}

func TestArchiveDestPath_Format(t *testing.T) {
	got := archiveDestPath("myapp/db", 3)
	want := "archive/myapp/db/v3"
	if got != want {
		t.Fatalf("want %s got %s", want, got)
	}
}

func TestArchiveResult_Fields(t *testing.T) {
	client := &stubLogical{
		readFn: func(path string) (*vaultapi.Secret, error) {
			return &vaultapi.Secret{
				Data: map[string]interface{}{
					"data": map[string]interface{}{"key": "val"},
				},
			}, nil
		},
		writeFn: func(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
			return &vaultapi.Secret{}, nil
		},
	}
	res, err := Archive(client, "secret/foo", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Version != 2 {
		t.Errorf("expected version 2, got %d", res.Version)
	}
	if res.Path != "archive/secret/foo/v2" {
		t.Errorf("unexpected path: %s", res.Path)
	}
	if res.ArchivedAt.IsZero() {
		t.Error("expected ArchivedAt to be set")
	}
}
