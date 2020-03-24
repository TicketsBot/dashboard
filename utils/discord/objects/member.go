package objects

import "github.com/TicketsBot/GoPanel/utils/types"

type Member struct {
	User     User
	Nick     string
	Roles    types.Int64StringSlice `json:"roles,string"`
	JoinedAt string
	Deaf     bool
	Mute     bool
}
