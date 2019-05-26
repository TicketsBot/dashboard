package objects

type Role struct {
	Id          string
	Name        string
	Color       int
	Hoist       bool
	Position    int
	Permissions int
	Managed     bool
	Mentionable bool
}
