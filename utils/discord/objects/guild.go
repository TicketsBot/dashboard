package objects

type Guild struct {
	Id          string
	Name        string
	Icon        string
	Owner       bool
	OwnerId     uint64 `json:"id,string"`
	Permissions int
}
