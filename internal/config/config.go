/*
Package config centralizes shared configuration across the entire application.

Key functions:
	LoadEmbedFile(fp string) - Loads a file by name from data embedded in the executable.
	RootDir() - Returns the root folder from where the binary executable resides.
*/
package config

import (
	"embed"
	"io/fs"
	"path"
	"path/filepath"
	"runtime"
)

// Port where the API is served.
const CONFIG_API_PORT = 8080

// const CONFIG_DICTIONARY_FILENAME = "google-10000-english-usa-no-swears-medium.txt"

// Filename of the default dictionary.
const CONFIG_DICTIONARY_FILENAME = "corncob_lowercase.txt"

// Filepath of the default dictionary.
const CONFIG_DICTIONARY_FILEPATH = "data/" + CONFIG_DICTIONARY_FILENAME

// Game word length.
const CONFIG_GAME_WORDLENGTH = 5

// Game maximum permitted attempts.
const CONFIG_GAME_MAXATTEMPTS = 12

// Game maximum valid attempts.
const CONFIG_GAME_MAXVALIDATTEMPTS = 6

// RootDir returns the runtime root folder.
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

//go:embed data/*
var embFS embed.FS

// LoadEmbedFile loads the given file path from the data embedded in the binary.
func LoadEmbedFile(fp string) (fs.File, error) {
	// Validate fp
	if len(fp) < 1 {
		return nil, ErrFilepath
	}

	// Open the file from the embedded file system.
	f, err := embFS.Open(fp)
	if err != nil {
		return nil, err
	}

	return f, nil
}
