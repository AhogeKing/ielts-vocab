package middleware

import (
	"net/http"
	"strings"

	"Ielts-vocab/internal/auth"
	"Ielts-vocab/internal/common"

	"github.com/gin-gonic/gin"
)

// RequireJWT protects every route in its group except the explicitly listed
// public paths. A valid token's claims are available under auth.ClaimsContextKey.
func RequireJWT(tokens *auth.TokenManager, publicPaths ...string) gin.HandlerFunc {
	public := make(map[string]struct{}, len(publicPaths))
	for _, path := range publicPaths {
		public[path] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := public[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		header := c.GetHeader("Authorization")
		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			common.Fail(c, http.StatusUnauthorized, "AUTHORIZATION_REQUIRED", "请提供 Bearer Token")
			c.Abort()
			return
		}

		claims, err := tokens.Parse(parts[1])
		if err != nil {
			common.Fail(c, http.StatusUnauthorized, "INVALID_TOKEN", "Token 无效或已过期")
			c.Abort()
			return
		}

		c.Set(auth.ClaimsContextKey, claims)
		c.Next()
	}
}
