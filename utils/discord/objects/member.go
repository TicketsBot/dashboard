package objects

type Member struct {
	User     User
	Nick     string
	Roles    []string
	JoinedAt string
	Deaf     bool
	Mute     bool
}
