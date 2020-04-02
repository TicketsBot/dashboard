package objects

type Activity struct {
	Name          string
	Type          int
	Url           string
	Timestamps    Timestamp
	ApplicationId string
	Details       string
	State         string
	Party         Party
	Assets        Asset
	Secrets       Secret
	Instance      bool
	Flags         int
}
