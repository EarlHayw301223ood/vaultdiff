package vault

import (
	"testing"
)

func TestQuotaMetaPath_Format(t *testing.T) {
	got := quotaMetaPath("secret/myapp/db")
	want := "vaultdiff/meta/secret/myapp/db/__quota"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestQuotaMetaPath_TrimsSlashes(t *testing.T) {
	got := quotaMetaPath("/secret/myapp/")
	want := "vaultdiff/meta/secret/myapp/__quota"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetQuota_EmptyPath(t *testing.T) {
	c := newStubClient()
	err := SetQuota(c, "", 10, 60)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSetQuota_ZeroMaxWrites(t *testing.T) {
	c := newStubClient()
	err := SetQuota(c, "secret/app", 0, 60)
	if err == nil {
		t.Fatal("expected error for zero max_writes")
	}
}

func TestSetQuota_ZeroWindowSec(t *testing.T) {
	c := newStubClient()
	err := SetQuota(c, "secret/app", 10, 0)
	if err == nil {
		t.Fatal("expected error for zero window_sec")
	}
}

func TestGetQuota_EmptyPath(t *testing.T) {
	c := newStubClient()
	_, err := GetQuota(c, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetQuota_NilWhenNotFound(t *testing.T) {
	c := newStubClient()
	record, err := GetQuota(c, "secret/missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record != nil {
		t.Errorf("expected nil record, got %+v", record)
	}
}

func TestQuotaRecord_Exceeded(t *testing.T) {
	q := QuotaRecord{MaxWrites: 5, Writes: 5}
	if !q.Exceeded() {
		t.Error("expected Exceeded() to be true")
	}
}

func TestQuotaRecord_NotExceeded(t *testing.T) {
	q := QuotaRecord{MaxWrites: 5, Writes: 3}
	if q.Exceeded() {
		t.Error("expected Exceeded() to be false")
	}
}

func TestQuotaRecord_Remaining(t *testing.T) {
	q := QuotaRecord{MaxWrites: 10, Writes: 3}
	if got := q.Remaining(); got != 7 {
		t.Errorf("got %d, want 7", got)
	}
}

func TestQuotaRecord_RemainingClampedAtZero(t *testing.T) {
	q := QuotaRecord{MaxWrites: 5, Writes: 8}
	if got := q.Remaining(); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}
