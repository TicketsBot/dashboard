package objects

type User struct {
	Id string
	Username string
	Discriminator string
	Avatar string
	Verified bool
	Email string
	Flags int
	PremiumType int
}
