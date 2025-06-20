package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"Accept",
			"Accept-Language",
			"Content-Language",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
		},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	}
}

// ProductionCORSConfig returns a production-safe CORS configuration
func ProductionCORSConfig(allowedOrigins []string) CORSConfig {
	config := DefaultCORSConfig()
	config.AllowedOrigins = allowedOrigins
	config.AllowCredentials = true
	return config
}

// CORS returns a CORS middleware with default configuration
func CORS() gin.HandlerFunc {
	return CORSWithConfig(DefaultCORSConfig())
}

// CORSWithConfig returns a CORS middleware with custom configuration
func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		method := c.Request.Method

		// Set CORS headers
		setCORSHeaders(c, config, origin)

		// Handle preflight requests
		if method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// setCORSHeaders sets the appropriate CORS headers based on configuration
func setCORSHeaders(c *gin.Context, config CORSConfig, origin string) {
	// Access-Control-Allow-Origin
	if len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*" {
		c.Header("Access-Control-Allow-Origin", "*")
	} else if isOriginAllowed(origin, config.AllowedOrigins) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Vary", "Origin")
	}

	// Access-Control-Allow-Methods
	if len(config.AllowedMethods) > 0 {
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
	}

	// Access-Control-Allow-Headers
	if len(config.AllowedHeaders) > 0 {
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
	}

	// Access-Control-Expose-Headers
	if len(config.ExposedHeaders) > 0 {
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
	}

	// Access-Control-Allow-Credentials
	if config.AllowCredentials {
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	// Access-Control-Max-Age
	if config.MaxAge > 0 {
		c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
	}
}

// isOriginAllowed checks if the origin is in the allowed origins list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
		// Support for wildcard subdomains (e.g., *.example.com)
		if strings.Contains(allowedOrigin, "*") {
			if matchWildcard(allowedOrigin, origin) {
				return true
			}
		}
	}
	return false
}

// matchWildcard checks if origin matches wildcard pattern
func matchWildcard(pattern, origin string) bool {
	if pattern == "*" {
		return true
	}

	// Simple wildcard matching for patterns like *.example.com
	if strings.HasPrefix(pattern, "*.") {
		domain := strings.TrimPrefix(pattern, "*.")
		return strings.HasSuffix(origin, domain)
	}

	return pattern == origin
}

// DevelopmentCORS returns a permissive CORS middleware for development
func DevelopmentCORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}
