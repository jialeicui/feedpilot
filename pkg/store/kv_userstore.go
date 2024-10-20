package store

import (
	"github.com/jialeicui/feedpilot/pkg/meta"
)

var _ UserStore = (*userStore)(nil)

type userStore struct {
	store KvStore
}

func (u *userStore) Put(id meta.UserID, user *meta.User) error {
	return u.store.Put(id.String(), user)
}

func (u *userStore) Get(id meta.UserID) (*meta.User, error) {
	val, err := u.store.Get(id.String())
	if err != nil {
		return nil, err
	}
	var user = new(meta.User)
	if err := user.Load(val); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userStore) Delete(id meta.UserID) error {
	return u.store.Delete(id.String())
}

func (u *userStore) List(offset, limit int) ([]*meta.User, error) {
	keys, err := u.store.List(offset, limit)
	if err != nil {
		return nil, err
	}
	users := make([]*meta.User, 0, len(keys))
	for _, key := range keys {
		val, err := u.store.Get(key)
		if err != nil {
			return nil, err
		}
		var user = new(meta.User)
		if err := user.Load(val); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func NewUserStore(kv KvStore) UserStore {
	return &userStore{
		store: kv,
	}
}
