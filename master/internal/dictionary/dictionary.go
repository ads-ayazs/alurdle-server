package dictionary

import (
	"bufio"
	"math/rand"
	"os"
)

const CONFIG_DICTIONARY_FILENAME = "google-10000-english-usa-no-swears-medium.txt"

func GenerateWord() (string, error) {
	word := "blank"
	if max := wordleDict.size(); max > 0 {
		index := rand.Intn(max)
		word = wordleDict.words[index]
	}

	return word, nil
}

func IsWordValid(w string) bool {
	if member, ok := wordleDict.wordMap[w]; ok {
		return member
	}

	return false
}

func Initialize(filename string) error {
	if len(filename) < 1 {
		filename = CONFIG_DICTIONARY_FILENAME
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	wordleDict.words = []string{}
	wordleDict.wordMap = make(map[string]bool)

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
		return err
	}

	return nil
}

type dict struct {
	words   []string
	wordMap map[string]bool
}

func (d dict) size() int {
	return len(d.words)
}

var wordleDict = &dict{}
