package vault

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/api"
)

const signMetaPrefix = "vaultdiff/signatures"

func signMetaPath(path string) string {
	return fmt.Sprintf("%s/%s", signMetaPrefix, strings.Trim(path, "/"))
}

// SignResult holds the outcome of a sign operation.
type SignResult struct {
	Path      string
	Version   int
	Signature string
}

// SignSecret computes an HMAC-SHA256 over the secret's key=value pairs and stores it.
func SignSecret(client *api.Client, path string, version int, hmacKey string) (*SignResult, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	if version < 1 {
		return nil, errors.New("version must be >= 1")
	}
	if hmacKey == "" {
		return nil, errors.New("hmac key must not be empty")
	}

	secret, err := FetchAtRef(client, path, fmt.Sprintf("%d", version))
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	if secret == nil {
		return nil, errors.New("secret not found")
	}

	data := extractStringMap(secret)
	sig := computeSignature(data, hmacKey)

	_, err = client.Logical().Write(signMetaPath(path), map[string]interface{}{
		"version":   version,
		"signature": sig,
	})
	if err != nil {
		return nil, fmt.Errorf("write signature: %w", err)
	}

	return &SignResult{Path: path, Version: version, Signature: sig}, nil
}

// VerifySecret recomputes the signature and compares it to the stored one.
func VerifySecret(client *api.Client, path string, version int, hmacKey string) (bool, error) {
	if path == "" {
		return false, errors.New("path must not be empty")
	}
	if version < 1 {
		return false, errors.New("version must be >= 1")
	}

	stored, err := client.Logical().Read(signMetaPath(path))
	if err != nil {
		return false, fmt.Errorf("read signature: %w", err)
	}
	if stored == nil || stored.Data == nil {
		return false, errors.New("no signature found")
	}

	storedSig, _ := stored.Data["signature"].(string)

	secret, err := FetchAtRef(client, path, fmt.Sprintf("%d", version))
	if err != nil {
		return false, fmt.Errorf("fetch: %w", err)
	}
	if secret == nil {
		return false, errors.New("secret not found")
	}

	data := extractStringMap(secret)
	sig := computeSignature(data, hmacKey)
	return hmac.Equal([]byte(sig), []byte(storedSig)), nil
}

func computeSignature(data map[string]string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		h.Write([]byte(k + "=" + data[k] + ";"))
	}
	return hex.EncodeToString(h.Sum(nil))
}
