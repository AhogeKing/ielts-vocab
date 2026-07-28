package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Ielts-vocab/internal/auth"

	"github.com/gin-gonic/gin"
)

func TestRequireJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tokens, err := auth.NewTokenManager("01234567890123456789012345678901", "ielts-vocab-test", time.Hour)
	if err != nil {
		t.Fatalf("NewTokenManager returned an error: %v", err)
	}

	router := gin.New()
	router.GET("/private", RequireJWT(tokens), func(c *gin.Context) {
		claims, ok := c.Get(auth.ClaimsContextKey)
		if !ok || claims.(auth.Claims).Subject != "42" {
			t.Fatal("authenticated request is missing expected claims")
		}
		c.Status(http.StatusNoContent)
	})

	unauthenticated := httptest.NewRecorder()
	router.ServeHTTP(unauthenticated, httptest.NewRequest(http.MethodGet, "/private", nil))
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d for a missing token, got %d", http.StatusUnauthorized, unauthenticated.Code)
	}

	accessToken, err := tokens.Generate(42, "xavier")
	if err != nil {
		t.Fatalf("Generate returned an error: %v", err)
	}

	authenticated := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/private", nil)
	request.Header.Set("Authorization", "Bearer "+accessToken.AccessToken)
	router.ServeHTTP(authenticated, request)
	if authenticated.Code != http.StatusNoContent {
		t.Fatalf("expected %d for a valid token, got %d", http.StatusNoContent, authenticated.Code)
	}
}
