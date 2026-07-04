package tests

import (
	"BlockCertify/internal/config"
	"os"
	"testing"
)

// TestConfigRSAKeysParse verifies that the RSA key pair in internal/config.yaml
// is valid PEM and can back the JWT helper. Skipped when the (gitignored)
// config file is not present, e.g. in CI.
func TestConfigRSAKeysParse(t *testing.T) {
	if _, err := os.Stat("../config.yaml"); os.IsNotExist(err) {
		t.Skip("internal/config.yaml not present")
	}

	if err := config.InitConfigFile(".."); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if err := config.InitJWT(); err != nil {
		t.Fatalf("failed to init JWT helper from config RSA keys: %v", err)
	}
	if config.JWT == nil {
		t.Fatal("config.JWT is nil after InitJWT")
	}
}
