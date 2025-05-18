package store

import (
	"github.com/jialeicui/feedpilot/pkg/meta"
)

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -typed
type Stringer interface {
	String() (string, error)
	Load(string) error
}

type UserStore interface {
	Put(ID meta.UserID, user *meta.User) error
	Get(ID meta.UserID) (*meta.User, error)
	Delete(ID meta.UserID) error
	// List returns a list of users with offset and limit
	// offset is the starting index of the list
	// limit is the maximum number of users to return, 0 means no limit
	List(offset, limit int) ([]*meta.User, error)
}

type PostStore interface {
	Put(ID meta.PostID, post *meta.Post) error
	Get(ID meta.PostID) (*meta.Post, error)
	Delete(ID meta.PostID) error
	List(offset, limit int) ([]*meta.Post, error)
}

type ObjectStore interface {
	Put(key string, value []byte)
	Get(key string) []byte
}

type Store struct {
	UserStore   UserStore
	PostStore   PostStore
	ObjectStore ObjectStore
}

type KvStore interface {
	Put(key string, value Stringer) error
	Get(key string) (string, error)
	Delete(key string) error
	List(offset, limit int) ([]string, error)
}
