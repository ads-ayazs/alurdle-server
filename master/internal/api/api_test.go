package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGame(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"Id\":\"c7llkdvr2g4ksqso2fp0\",\"Status\":0,\"SecretWord\":\"BLANK\",\"Attempts\":[]}", w.Body.String())
}

func TestGetPlay(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	params := url.Values{}
	params.Add("id", "xxxxx")
	params.Add("guess", "blues")
	req, _ := http.NewRequest("GET", "/play", strings.NewReader(params.Encode()))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
