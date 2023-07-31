package utils

import (
	_ "embed"
	"encoding/json"
	"github.com/apex/log"
)

//go:embed emojis.json
var emojisFile []byte

var emojisByName map[string]string
var emojis []string

func LoadEmoji() {
	if err := json.Unmarshal(emojisFile, &emojisByName); err != nil {
		log.Error("Couldn't load emoji: " + err.Error())
		return
	}

	emojis = make([]string, len(emojisByName))
	i := 0
	for _, emoji := range emojisByName {
		emojis[i] = emoji
		i++
	}
}

func GetEmoji(input string) (emoji string, ok bool) {
	// try by name first
	emoji, ok = emojisByName[input]
	if !ok { // else try by the actual unicode char
		for _, unicode := range emojis { // TODO: Optimise
			if unicode == input {
				emoji = unicode
				ok = true
				break
			}
		}
	}

	return
}
