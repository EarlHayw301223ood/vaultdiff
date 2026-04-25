package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend_SuccessfulDelivery(t *testing.T) {
	var received WebhookEvent
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	d := NewDispatcher()
	event := WebhookEvent{
		Path:      "secret/app/db",
		Operation: "write",
		Version:   3,
		Timestamp: time.Now().UTC(),
	}
	cfgs := []WebhookConfig{{URL: ts.URL, Timeout: 5 * time.Second}}
	results := d.Send(context.Background(), event, cfgs)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", results[0].StatusCode)
	}
	if received.Path != event.Path {
		t.Errorf("payload path mismatch: got %q", received.Path)
	}
}

func TestSend_Non2xxSetsErr(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	d := NewDispatcher()
	results := d.Send(context.Background(), WebhookEvent{}, []WebhookConfig{{URL: ts.URL}})
	if results[0].Err == nil {
		t.Error("expected error for 500 response")
	}
	if results[0].StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", results[0].StatusCode)
	}
}

func TestSend_InvalidURL(t *testing.T) {
	d := NewDispatcher()
	results := d.Send(context.Background(), WebhookEvent{}, []WebhookConfig{{URL: "not-a-url"}})
	if results[0].Err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestSend_CustomHeaders(t *testing.T) {
	var authHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	d := NewDispatcher()
	cfg := WebhookConfig{
		URL:     ts.URL,
		Headers: map[string]string{"Authorization": "Bearer token123"},
	}
	d.Send(context.Background(), WebhookEvent{}, []WebhookConfig{cfg})
	if authHeader != "Bearer token123" {
		t.Errorf("expected auth header, got %q", authHeader)
	}
}

func TestSend_MultipleTargets(t *testing.T) {
	makeServer := func(code int) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
		}))
	}
	s1 := makeServer(200)
	s2 := makeServer(201)
	defer s1.Close()
	defer s2.Close()

	d := NewDispatcher()
	cfgs := []WebhookConfig{{URL: s1.URL}, {URL: s2.URL}}
	results := d.Send(context.Background(), WebhookEvent{}, cfgs)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].StatusCode != 200 || results[1].StatusCode != 201 {
		t.Errorf("unexpected status codes: %d, %d", results[0].StatusCode, results[1].StatusCode)
	}
}
