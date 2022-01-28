package config

import (
	"embed"
	"io/fs"
	"path"
	"path/filepath"
	"runtime"
)

const CONFIG_API_PORT = 8080

// const CONFIG_DICTIONARY_FILENAME = "google-10000-english-usa-no-swears-medium.txt"
const CONFIG_DICTIONARY_FILENAME = "corncob_lowercase.txt"
const CONFIG_DICTIONARY_FILEPATH = "data/" + CONFIG_DICTIONARY_FILENAME
const CONFIG_GAME_WORDLENGTH = 5
const CONFIG_GAME_MAXATTEMPTS = 12
const CONFIG_GAME_MAXVALIDATTEMPTS = 6

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

//go:embed data/*
var embFS embed.FS

func LoadEmbedFile(fp string) (fs.File, error) {
	if len(fp) < 1 {
		return nil, ErrFilepath
	}

	f, err := embFS.Open(fp)
	if err != nil {
		return nil, err
	}

	return f, nil
}
