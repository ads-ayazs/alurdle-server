package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"aluance.io/wordleplayer/internal/config"
)

func main() {
	out, err := startGame()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}

func startGame() (string, error) {
	resp, err := http.Get(config.CONFIG_SERVER_URL + "/game")
	if err != nil {
		log.Error(err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return string(body), nil
}
