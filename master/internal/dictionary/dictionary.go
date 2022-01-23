package dictionary

import (
	"bufio"
	"fmt"
	"os"
)

func GenerateWord() (string, error) {
	return "", nil
}

func IsWordValid(w string) bool {
	return false
}

func Initialize(file string) error {
	if len(file) < 1 {
		return fmt.Errorf("invalid filename")
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		word := scanner.Text()
		if len(word) == 5 {
			wordleDict.words = append(wordleDict.words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type dict struct {
	words []string
}

var wordleDict = &dict{}
