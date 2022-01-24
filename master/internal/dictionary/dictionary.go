package dictionary

import (
	"bufio"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"aluance.io/wordle/internal/config"
	"github.com/matryer/resync"
)

const CONFIG_DICTIONARY_FILEPATH = "data/google-10000-english-usa-no-swears-medium.txt"

func GenerateWord() (string, error) {
	if err := Initialize(""); err != nil {
		return "", err
	}

	word := "blank"
	if max := wordleDict.size(); max > 0 {
		index := rand.Intn(max)
		word = wordleDict.words[index]
	}

	return word, nil
}

func Foo(bar string) string {
	return "foo.bar"
}

func IsWordValid(w string) bool {
	if err := Initialize(""); err != nil {
		return false
	}

	if member, ok := wordleDict.wordMap[strings.ToLower(w)]; ok {
		return member
	}

	return false
}

func Initialize(filename string) error {

	// Only initialized dictionary once
	if wordleDict.initalized {
		return nil
	}

	if len(filename) < 1 {
		filename = path.Join(config.RootDir(), CONFIG_DICTIONARY_FILEPATH)
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Do this only once (unless reset)
	wordleDict.init_once.Do(func() {
		rand.Seed(time.Now().UnixNano())

		// Load only 5-letter words from the file
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			word := scanner.Text()
			if len(word) == 5 {
				wordleDict.words = append(wordleDict.words, word)
				wordleDict.wordMap[word] = true
			}
		}

		if err := scanner.Err(); err != nil {
			return
		}

		wordleDict.initalized = true
	})

	return nil
}

type dict struct {
	init_once  resync.Once
	initalized bool
	words      []string
	wordMap    map[string]bool
}

func (d dict) size() int {
	return len(d.words)
}

func (d *dict) reset() {
	d.words = []string{}
	d.wordMap = make(map[string]bool)
	d.init_once.Reset()
	d.initalized = false
}

var wordleDict = &dict{initalized: false, words: []string{}, wordMap: make(map[string]bool)}