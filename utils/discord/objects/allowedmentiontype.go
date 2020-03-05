package objects

type AllowedMentionType string

const(
	EVERYONE AllowedMentionType = "everyone"
	USERS AllowedMentionType = "users"
	ROLES AllowedMentionType = "roles"
)
