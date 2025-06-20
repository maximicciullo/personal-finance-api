package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger initializes the global zap logger
func InitLogger(environment string) error {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else if environment == "test" {
		// Silent logger for tests
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.FatalLevel) // Only fatal logs
		config.OutputPaths = []string{"/dev/null"}
	} else {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Customize encoding
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = "message"

	var err error
	Logger, err = config.Build(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		return err
	}

	// Replace global logger
	zap.ReplaceGlobals(Logger)
	
	return nil
}

// LogConfig holds logging configuration
type LogConfig struct {
	ShowRequest  bool
	ShowResponse bool
	ShowHeaders  bool
	SkipPaths    []string
	MaxBodySize  int64
}

// DefaultLogConfig returns a default logging configuration
func DefaultLogConfig() LogConfig {
	return LogConfig{
		ShowRequest:  true,
		ShowResponse: true,
		ShowHeaders:  false,
		SkipPaths:    []string{"/health"},
		MaxBodySize:  2048, // 2KB
	}
}

// ProductionLogConfig returns a production-safe logging configuration
func ProductionLogConfig() LogConfig {
	return LogConfig{
		ShowRequest:  false,
		ShowResponse: false,
		ShowHeaders:  false,
		SkipPaths:    []string{"/health", "/metrics"},
		MaxBodySize:  1024, // 1KB
	}
}

// DebugLogConfig returns a debug logging configuration
func DebugLogConfig() LogConfig {
	return LogConfig{
		ShowRequest:  true,
		ShowResponse: true,
		ShowHeaders:  true,
		SkipPaths:    []string{},
		MaxBodySize:  4096, // 4KB
	}
}

// ZapLogger returns a gin middleware that uses zap for logging
func ZapLogger() gin.HandlerFunc {
	return ZapLoggerWithConfig(DefaultLogConfig())
}

// ZapLoggerWithConfig returns a gin middleware that uses zap with custom config
func ZapLoggerWithConfig(config LogConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip paths if configured
		for _, path := range config.SkipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read request body if needed
		var requestBody []byte
		if config.ShowRequest && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Log request
		logRequest(c, requestBody, config)

		// Create response body writer
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get final path
		if raw != "" {
			path = path + "?" + raw
		}

		// Log response
		logResponse(c, blw.body.Bytes(), latency, path, config)

		// Log errors if any
		if len(c.Errors) > 0 {
			logErrors(c)
		}
	}
}

// bodyLogWriter captures response body for logging
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// logRequest logs incoming request details using zap
func logRequest(c *gin.Context, body []byte, config LogConfig) {
	fields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
	}

	// Add query parameters
	if len(c.Request.URL.Query()) > 0 {
		fields = append(fields, zap.String("query", c.Request.URL.RawQuery))
	}

	// Add headers if enabled
	if config.ShowHeaders {
		headers := make(map[string]string)
		for name, values := range c.Request.Header {
			if !isSecretHeader(name) {
				headers[name] = strings.Join(values, ", ")
			}
		}
		fields = append(fields, zap.Any("headers", headers))
	}

	// Add request body if enabled and exists
	if config.ShowRequest && len(body) > 0 && int64(len(body)) <= config.MaxBodySize {
		if isJSONContent(c.Request.Header.Get("Content-Type")) {
			var jsonBody interface{}
			if err := json.Unmarshal(body, &jsonBody); err == nil {
				fields = append(fields, zap.Any("request_body", jsonBody))
			} else {
				fields = append(fields, zap.String("request_body_raw", string(body)))
			}
		} else {
			fields = append(fields, zap.String("request_body", string(body)))
		}
	}

	Logger.Info("📥 HTTP Request", fields...)
}

// logResponse logs response details using zap
func logResponse(c *gin.Context, body []byte, latency time.Duration, path string, config LogConfig) {
	statusCode := c.Writer.Status()
	
	fields := []zap.Field{
		zap.Int("status_code", statusCode),
		zap.String("path", path),
		zap.Duration("latency", latency),
		zap.Int("response_size", len(body)),
		zap.String("status_icon", getStatusIcon(statusCode)),
	}

	// Add response headers if enabled
	if config.ShowHeaders {
		headers := make(map[string]string)
		for name, values := range c.Writer.Header() {
			headers[name] = strings.Join(values, ", ")
		}
		fields = append(fields, zap.Any("response_headers", headers))
	}

	// Add response body if enabled and not too large
	if config.ShowResponse && len(body) > 0 && int64(len(body)) <= config.MaxBodySize {
		if isJSONContent(c.Writer.Header().Get("Content-Type")) {
			var jsonBody interface{}
			if err := json.Unmarshal(body, &jsonBody); err == nil {
				fields = append(fields, zap.Any("response_body", jsonBody))
			} else {
				fields = append(fields, zap.String("response_body_raw", string(body)))
			}
		} else {
			fields = append(fields, zap.String("response_body", string(body)))
		}
	}

	// Choose log level based on status code
	switch {
	case statusCode >= 500:
		Logger.Error("📤 HTTP Response", fields...)
	case statusCode >= 400:
		Logger.Warn("📤 HTTP Response", fields...)
	default:
		Logger.Info("📤 HTTP Response", fields...)
	}
}

// logErrors logs any errors that occurred during request processing
func logErrors(c *gin.Context) {
	for _, err := range c.Errors {
		Logger.Error("❌ Request Error",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("error_type", getErrorTypeString(err.Type)),
			zap.Error(err.Err),
		)
	}
}

// getErrorTypeString converts gin.ErrorType to string
func getErrorTypeString(errorType gin.ErrorType) string {
	switch errorType {
	case gin.ErrorTypeBind:
		return "BIND"
	case gin.ErrorTypeRender:
		return "RENDER"
	case gin.ErrorTypePublic:
		return "PUBLIC"
	case gin.ErrorTypePrivate:
		return "PRIVATE"
	case gin.ErrorTypeAny:
		return "ANY"
	default:
		return "UNKNOWN"
	}
}

// getStatusIcon returns an appropriate icon for HTTP status codes
func getStatusIcon(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "✅" // Success
	case statusCode >= 300 && statusCode < 400:
		return "🔄" // Redirect
	case statusCode >= 400 && statusCode < 500:
		return "⚠️" // Client error
	case statusCode >= 500:
		return "🚨" // Server error
	default:
		return "❓" // Unknown
	}
}

// isJSONContent checks if content type is JSON
func isJSONContent(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "application/json")
}

// isSecretHeader checks if header contains sensitive information
func isSecretHeader(headerName string) bool {
	lowerName := strings.ToLower(headerName)
	secretHeaders := []string{
		"authorization",
		"cookie",
		"x-api-key",
		"x-auth-token",
		"password",
	}
	
	for _, secret := range secretHeaders {
		if strings.Contains(lowerName, secret) {
			return true
		}
	}
	return false
}

// DevelopmentLogger returns a verbose logger for development
func DevelopmentLogger() gin.HandlerFunc {
	return ZapLoggerWithConfig(DebugLogConfig())
}

// ProductionLogger returns a minimal logger for production
func ProductionLogger() gin.HandlerFunc {
	return ZapLoggerWithConfig(ProductionLogConfig())
}

// BusinessLogger logs business logic events in different layers
func BusinessLogger() *BusinessLoggerInstance {
	return &BusinessLoggerInstance{logger: Logger}
}

type BusinessLoggerInstance struct {
	logger *zap.Logger
}

// Controller logs controller layer events
func (bl *BusinessLoggerInstance) Controller(operation string, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.String("layer", "controller")}, fields...)
	bl.logger.Info("🎛️ "+operation, allFields...)
}

// Service logs service layer events
func (bl *BusinessLoggerInstance) Service(operation string, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.String("layer", "service")}, fields...)
	bl.logger.Info("⚙️ "+operation, allFields...)
}

// Repository logs repository layer events
func (bl *BusinessLoggerInstance) Repository(operation string, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.String("layer", "repository")}, fields...)
	bl.logger.Info("🗄️ "+operation, allFields...)
}

// Error logs error events in any layer
func (bl *BusinessLoggerInstance) Error(layer, operation string, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("layer", layer),
		zap.Error(err),
	}, fields...)
	bl.logger.Error("❌ "+operation, allFields...)
}

// Performance logs performance metrics
func (bl *BusinessLoggerInstance) Performance(operation string, duration time.Duration, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Duration("duration", duration)}, fields...)
	bl.logger.Info("⚡ "+operation, allFields...)
}

// Debug logs debug information
func (bl *BusinessLoggerInstance) Debug(layer, message string, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.String("layer", layer)}, fields...)
	bl.logger.Debug("🔍 "+message, allFields...)
}