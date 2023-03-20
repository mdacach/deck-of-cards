package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
)

func createTestDeck(router *gin.Engine, params string) uuid.UUID {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/deck/new"+params, nil)
	router.ServeHTTP(w, req)

	var createResponse CreateDeckResponse
	// This should never fail.
	_ = json.Unmarshal(w.Body.Bytes(), &createResponse)

	deckID := createResponse.DeckID

	return deckID
}

// TODO: Some way to make this run before each test?
//
// Update: TestMain does not work because we need access to the router created by setup().
func setup() *gin.Engine {
	// Lightweight mode for testing.
	gin.SetMode(gin.TestMode)

	server := NewServer()

	return server.router
}
