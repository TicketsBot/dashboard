package utils

import (
	"github.com/TicketsBot/common/collections"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/guild"
)

func Contains[T comparable](slice []T, value T) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}

	return false
}

func Reverse(slice []message.Message) []message.Message {
	for i := len(slice)/2 - 1; i >= 0; i-- {
		opp := len(slice) - 1 - i
		slice[i], slice[opp] = slice[opp], slice[i]
	}
	return slice
}

func Map[T comparable, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, elem := range slice {
		result[i] = f(elem)
	}

	return result
}

func ToSet[T comparable](slice []T) *collections.Set[T] {
	set := collections.NewSet[T]()

	for _, el := range slice {
		set.Add(el)
	}

	return set
}

func RoleToId(role guild.Role) uint64 {
	return role.Id
}
