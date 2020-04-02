package objects

type Embed struct {
	Title string
	Type string
	Description string
	Url string
	Timestamp string
	Color int
	Footer EmbedField
	Image EmbedImage
	Thumbnail EmbedThumbnail
	Video EmbedVideo
	Provider EmbedProvider
	Author EmbedAuthor
	Fields []EmbedField
}
