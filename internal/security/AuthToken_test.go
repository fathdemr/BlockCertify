package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func generateTestKeyPair(t *testing.T) (privatePEM, publicPEM string) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}

	privateDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}
	privatePEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateDER}))

	publicDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	publicPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER}))

	return privatePEM, publicPEM
}

func TestRSASignVerifyRoundtrip(t *testing.T) {
	privatePEM, publicPEM := generateTestKeyPair(t)

	helper, err := NewRSAJWTHelper(privatePEM, publicPEM, time.Hour)
	if err != nil {
		t.Fatalf("NewRSAJWTHelper failed: %v", err)
	}

	token, err := helper.Sign(jwt.MapClaims{"email": "test@example.com", "role": "admin"})
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	claims, err := helper.Verify(token)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if claims["email"] != "test@example.com" {
		t.Errorf("email claim = %v, want test@example.com", claims["email"])
	}
}

func TestVerifyRejectsExpiredToken(t *testing.T) {
	privatePEM, publicPEM := generateTestKeyPair(t)

	helper, err := NewRSAJWTHelper(privatePEM, publicPEM, time.Hour)
	if err != nil {
		t.Fatalf("NewRSAJWTHelper failed: %v", err)
	}

	token, err := helper.Sign(jwt.MapClaims{
		"email": "test@example.com",
		"exp":   time.Now().Add(-time.Hour).Unix(),
	})
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if _, err := helper.Verify(token); err == nil {
		t.Error("Verify accepted an expired token")
	}
}

// TestVerifyRejectsHMACToken guards against the classic RS256→HS256 algorithm
// confusion attack: a token signed with HMAC using the public key as secret
// must not validate.
func TestVerifyRejectsHMACToken(t *testing.T) {
	privatePEM, publicPEM := generateTestKeyPair(t)

	helper, err := NewRSAJWTHelper(privatePEM, publicPEM, time.Hour)
	if err != nil {
		t.Fatalf("NewRSAJWTHelper failed: %v", err)
	}

	forged := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "attacker@example.com",
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	forgedStr, err := forged.SignedString([]byte(publicPEM))
	if err != nil {
		t.Fatalf("failed to sign forged token: %v", err)
	}

	if _, err := helper.Verify(forgedStr); err == nil {
		t.Error("Verify accepted an HS256-signed token")
	}
}

func TestVerifyRejectsTamperedToken(t *testing.T) {
	privatePEM, publicPEM := generateTestKeyPair(t)
	_, otherPublicPEM := generateTestKeyPair(t)

	signer, err := NewRSAJWTHelper(privatePEM, publicPEM, time.Hour)
	if err != nil {
		t.Fatalf("NewRSAJWTHelper failed: %v", err)
	}
	verifier, err := NewRSAJWTHelper("", otherPublicPEM, time.Hour)
	if err != nil {
		t.Fatalf("NewRSAJWTHelper (verify-only) failed: %v", err)
	}

	token, err := signer.Sign(jwt.MapClaims{"email": "test@example.com"})
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if _, err := verifier.Verify(token); err == nil {
		t.Error("Verify accepted a token signed by a different key")
	}
}

func TestSignWithoutPrivateKeyFails(t *testing.T) {
	_, publicPEM := generateTestKeyPair(t)

	helper, err := NewRSAJWTHelper("", publicPEM, time.Hour)
	if err != nil {
		t.Fatalf("NewRSAJWTHelper failed: %v", err)
	}

	if _, err := helper.Sign(jwt.MapClaims{"email": "x@example.com"}); err == nil {
		t.Error("Sign succeeded without a private key")
	}
}
