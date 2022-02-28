/*
Package dictionary implements a thread-safe dictionary used to generate and
validate game words.

Key functions:
	GenerateWord() - Returns a random word from the dictionary.
	Initialize(filename string) - Loads a dictionary from the given file name.
	IsWordValid(w string) - Returns true if the given word exists in the dictionary.
*/
package dictionary

import (
	"bufio"
	"math/rand"
	"strings"
	"sync"
	"time"

	"aluance.io/wordleserver/internal/config"
	"github.com/matryer/resync"
)

// GenerateWord returns a random word from the dictionary.
//
// Returns an error if the dictionary fails to initialize.
func GenerateWord() (string, error) {
	// Initialize the dictionary
	if err := Initialize(""); err != nil {
		return "", err
	}

	// Obtain a read lock
	wordleDict.mu.RLock()
	defer wordleDict.mu.RUnlock()

	// Return a random word or "blank" by default.
	word := "blank"
	if max := wordleDict.size(); max > 0 {
		index := rand.Intn(max)
		word = wordleDict.words[index]
	}

	return word, nil
}

// IsWordValid returns true if the given word exists in the dictionary.
func IsWordValid(w string) bool {
	// Initialize the dictionary
	if err := Initialize(""); err != nil {
		return false
	}

	// Ontain a read lock
	wordleDict.mu.RLock()
	defer wordleDict.mu.RUnlock()

	// Check if the word is in the dictionary
	if member, ok := wordleDict.wordMap[strings.ToLower(w)]; ok {
		return member
	}

	return false
}

// Initialize loads a dictionary from the given file name.
//
// When filename is empty, it attempts to load the default dictionary. The
// function ensures that a dictionary is only loaded once. Returns an error if
// it is unable to load from the file.
func Initialize(filename string) error {
	// Ontain a write lock
	wordleDict.mu.Lock()
	defer wordleDict.mu.Unlock()

	// Only initialize dictionary once
	if wordleDict.initalized {
		return nil
	}

	// Open the dictonary file
	if len(filename) < 1 {
		filename = config.CONFIG_DICTIONARY_FILEPATH
	}

	f, err := config.LoadEmbedFile(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Do this only once (unless reset)
	wordleDict.init_once.Do(func() {
		rand.Seed(time.Now().UnixNano())

		// Load only words of configured length from the file
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			word := scanner.Text()
			if len(word) == config.CONFIG_GAME_WORDLENGTH {
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

// dict struct stores the entire dictionary in memory.
type dict struct {
	init_once  resync.Once
	initalized bool
	mu         sync.RWMutex
	words      []string
	wordMap    map[string]bool
}

// size returns the number of words in the dictionary
func (d *dict) size() int {
	return len(d.words)
}

// reset clears and resets the loaded dictionary for testing.
func (d *dict) reset() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.words = []string{}
	d.wordMap = make(map[string]bool)
	d.init_once.Reset()
	d.initalized = false
}

// wordleDict is the singleton dict variable.
var wordleDict = &dict{initalized: false, words: []string{}, wordMap: make(map[string]bool)}
