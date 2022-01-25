package api

import (
	"fmt"
	"net/http"

	"aluance.io/wordle/internal/game"
	"github.com/gin-gonic/gin"
)

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

	router.Run(":8080")

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

	c.String(http.StatusOK, out)
}

func getPlay(c *gin.Context) {
	gameId := c.Query("id")
	guessWord := c.Query("guess")

	if len(gameId) < 1 {
		c.String(http.StatusBadRequest, "invalid ID")
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
				c.String(http.StatusOK, out)
				return
			}
		}

		handleError(c, err)
		return
	}

	c.String(http.StatusOK, out)
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

	c.String(http.StatusOK, out)
}

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("{\"error\": \"%s\"}", err.Error()))
		return true
	}

	return false
}
