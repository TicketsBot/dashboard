package utils

import (
	"encoding/json"
	"github.com/apex/log"
	"io/ioutil"
)

var emojisByName map[string]string
var emojis []string

func LoadEmoji() {
	bytes, err := ioutil.ReadFile("emojis.json")
	if err != nil {
		log.Error("Couldn't load emoji: " + err.Error())
		return
	}

	if err := json.Unmarshal(bytes, &emojisByName); err != nil {
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
