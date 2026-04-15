package vault

import (
	"testing"
)

func TestKvV2WritePath_Standard(t *testing.T) {
	got := kvV2WritePath("secret/myapp/config")
	want := "secret/data/myapp/config"
	if got != want {
		t.Errorf("kvV2WritePath: got %q, want %q", got, want)
	}
}

func TestKvV2WritePath_SingleSegment(t *testing.T) {
	got := kvV2WritePath("secret")
	// no slash — mount prefix is the whole string, suffix is empty
	want := "secret/data/"
	if got != want {
		t.Errorf("kvV2WritePath single segment: got %q, want %q", got, want)
	}
}

func TestKvV2WritePath_EmptyPath(t *testing.T) {
	got := kvV2WritePath("")
	if got != "" {
		t.Errorf("kvV2WritePath empty: expected empty string, got %q", got)
	}
}

func TestToInterfaceMap_PreservesValues(t *testing.T) {
	input := map[string]string{"key1": "val1", "key2": "val2"}
	out := toInterfaceMap(input)
	if len(out) != len(input) {
		t.Fatalf("toInterfaceMap: length mismatch: got %d, want %d", len(out), len(input))
	}
	for k, v := range input {
		got, ok := out[k]
		if !ok {
			t.Errorf("toInterfaceMap: missing key %q", k)
			continue
		}
		if got != v {
			t.Errorf("toInterfaceMap[%q]: got %v, want %v", k, got, v)
		}
	}
}

func TestToInterfaceMap_Empty(t *testing.T) {
	out := toInterfaceMap(map[string]string{})
	if len(out) != 0 {
		t.Errorf("toInterfaceMap empty: expected empty map, got %v", out)
	}
}

func TestPromoteResult_Fields(t *testing.T) {
	r := &PromoteResult{
		SourcePath:    "secret/staging/app",
		DestPath:      "secret/prod/app",
		SourceVersion: 3,
		DestVersion:   1,
		Keys:          5,
	}
	if r.SourcePath != "secret/staging/app" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if r.Keys != 5 {
		t.Errorf("unexpected Keys count: %d", r.Keys)
	}
}
