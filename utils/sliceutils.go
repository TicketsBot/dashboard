package utils

import (
	"github.com/rxdn/gdl/objects/channel/message"
)

func ContainsString(slice []string, target string) bool {
	for _, elem := range slice {
		if elem == target {
			return true
		}
	}

	return false
}

func Reverse(slice []message.Message) []message.Message {
	for i := len(slice)/2-1; i >= 0; i-- {
		opp := len(slice)-1-i
		slice[i], slice[opp] = slice[opp], slice[i]
	}
	return slice
}
