package utils

import (
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

func Reverse(slice []message.Message) []message.Message {
	for i := len(slice)/2-1; i >= 0; i-- {
		opp := len(slice)-1-i
		slice[i], slice[opp] = slice[opp], slice[i]
	}
	return slice
}
