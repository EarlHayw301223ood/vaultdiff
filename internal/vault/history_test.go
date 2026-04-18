package vault

import (
	"testing"
	"time"
)

func makeHistory(versions []VersionMeta) *VersionHistory {
	return &VersionHistory{
		Path:     "secret/data/app/config",
		Versions: versions,
	}
}

func TestVersionHistory_Latest_Empty(t *testing.T) {
	h := makeHistory(nil)
	_, ok := h.Latest()
	if ok {
		t.Fatal("expected ok=false for empty history")
	}
}

func TestVersionHistory_Latest_ReturnsHighest(t *testing.T) {
	h := makeHistory([]VersionMeta{
		{Version: 1},
		{Version: 2},
		{Version: 3},
	})
	v, ok := h.Latest()
	if !ok {
		t.Fatal("expected ok=true")
	}
	if v.Version != 3 {
		t.Fatalf("expected version 3, got %d", v.Version)
	}
}

func TestVersionHistory_ActiveVersions_FiltersDeleted(t *testing.T) {
	now := time.Now()
	h := makeHistory([]VersionMeta{
		{Version: 1},
		{Version: 2, DeletedTime: &now},
		{Version: 3, Destroyed: true},
		{Version: 4},
	})
	active := h.ActiveVersions()
	if len(active) != 2 {
		t.Fatalf("expected 2 active versions, got %d", len(active))
	}
	if active[0].Version != 1 || active[1].Version != 4 {
		t.Fatalf("unexpected active versions: %+v", active)
	}
}

func TestVersionHistory_ActiveVersions_AllActive(t *testing.T) {
	h := makeHistory([]VersionMeta{
		{Version: 1},
		{Version: 2},
	})
	active := h.ActiveVersions()
	if len(active) != 2 {
		t.Fatalf("expected 2 active versions, got %d", len(active))
	}
}

func TestVersionHistory_ActiveVersions_NoneActive(t *testing.T) {
	now := time.Now()
	h := makeHistory([]VersionMeta{
		{Version: 1, DeletedTime: &now},
		{Version: 2, Destroyed: true},
	})
	active := h.ActiveVersions()
	if len(active) != 0 {
		t.Fatalf("expected 0 active versions, got %d", len(active))
	}
}

func TestGetHistory_EmptyPath(t *testing.T) {
	_, err := GetHistory(nil, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}
