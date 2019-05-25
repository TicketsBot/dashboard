package objects

type VoiceState struct {
	GuildId string
	ChannelId string
	UserId string
	Member Member
	SessionId string
	Deaf bool
	Mute bool
	SelfDeaf bool
	SelfMute bool
	Suppress bool
}
