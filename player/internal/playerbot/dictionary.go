package playerbot

import (
	"bytes"
	"math/rand"
	"time"

	"aluance.io/wordleplayer/internal/config"
)

type PlayerDictionary interface {
	Remember(string, bool) error
	Generate() (string, error)
	IsValid(string) bool
}

type playerDictionary struct {
	lexicon map[string]bool
}

func (d *playerDictionary) initialize() {
	d.lexicon = make(map[string]bool)
}

func CreateDictionary() PlayerDictionary {
	rand.Seed(time.Now().UnixNano())

	d := new(playerDictionary)
	d.initialize()

	return d
}

func (d *playerDictionary) Remember(word string, valid bool) error {
	d.lexicon[word] = valid
	return nil
}

func (d playerDictionary) Generate() (string, error) {
	buf := bytes.NewBufferString("")

	for i := 0; i < config.CONFIG_GAME_WORDLENGTH; i++ {
		letterNum := rand.Intn(26) + 65
		c := string(rune(letterNum))
		buf.WriteString(c)
	}

	return buf.String(), nil
}

func (d playerDictionary) IsValid(word string) bool {
	if validity, ok := d.lexicon[word]; ok {
		return validity
	}

	return false
}
