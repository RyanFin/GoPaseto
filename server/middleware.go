package server

import (
	"RyanFin/GoPaseto/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey        = "Authorization"
	authorizationHeaderBearerType = "bearer"
)

// Authentication Middleware using a high-order function
func authMiddleware(maker token.PasetoMaker) gin.HandlerFunc {
	// actual authentication middleware code in this func
	return func(ctx *gin.Context) {
		// Header check for 'Authorization'.
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if authHeader == "" {
			// H is a shortcut for map[string]any
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No header was passed"})
			return
		}

		// Token format validation: Ensure that a valid bearer token exists.
		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or Missing Bearer Token"})
			return
		}

		// Authorization type check: Check that the token is of the type 'bearer'
		authType := fields[0]
		if strings.ToLower(authType) != authorizationHeaderBearerType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization Type Not Supported"})
			return
		}

		// Token verification: attempt to verify the authenticity of the provided token
		token := fields[1]
		_, err := maker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access Token Not Valid"})
			return
		}
		// Move onto next peice of middleware or actual endpoint function
		ctx.Next()
	}
}
