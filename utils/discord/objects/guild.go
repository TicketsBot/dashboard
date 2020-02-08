package objects

type Guild struct {
	Id                          string
	Name                        string
	Icon                        string
	Splash                      string
	Owner                       bool
	OwnerId                     string `json:"owner_id"`
	Permissions                 int
	Region                      string
	AfkChannelid                string `json:"afk_channel_id"`
	AfkTimeout                  int
	EmbedEnabled                bool   `json:"embed_enabled"`
	EmbedChannelId              string `json:"embed_channel_id"`
	VerificationLevel           int    `json:"verification_level"`
	DefaultMessageNotifications int    `json:"default_message_notifications"`
	ExplicitContentFilter       int    `json:"explicit_content_filter"`
	Roles                       []Role
	Emojis                      []Emoji
	Features                    []string
	MfaLevel                    int    `json:"mfa_level"`
	ApplicationId               string `json:"application_id"`
	WidgetEnabled               bool   `json:"widget_enabled"`
	WidgetChannelId             string `json:"widget_channel_id"`
	SystemChannelId             string `json:"system_channel_id"`
	JoinedAt                    string `json:"joined_at"`
	Large                       bool
	Unavailable                 bool
	MemberCount                 int `json:"member_count"`
	VoiceStates                 []VoiceState
	Members                     []Member
	Channels                    []Channel
	Presences                   []Presence
	MaxPresences                int    `json:"max_presences"`
	MaxMembers                  int    `json:"max_members"`
	VanityUrlCode               string `json:"vanity_url_code"`
	Description                 string
	Banner                      string
	PremiumTier                 int    `json:"premium_tier"`
	PremiumSubscriptionCount    int    `json:"premium_subscription_count"`
	PreferredLocale             string `json:"preferred_locale"`
}
