package meta

type PostID string
type Post struct {
	ID        PostID  `json:"id"`
	Text      string  `json:"text"`
	Medias    []Media `json:"media,omitempty"`
	Links     string  `json:"links,omitempty"`
	Location  string  `json:"location,omitempty"`
	Timestamp string  `json:"timestamp,omitempty"`
}

type PostMeta struct {
	Likes  []UserID `json:"likes"`
	Shares []UserID `json:"shares"`
	Tags   []string `json:"tags"`
}

type CommentID string
type Comment struct {
	ID        CommentID `json:"id"`
	PostID    PostID    `json:"post_id"`
	UserID    UserID    `json:"user_id"`
	Text      string    `json:"text,omitempty"`
	Timestamp string    `json:"timestamp,omitempty"`
}
