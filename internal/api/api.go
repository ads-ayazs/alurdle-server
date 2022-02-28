/*
Package api implements the REST API for alurdle-server.

This alurdle game server is designed to be played by automated bots. However, it
could equally be used by a front-end interface for interactive play.

The interface exposes three routes:

GET /game?id=game_id&word=start_word
	- Creates a new game when called without the "id" param.
	- Optional param "word" can be used to specify the start word for a new game.
	- Optional param "id" can be passed to obtain a JSON game description for
		the given ID.

GET /play?id=GAME_ID&guess=GUESS_WORD
	- Attempt a guess of the secret word for a game.
	- Required param "id" must be a valid game ID that is currently in play.
	- Required param "guess" must be a valid guess word.
	- Returns an error if the ID is invalid or the game not in play.
	- Returns an error if the guess word is not a valid word (due to length, etc.)

GET /resign?id=GAME_ID
	- Resign from a game that is currently in play.
	- Required param "id" must be a valid game ID that is currently in play.

Key functions:
	Initialize() - Sets up the game API router.
*/
package api

import (
	"fmt"
	"net/http"

	"aluance.io/alurdleserver/internal/config"
	"aluance.io/alurdleserver/internal/game"
	"github.com/gin-gonic/gin"
)

// TODO: Constants must be camelCase
const API_RESPONSE_CONTENT_TYPE = "application/json; charset=utf-8"

// Initialize sets up the API router.
func Initialize() {
	setupRouter()
}

// setupRouter creates and configures the GIN engine.
func setupRouter() *gin.Engine {
	router := gin.Default()
	// TODO: Enable security | https://github.com/gin-contrib/secure
	// router.Use(secure.New(secure.DefaultConfig()))

	router.GET("/game", getGame)
	router.GET("/play", getPlay)
	router.GET("/resign", getResign)

	router.Run(fmt.Sprintf(":%d", config.CONFIG_API_PORT))

	return router
}

// getGame handles requests to create a new game and obtain a JSON description
// for a game ID.
func getGame(c *gin.Context) {
	gameId := c.Query("id")
	startWord := c.Query("word")

	var g game.Game
	var err error
	if len(gameId) < 1 {
		g, err = game.Create(startWord)
	} else {
		g, err = game.Retrieve(gameId)
	}
	if handleError(c, err) {
		return
	}

	out, err := g.Describe()
	if handleError(c, err) {
		return
	}

	c.Data(http.StatusOK, API_RESPONSE_CONTENT_TYPE, []byte(out))
}

// getPlay handles game play requests for an existing game ID.
func getPlay(c *gin.Context) {
	gameId := c.Query("id")
	guessWord := c.Query("guess")

	if len(gameId) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	g, err := game.Retrieve(gameId)
	if handleError(c, err) {
		return
	}

	out, err := g.Play(guessWord)
	if err != nil {
		safeErrors := []error{game.ErrGameOver, game.ErrInvalidWord, game.ErrOutOfTurns}
		for _, safe := range safeErrors {
			if err == safe {
				c.Data(http.StatusOK, API_RESPONSE_CONTENT_TYPE, []byte(out))
				return
			}
		}

		handleError(c, err)
		return
	}

	c.Data(http.StatusOK, API_RESPONSE_CONTENT_TYPE, []byte(out))
}

// getResign handles requests to resign from a game that is currently in play.
func getResign(c *gin.Context) {
	gameId := c.Query("id")

	if len(gameId) < 1 {
		handleError(c, ErrInvalidId)
		return
	}
	g, err := game.Retrieve(gameId)
	if handleError(c, err) {
		return
	}

	out, err := g.Resign()
	if handleError(c, err) {
		return
	}

	c.Data(http.StatusOK, API_RESPONSE_CONTENT_TYPE, []byte(out))
}

// handleError provides default error handling across all API routes.
func handleError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return true
	}

	return false
}
