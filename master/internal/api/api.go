package api

import (
	"net/http"

	"aluance.io/wordle/master/internal/game"
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

	var g game.Game
	if len(gameId) < 1 {
		g, _ = game.Create("blank")
	} else {
		g, _ = game.Retrieve(gameId)
	}

	out, _ := g.Describe()

	c.String(http.StatusOK, out)
}

func getPlay(c *gin.Context) {
	gameId := c.Query("id")
	guessWord := c.Query("guess")

	g, err := game.Retrieve(gameId)
	if err != nil {
		c.String(http.StatusFailedDependency, "game.Retrieve() failed")
	}
	out, err := g.Play(guessWord)
	if err != nil {
		c.String(http.StatusFailedDependency, "game.Play() failed")
	}

	c.String(http.StatusOK, out)
}

func getResign(c *gin.Context) {
	gameId := c.Query("id")

	g, err := game.Retrieve(gameId)
	if err != nil {
		c.String(http.StatusFailedDependency, "game.Retrieve() failed")
	}
	out, err := g.Resign()
	if err != nil {
		c.String(http.StatusFailedDependency, "game.Resign() failed")
	}

	c.String(http.StatusOK, out)
}
