package objects

type Message struct {
	Id string
	ChannelId string
	GuildId string
	Author User
	Member Member
	Content string
	Timestamp string
	EditedTimestamp string
	Tts bool
	MentionEveryone bool
	Mentions []interface{}
	MentionsRoles []int64
	Attachments []Attachment
	Embeds []Embed
	Reactions []Reaction
	Nonce string
	Pinned bool
	WebhookId string
	Type int
	Activity MessageActivity
	Application MessageApplication
}
