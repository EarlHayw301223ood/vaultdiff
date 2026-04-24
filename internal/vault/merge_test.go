package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestMerge_EmptySourcePath(t *testing.T) {
	c := &stubLogical{}
	_, err := Merge(c, "", "secret/data/dst", MergeOptions{Strategy: MergeStrategyOurs})
	if err == nil || err.Error() != "source path must not be empty" {
		t.Fatalf("expected source empty error, got %v", err)
	}
}

func TestMerge_EmptyDestPath(t *testing.T) {
	c := &stubLogical{}
	_, err := Merge(c, "secret/data/src", "", MergeOptions{Strategy: MergeStrategyOurs})
	if err == nil || err.Error() != "destination path must not be empty" {
		t.Fatalf("expected dest empty error, got %v", err)
	}
}

func TestMerge_SameSourceAndDest(t *testing.T) {
	c := &stubLogical{}
	_, err := Merge(c, "secret/data/x", "secret/data/x", MergeOptions{})
	if err == nil {
		t.Fatal("expected error for same paths")
	}
}

func TestMerge_OursStrategy_KeepsDestOnConflict(t *testing.T) {
	srcData := map[string]interface{}{"key": "src-value", "new": "added"}
	dstData := map[string]interface{}{"key": "dst-value"}

	c := &stubLogical{
		readResponses: map[string]*vaultapi.Secret{
			"secret/data/src": makeDataSecret(srcData, 1),
			"secret/data/dst": makeDataSecret(dstData, 1),
		},
		writeResponses: map[string]*vaultapi.Secret{
			"secret/data/data/dst": {},
		},
	}

	res, err := Merge(c, "secret/data/src", "secret/data/dst", MergeOptions{Strategy: MergeStrategyOurs})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Data["key"] != "dst-value" {
		t.Errorf("expected dst-value, got %s", res.Data["key"])
	}
	if res.Data["new"] != "added" {
		t.Errorf("expected new key to be added, got %s", res.Data["new"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "key" {
		t.Errorf("expected one conflict on 'key', got %v", res.Conflicts)
	}
}

func TestMerge_TheirsStrategy_OverwritesOnConflict(t *testing.T) {
	srcData := map[string]interface{}{"key": "src-value"}
	dstData := map[string]interface{}{"key": "dst-value"}

	c := &stubLogical{
		readResponses: map[string]*vaultapi.Secret{
			"secret/data/src": makeDataSecret(srcData, 1),
			"secret/data/dst": makeDataSecret(dstData, 1),
		},
		writeResponses: map[string]*vaultapi.Secret{
			"secret/data/data/dst": {},
		},
	}

	res, err := Merge(c, "secret/data/src", "secret/data/dst", MergeOptions{Strategy: MergeStrategyTheirs})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Data["key"] != "src-value" {
		t.Errorf("expected src-value, got %s", res.Data["key"])
	}
}

func TestMerge_DryRun_DoesNotWrite(t *testing.T) {
	srcData := map[string]interface{}{"only": "here"}
	dstData := map[string]interface{}{}

	c := &stubLogical{
		readResponses: map[string]*vaultapi.Secret{
			"secret/data/src": makeDataSecret(srcData, 1),
			"secret/data/dst": makeDataSecret(dstData, 1),
		},
	}

	res, err := Merge(c, "secret/data/src", "secret/data/dst", MergeOptions{Strategy: MergeStrategyOurs, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "only" {
		t.Errorf("expected 'only' in Added, got %v", res.Added)
	}
	if len(c.writeLog) != 0 {
		t.Errorf("expected no writes in dry-run, got %d", len(c.writeLog))
	}
}
