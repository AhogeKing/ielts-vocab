package auth

import (
	"testing"
	"time"
)

func TestTokenManagerGeneratesAndParsesToken(t *testing.T) {
	manager, err := NewTokenManager("01234567890123456789012345678901", "ielts-vocab-test", time.Hour)
	if err != nil {
		t.Fatalf("NewTokenManager returned an error: %v", err)
	}

	accessToken, err := manager.Generate(42, "xavier")
	if err != nil {
		t.Fatalf("Generate returned an error: %v", err)
	}

	claims, err := manager.Parse(accessToken.AccessToken)
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}
	if claims.Subject != "42" || claims.Username != "xavier" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}
