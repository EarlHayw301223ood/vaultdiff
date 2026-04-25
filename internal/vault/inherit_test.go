package vault

import (
	"testing"
)

func TestInheritSecret_EmptyParentPath(t *testing.T) {
	_, err := InheritSecret(nil, "", "child/path", "latest", "latest")
	if err == nil || err.Error() != "parent path must not be empty" {
		t.Fatalf("expected empty parent path error, got %v", err)
	}
}

func TestInheritSecret_EmptyChildPath(t *testing.T) {
	_, err := InheritSecret(nil, "parent/path", "", "latest", "latest")
	if err == nil || err.Error() != "child path must not be empty" {
		t.Fatalf("expected empty child path error, got %v", err)
	}
}

func TestInheritSecret_SamePaths(t *testing.T) {
	_, err := InheritSecret(nil, "secret/app", "secret/app", "latest", "latest")
	if err == nil {
		t.Fatal("expected error for identical paths")
	}
}

func TestInheritSecret_ChildOverridesParent(t *testing.T) {
	parent := makeKVSecret(map[string]interface{}{"HOST": "prod.example.com", "PORT": "5432"})
	child := makeKVSecret(map[string]interface{}{"HOST": "staging.example.com", "DEBUG": "true"})

	calls := 0
	stub := &stubVersionClient{
		readFn: func(path string) (*stubSecret, error) {
			calls++
			if calls == 1 {
				return parent, nil
			}
			return child, nil
		},
	}

	res, err := InheritSecret(stub, "secret/base", "secret/staging", "latest", "latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Data["HOST"] != "staging.example.com" {
		t.Errorf("expected child HOST, got %q", res.Data["HOST"])
	}
	if res.Data["PORT"] != "5432" {
		t.Errorf("expected inherited PORT, got %q", res.Data["PORT"])
	}
	if res.Data["DEBUG"] != "true" {
		t.Errorf("expected child-only DEBUG, got %q", res.Data["DEBUG"])
	}
	if res.Override != 1 {
		t.Errorf("expected 1 override, got %d", res.Override)
	}
	if res.Inherited != 1 {
		t.Errorf("expected 1 inherited, got %d", res.Inherited)
	}
}

func TestInheritSecret_NoConflicts(t *testing.T) {
	parent := makeKVSecret(map[string]interface{}{"A": "1", "B": "2"})
	child := makeKVSecret(map[string]interface{}{"C": "3"})

	calls := 0
	stub := &stubVersionClient{
		readFn: func(_ string) (*stubSecret, error) {
			calls++
			if calls == 1 {
				return parent, nil
			}
			return child, nil
		},
	}

	res, err := InheritSecret(stub, "secret/base", "secret/child", "latest", "latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Data) != 3 {
		t.Errorf("expected 3 keys, got %d", len(res.Data))
	}
	if res.Override != 0 {
		t.Errorf("expected 0 overrides, got %d", res.Override)
	}
	if res.Inherited != 2 {
		t.Errorf("expected 2 inherited, got %d", res.Inherited)
	}
}
