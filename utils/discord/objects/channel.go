package objects

type Channel struct {
	Id string
	Type int
	GuildId string
	Position int
	PermissionsOverwrites []Overwrite
	Name string
	Topic string
	Nsfw bool
	LastMessageId string
	Bitrate int
	userLimit int
	RateLimitPerUser int
	Recipients []User
	Icon string
	Ownerid string
	ApplicationId string
	ParentId string
	LastPinTimestamp string
}
