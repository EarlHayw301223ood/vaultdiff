package vault

import (
	"testing"
)

func TestCheckCompliance_PassesCleanSecret(t *testing.T) {
	data := map[string]string{
		"api_key": "abc123",
		"region":  "us-east-1",
	}
	result := CheckCompliance("secret/app", 1, data, DefaultComplianceRules())
	if result.Status != CompliancePass {
		t.Errorf("expected pass, got %s", result.Status)
	}
	if len(result.Violations) != 0 {
		t.Errorf("expected no violations, got %d", len(result.Violations))
	}
}

func TestCheckCompliance_FailsEmptyValue(t *testing.T) {
	data := map[string]string{
		"api_key": "",
	}
	result := CheckCompliance("secret/app", 2, data, DefaultComplianceRules())
	if result.Status != ComplianceFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
	found := false
	for _, v := range result.Violations {
		if v.Rule == "no-empty-values" {
			found = true
		}
	}
	if !found {
		t.Error("expected no-empty-values violation")
	}
}

func TestCheckCompliance_FailsEmptySecret(t *testing.T) {
	data := map[string]string{}
	result := CheckCompliance("secret/empty", 1, data, DefaultComplianceRules())
	if result.Status != ComplianceFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
	found := false
	for _, v := range result.Violations {
		if v.Rule == "min-secret-count" {
			found = true
		}
	}
	if !found {
		t.Error("expected min-secret-count violation")
	}
}

func TestCheckCompliance_FailsPlaintextPassword(t *testing.T) {
	data := map[string]string{
		"password": "password",
	}
	result := CheckCompliance("secret/creds", 1, data, DefaultComplianceRules())
	if result.Status != ComplianceFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
	found := false
	for _, v := range result.Violations {
		if v.Rule == "no-plaintext-password-key" {
			found = true
		}
	}
	if !found {
		t.Error("expected no-plaintext-password-key violation")
	}
}

func TestCheckCompliance_PathAndVersionStored(t *testing.T) {
	data := map[string]string{"k": "v"}
	result := CheckCompliance("secret/myapp", 7, data, DefaultComplianceRules())
	if result.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", result.Path)
	}
	if result.Version != 7 {
		t.Errorf("expected version 7, got %d", result.Version)
	}
}

func TestCheckCompliance_CheckedAtIsUTC(t *testing.T) {
	data := map[string]string{"k": "v"}
	result := CheckCompliance("secret/ts", 1, data, DefaultComplianceRules())
	if result.CheckedAt.Location().String() != "UTC" {
		t.Errorf("expected UTC timestamp, got %s", result.CheckedAt.Location())
	}
}

func TestCheckCompliance_CustomRule(t *testing.T) {
	rule := ComplianceRule{
		Name: "must-have-owner",
		Check: func(data map[string]string) (bool, string) {
			if _, ok := data["owner"]; !ok {
				return false, "missing required key 'owner'"
			}
			return true, ""
		},
	}
	data := map[string]string{"api_key": "abc"}
	result := CheckCompliance("secret/app", 1, data, []ComplianceRule{rule})
	if result.Status != ComplianceFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
	if result.Violations[0].Rule != "must-have-owner" {
		t.Errorf("unexpected rule name: %s", result.Violations[0].Rule)
	}
}
