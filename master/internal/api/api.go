package api

import (
	"fmt"
	"net/http"

	"aluance.io/wordleserver/internal/config"
	"aluance.io/wordleserver/internal/game"
	"github.com/gin-gonic/gin"
)

const API_RESPONSE_CONTENT_TYPE = "application/json; charset=utf-8"

func Initialize() {
	setupRouter()
}

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

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return true
	}

	return false
}
