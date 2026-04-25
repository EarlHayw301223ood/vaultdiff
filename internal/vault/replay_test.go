package vault

import (
	"testing"
	"time"
)

func makeReplayLog(versions []int) *ReplayLog {
	log := &ReplayLog{Path: "secret/app/config"}
	for _, v := range versions {
		log.Entries = append(log.Entries, ReplayEntry{
			Version:   v,
			Data:      map[string]string{"key": fmt.Sprintf("val%d", v)},
			CreatedAt: time.Now().UTC(),
		})
	}
	return log
}

func TestReplayLog_Latest_Empty(t *testing.T) {
	log := &ReplayLog{Path: "secret/empty"}
	_, ok := log.Latest()
	if ok {
		t.Fatal("expected Latest to return false for empty log")
	}
}

func TestReplayLog_Latest_ReturnsLast(t *testing.T) {
	log := makeReplayLog([]int{1, 2, 3})
	entry, ok := log.Latest()
	if !ok {
		t.Fatal("expected Latest to return true")
	}
	if entry.Version != 3 {
		t.Fatalf("expected version 3, got %d", entry.Version)
	}
}

func TestReplayLog_At_Found(t *testing.T) {
	log := makeReplayLog([]int{1, 2, 3})
	entry, err := log.At(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Version != 2 {
		t.Fatalf("expected version 2, got %d", entry.Version)
	}
}

func TestReplayLog_At_NotFound(t *testing.T) {
	log := makeReplayLog([]int{1, 3})
	_, err := log.At(2)
	if err == nil {
		t.Fatal("expected error for missing version")
	}
}

func TestReplaySecret_EmptyPath(t *testing.T) {
	client := newStubClient()
	_, err := ReplaySecret(client, "", "secret")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestReplaySecret_EmptyMount(t *testing.T) {
	client := newStubClient()
	_, err := ReplaySecret(client, "app/config", "")
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}
