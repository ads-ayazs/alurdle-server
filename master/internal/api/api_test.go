package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"aluance.io/wordleserver/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGame(t *testing.T) {
	tests := []struct {
		id     string
		word   string
		errMsg string
	}{
		{id: ""},
		{id: "", word: "happy"},
		{id: "<ID>"},
	}

	assert := assert.New(t)

	router := setupRouter()

	gameId := ""
	for _, test := range tests {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/game", nil)
		assert.NoError(err)

		q := req.URL.Query()
		if len(test.id) > 0 {
			q.Add("id", strings.Replace(test.id, "<ID>", gameId, 1))
		}
		if len(test.word) > 0 {
			q.Add("word", test.word)
		}
		req.URL.RawQuery = q.Encode()

		router.ServeHTTP(w, req)
		assert.Equal(http.StatusOK, w.Code)
		assert.Contains(w.Result().Header["Content-Type"], API_RESPONSE_CONTENT_TYPE)
		assert.NotEmpty(w.Body.Bytes())

		mapResult := map[string]interface{}{}
		assert.NoError(json.Unmarshal(w.Body.Bytes(), &mapResult))

		testElements := []string{"id", "gameStatus", "attempts"}
		for _, elem := range testElements {
			assert.Contains(mapResult, elem)
		}
		if v, ok := mapResult["id"]; ok {
			gameId = v.(string)
		}
	}
}

func TestGetPlay(t *testing.T) {
	tests := []struct {
		id     string
		guess  string
		errMsg string
	}{
		{id: "", guess: "", errMsg: "invalid id"},
		{id: "<ID>", guess: "", errMsg: "invalid guess"},
		{id: "<ID>", guess: "alphabet", errMsg: "invalid guess"},
		{id: "<ID>", guess: "adieu"},
		{id: "<ID>", guess: "handy"},
		{id: "<ID>", guess: "happy"},
		{id: "<ID>", guess: "bless"},
		{id: "<ID>", guess: "grand"},
		{id: "<ID>", guess: "smile"},
		{id: "<ID>", guess: "poems"},
		{id: "<ID>", guess: "imply"},
		{id: "<ID>", guess: "sugar"},
		{id: "<ID>", guess: "whack"},
		{id: "<ID>", guess: "blink"},
		{id: "<ID>", guess: "games"},
		{id: "<ID>", guess: "scent"},
	}
	tests2 := []struct {
		id     string
		guess  string
		errMsg string
	}{
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "xxxxx"},
		{id: "<ID>", guess: "douze"},
		{id: "<ID>", guess: "xxxxx"},
	}
	startWord := "poems"

	assert := assert.New(t)
	require := require.New(t)

	router := setupRouter()

	// Create game
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/game", nil)
	require.NoError(err)
	q := req.URL.Query()
	q.Add("word", startWord)
	req.URL.RawQuery = q.Encode()

	router.ServeHTTP(w, req)
	require.Equal(http.StatusOK, w.Code)
	require.NotEmpty(w.Body.Bytes())

	mapResult := map[string]interface{}{}
	require.NoError(json.Unmarshal(w.Body.Bytes(), &mapResult))
	gameId := mapResult["id"].(string)
	require.NotEmpty(gameId)

	// Test plays
	for _, test := range tests {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/play", nil)
		assert.NoError(err)

		q := req.URL.Query()
		q.Add("id", strings.Replace(test.id, "<ID>", gameId, 1))
		q.Add("guess", test.guess)
		req.URL.RawQuery = q.Encode()

		router.ServeHTTP(w, req)
		if len(test.errMsg) > 0 {
			assert.NotEqual(http.StatusOK, w.Code, test.guess)
			assert.Contains(w.Result().Header["Content-Type"], API_RESPONSE_CONTENT_TYPE)
		} else if assert.Equal(http.StatusOK, w.Code, fmt.Sprintf("\"%s\": %s", test.guess, w.Body.String())) {
			assert.Contains(w.Result().Header["Content-Type"], API_RESPONSE_CONTENT_TYPE)
			assert.NotEmpty(w.Body.Bytes())

			mapResult := map[string]interface{}{}
			assert.NoError(json.Unmarshal(w.Body.Bytes(), &mapResult), test.guess)

			attemptsUsed := int(mapResult["attemptsUsed"].(float64))
			gameStatus := mapResult["gameStatus"].(string)
			if attemptsUsed == config.CONFIG_GAME_MAXATTEMPTS {
				assert.NotEqual("InPlay", gameStatus)
			}
		}
	}

	// Create another game
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/game", nil)
	require.NoError(err)
	q = req.URL.Query()
	q.Add("word", startWord)
	req.URL.RawQuery = q.Encode()

	router.ServeHTTP(w, req)
	require.Equal(http.StatusOK, w.Code)
	require.NotEmpty(w.Body.Bytes())

	mapResult = map[string]interface{}{}
	require.NoError(json.Unmarshal(w.Body.Bytes(), &mapResult))
	gameId = mapResult["id"].(string)
	require.NotEmpty(gameId)

	// Test plays
	for _, test := range tests2 {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/play", nil)
		assert.NoError(err)

		q := req.URL.Query()
		q.Add("id", strings.Replace(test.id, "<ID>", gameId, 1))
		q.Add("guess", test.guess)
		req.URL.RawQuery = q.Encode()

		router.ServeHTTP(w, req)
		if len(test.errMsg) > 0 {
			assert.NotEqual(http.StatusOK, w.Code, test.guess)
			assert.Contains(w.Result().Header["Content-Type"], API_RESPONSE_CONTENT_TYPE)
		} else if assert.Equal(http.StatusOK, w.Code, fmt.Sprintf("\"%s\": %s", test.guess, w.Body.String())) {
			assert.Contains(w.Result().Header["Content-Type"], API_RESPONSE_CONTENT_TYPE)
			assert.NotEmpty(w.Body.Bytes())

			mapResult := map[string]interface{}{}
			assert.NoError(json.Unmarshal(w.Body.Bytes(), &mapResult), test.guess)

			attemptsUsed := int(mapResult["attemptsUsed"].(float64))
			gameStatus := mapResult["gameStatus"].(string)
			if attemptsUsed == config.CONFIG_GAME_MAXATTEMPTS {
				assert.NotEqual("InPlay", gameStatus)
			}
		}
	}
}

func TestGetResign(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	router := setupRouter()

	// Create game
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/game", nil)
	router.ServeHTTP(w, req)
	require.Equal(http.StatusOK, w.Code)
	require.NotEmpty(w.Body.Bytes())

	mapResult := map[string]interface{}{}
	require.NoError(json.Unmarshal(w.Body.Bytes(), &mapResult))

	v, ok := mapResult["id"]
	require.True(ok)
	gameId := v.(string)

	// test Resign
	w = httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/resign", nil)
	require.NoError(err)
	q := req.URL.Query()
	q.Add("id", gameId)
	req.URL.RawQuery = q.Encode()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code)
	assert.Contains(w.Result().Header["Content-Type"], API_RESPONSE_CONTENT_TYPE)
	assert.NotEmpty(w.Body.Bytes())

	mapResult = map[string]interface{}{}
	require.NoError(json.Unmarshal(w.Body.Bytes(), &mapResult))
	testElements := []string{"attemptsUsed", "gameStatus"}
	for _, elem := range testElements {
		assert.Contains(mapResult, elem)
	}
	assert.EqualValues("Resigned", mapResult["gameStatus"])
}
