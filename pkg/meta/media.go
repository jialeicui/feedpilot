package meta

type MediaBase struct {
	ID   MediaID   `json:"id"`
	Link string    `json:"link"`
	Type MediaType `json:"type"`
}

type MediaID string
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)
