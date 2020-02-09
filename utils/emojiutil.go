package utils

import (
	"encoding/json"
	"github.com/apex/log"
	"io/ioutil"
)

var emojis map[string]interface{}

func LoadEmoji() {
	bytes, err := ioutil.ReadFile("emojis.json"); if err != nil {
		log.Error("Couldn't load emoji: " + err.Error())
		return
	}

	if err := json.Unmarshal(bytes, &emojis); err != nil {
		log.Error("Couldn't load emoji: " + err.Error())
		return
	}
}

func GetEmojiByName(name string) string {
	emoji, ok := emojis[name]; if !ok {
		return ""
	}

	str, _ := emoji.(string)
	return str
}
