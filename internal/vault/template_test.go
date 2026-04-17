package vault

import (
	"testing"
)

func TestInterpolate_AllKeysPresent(t *testing.T) {
	tmpl := "host={{ HOST }} port={{ PORT }}"
	data := map[string]string{"HOST": "localhost", "PORT": "5432"}

	got, err := interpolate(tmpl, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "host=localhost port=5432"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInterpolate_MissingKey(t *testing.T) {
	tmpl := "user={{ USER }} pass={{ PASS }}"
	data := map[string]string{"USER": "admin"}

	_, err := interpolate(tmpl, data)
	if err == nil {
		t.Fatal("expected error for missing placeholder, got nil")
	}
}

func TestInterpolate_NoPlaceholders(t *testing.T) {
	tmpl := "static string"
	data := map[string]string{"KEY": "val"}

	got, err := interpolate(tmpl, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != tmpl {
		t.Errorf("got %q, want %q", got, tmpl)
	}
}

func TestInterpolate_EmptyData(t *testing.T) {
	tmpl := "value={{ SECRET }}"
	data := map[string]string{}

	_, err := interpolate(tmpl, data)
	if err == nil {
		t.Fatal("expected error for empty data with placeholder")
	}
}

func TestRenderTemplate_EmptyPath(t *testing.T) {
	_, err := RenderTemplate(nil, "", "latest", "{{ KEY }}")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRenderTemplate_EmptyTemplate(t *testing.T) {
	_, err := RenderTemplate(nil, "secret/myapp", "latest", "")
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestInterpolate_MultipleOccurrences(t *testing.T) {
	tmpl := "{{ HOST }}:{{ PORT }}/{{ HOST }}"
	data := map[string]string{"HOST": "db", "PORT": "3306"}

	got, err := interpolate(tmpl, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "db:3306/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
