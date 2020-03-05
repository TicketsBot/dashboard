package objects

type AllowedMention struct {
	Parse []AllowedMentionType
	Roles []string
	Users []string
}
