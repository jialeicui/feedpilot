package meta

import (
	"encoding/json"
)

type UserID string

func (u *UserID) String() string {
	return string(*u)
}

type User struct {
	ID             UserID `json:"id"`
	Username       string `json:"username,omitempty"`
	DisplayName    string `json:"display_name,omitempty"`
	ProfilePic     string `json:"profile_pic,omitempty"`
	FollowersCount int    `json:"followers_count,omitempty"`
	FollowingCount int    `json:"following_count,omitempty"`
	Bio            string `json:"bio,omitempty"`
}

func (u *User) String() (string, error) {
	content, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (u *User) Load(data string) error {
	return json.Unmarshal([]byte(data), u)
}
