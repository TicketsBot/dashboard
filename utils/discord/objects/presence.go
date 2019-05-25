package objects

type Presence struct {
	User User
	Roles []string
	Game Activity
	GuildId string
	Status string
	Activities []Activity
	ClientStatus ClientStatus
}
