package controllers_test

import (
	"net/http"
	"testing"

	"github.com/maximicciullo/personal-finance-api/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HealthControllerTestSuite struct {
	suite.Suite
	server *test.TestServer
}

func (suite *HealthControllerTestSuite) SetupTest() {
	suite.server = test.NewTestServer()
}

func (suite *HealthControllerTestSuite) TestHealthCheck_Success() {
	// When
	w := suite.server.MakeRequest("GET", "/health", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Check required fields are present
	test.AssertJSONContains(suite.T(), w, map[string]interface{}{
		"status":  "healthy",
		"service": "personal-finance-api",
		"version": "1.0.0",
	})

	// Check that timestamp and uptime fields exist
	response := test.GetResponseJSON(suite.T(), w)
	assert.Contains(suite.T(), response, "timestamp")
	assert.Contains(suite.T(), response, "uptime")
	assert.Equal(suite.T(), "running", response["uptime"])
}

func (suite *HealthControllerTestSuite) TestHealthCheck_ResponseFormat() {
	// When
	w := suite.server.MakeRequest("GET", "/health", nil)

	// Then
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// Verify JSON structure
	response := test.GetResponseJSON(suite.T(), w)
	assert.Len(suite.T(), response, 5, "Health response should have exactly 5 fields")
}

func (suite *HealthControllerTestSuite) TestHealthCheck_MultipleRequests() {
	// Test that multiple health checks work consistently
	for i := 0; i < 3; i++ {
		w := suite.server.MakeRequest("GET", "/health", nil)
		assert.Equal(suite.T(), http.StatusOK, w.Code)
		
		test.AssertJSONContains(suite.T(), w, map[string]interface{}{
			"status": "healthy",
		})
	}
}

func TestHealthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(HealthControllerTestSuite))
}