package dictionary

import (
	"bufio"
	"math/rand"
	"strings"
	"time"

	"aluance.io/wordle/internal/config"
	"github.com/matryer/resync"
)

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
		// filename = path.Join(config.RootDir(), config.CONFIG_DICTIONARY_FILEPATH)
		filename = config.CONFIG_DICTIONARY_FILEPATH
	}

	// f, err := os.Open(filename)
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

/*
//go:embed data/*
var embFS embed.FS

func initializeEmbeded() (io.Reader, error) {

	r, err := embFS.Open(config.CONFIG_DICTIONARY_FILEPATH)
	if err != nil {
		return nil, err
	}

	return r, nil
}
*/

type dict struct {
	init_once  resync.Once
	initalized bool
	words      []string
	wordMap    map[string]bool
}

func (d *dict) size() int {
	return len(d.words)
}

func (d *dict) reset() {
	d.words = []string{}
	d.wordMap = make(map[string]bool)
	d.init_once.Reset()
	d.initalized = false
}

var wordleDict = &dict{initalized: false, words: []string{}, wordMap: make(map[string]bool)}
