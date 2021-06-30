package session

type SessionData struct {
	AccessToken  string `json:"access_token"`
	Expiry       int64  `json:"expiry"`
	RefreshToken string `json:"refresh_token"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar_hash"`
	HasGuilds    bool   `json:"has_guilds"`
}
