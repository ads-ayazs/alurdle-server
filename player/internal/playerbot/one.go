package playerbot

import (
	"encoding/json"

	"github.com/rs/xid"
)

const ONEBOT_NAME = "one"

type oneBot struct {
	id         string
	game       oneGame
	dictionary PlayerDictionary
}

func createOne() (Playerbot, error) {
	bot := new(oneBot)
	bot.dictionary = CreateDictionary()

	bot.id = xid.New().String()

	return bot, nil
}

func (b oneBot) PlayGame(ch *chan string) {
	// Avoid sending to nil channel
	if ch == nil {
		return
	}

	// Start the game
	if err := b.startGame(); err != nil {
		*ch <- ""
		return
	}

	// Play while the game remains "InPlay"
	for b.isGameInPlay() {
		if err := b.playTurn(); err != nil {
			*ch <- ""
			return
		}
	}

	// Finish the game and write the output to ch
	*ch <- b.finishGame()
}

type oneGame struct {
	gameId     string
	gameStatus string
	turns      []oneTurn
}

type oneTurn struct {
	guess   string
	isValid bool
}

func (bot *oneBot) startGame() error {
	ge := GetGameEngine()

	// Create a new game and save the game id
	out, err := ge.NewGame()
	if err != nil {
		return err
	}

	// Unmarshall the output
	outmap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(out), &outmap); err != nil {
		return err
	}

	// Save essential information
	bot.game.gameId = outmap["id"].(string)
	bot.game.gameStatus = outmap["gameStatus"].(string)

	return nil
}

func (bot oneBot) isGameInPlay() bool {
	return bot.game.gameStatus == "InPlay"
}

func (bot *oneBot) playTurn() error {
	ge := GetGameEngine()

	// Generate a word
	guessWord, err := bot.dictionary.Generate()
	if err != nil {
		return err
	}

	// Play the guess word
	out, err := ge.PlayTurn(bot.game.gameId, guessWord)
	if err != nil {
		return err
	}

	// Unmarshall the output
	outmap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(out), &outmap); err != nil {
		return err
	}

	// Find the latest attempt
	attempts := outmap["attempts"].([]interface{})
	lastAttempt := attempts[len(attempts)-1].(map[string]interface{})
	tw := lastAttempt["tryWord"].(string)
	if tw != guessWord {
		return ErrFailedAttempt
	}

	// Create and save a turn record
	turn := new(oneTurn)
	turn.guess = guessWord
	turn.isValid = lastAttempt["isValidWord"] == "true"
	bot.game.turns = append(bot.game.turns, *turn)

	// Update the dictionary
	if err := bot.dictionary.Remember(turn.guess, turn.isValid); err != nil {
		return err
	}

	// Save essential information
	bot.game.gameStatus = outmap["gameStatus"].(string)

	return nil
}

func (bot oneBot) finishGame() string {
	return bot.game.gameId
}
