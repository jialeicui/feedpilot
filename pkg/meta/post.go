package meta

import (
	"encoding/json"
)

type PostID string

func (p *PostID) String() string {
	return string(*p)
}

type Post struct {
	ID        PostID  `json:"id"`
	Text      string  `json:"text"`
	Medias    []Media `json:"media,omitempty"`
	Links     string  `json:"links,omitempty"`
	Location  string  `json:"location,omitempty"`
	Timestamp string  `json:"timestamp,omitempty"`
}

func (p *Post) String() (string, error) {
	content, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (p *Post) Load(data string) error {
	return json.Unmarshal([]byte(data), p)
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
