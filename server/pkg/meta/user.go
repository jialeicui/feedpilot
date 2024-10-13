package meta

type UserID string

type User struct {
	ID             UserID `json:"id"`
	Username       string `json:"username,omitempty"`
	DisplayName    string `json:"display_name,omitempty"`
	ProfilePic     string `json:"profile_pic,omitempty"`
	FollowersCount int    `json:"followers_count,omitempty"`
	FollowingCount int    `json:"following_count,omitempty"`
	Bio            string `json:"bio,omitempty"`
}
