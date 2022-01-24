package api

import (
	"net/http"

	"aluance.io/wordle/internal/game"
	"github.com/gin-gonic/gin"
)

func Initialize() {
	setupRouter()
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/game", getGame)
	router.GET("/play", getPlay)
	router.GET("/resign", getResign)

	router.Run(":8080")

	return router
}

func getGame(c *gin.Context) {
	gameId := c.Query("id")
	startWord := c.Query("word")

	if len(startWord) < 1 {
		startWord = "blank"
	}

	var g game.Game
	if len(gameId) < 1 {
		g, _ = game.Create(startWord)
	} else {
		g, _ = game.Retrieve(gameId)
	}

	out, _ := g.Describe()

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
	if err != nil {
		c.String(http.StatusBadRequest, "incorrect ID - game.Retrieve() failed")
		return
	}
	out, err := g.Play(guessWord)
	if err != nil {
		c.String(http.StatusBadRequest, "guess word error - game.Play() failed")
		return
	}

	c.String(http.StatusOK, out)
}

func getResign(c *gin.Context) {
	gameId := c.Query("id")

	if len(gameId) < 1 {
		c.String(http.StatusBadRequest, "{ \"error\": \"invalid id\" }")
		return
	}
	g, err := game.Retrieve(gameId)
	if err != nil {
		c.String(http.StatusBadRequest, "{ \"error\": \"game.Retrieve() failed\" }")
		return
	}
	out, err := g.Resign()
	if err != nil {
		c.String(http.StatusInternalServerError, "{ \"error\": \"game.Resign() failed\" }")
		return
	}

	c.String(http.StatusOK, out)
}