package objects

type Member struct {
	User     User
	Nick     string
	Roles    []int64 `json:"roles,string"`
	JoinedAt string
	Deaf     bool
	Mute     bool
}
