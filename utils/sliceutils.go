package utils

import (
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/rxdn/gdl/objects/channel/message"
	"reflect"
)

func Contains(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {

			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}

	return false
}

func Insert(slice []objects.Guild, index int, value objects.Guild) []objects.Guild {
	// Grow the slice by one element.
	slice = slice[0 : len(slice)+1]
	// Use copy to move the upper part of the slice out of the way and open a hole.
	copy(slice[index+1:], slice[index:])
	// Store the new value.
	slice[index] = value
	// Return the result.
	return slice
}

func Reverse(slice []message.Message) []message.Message {
	for i := len(slice)/2-1; i >= 0; i-- {
		opp := len(slice)-1-i
		slice[i], slice[opp] = slice[opp], slice[i]
	}
	return slice
}
