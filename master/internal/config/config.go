package config

import (
	"path"
	"path/filepath"
	"runtime"
)

const CONFIG_DICTIONARY_FILEPATH = "data/google-10000-english-usa-no-swears-medium.txt"
const CONFIG_GAME_WORDLENGTH = 5

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
