package api

import (
	"encoding/json"
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

	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code)
	assert.NotEmpty(w.Body.Bytes())

	mapResult := map[string]interface{}{}
	assert.NoError(json.Unmarshal(w.Body.Bytes(), &mapResult))

	testElements := []string{"Id", "Status", "SecretWord", "Attempts"}
	for _, elem := range testElements {
		assert.Contains(mapResult, elem)
	}
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

func TestGetResign(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	params := url.Values{}
	params.Add("id", "xxxxx")
	req, _ := http.NewRequest("GET", "/resign", strings.NewReader(params.Encode()))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
