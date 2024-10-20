package weibo

import (
	"github.com/jialeicui/feedpilot/pkg/platform"
	"github.com/jialeicui/feedpilot/pkg/store"
)

var _ platform.Platform = (*Weibo)(nil)

type Weibo struct {
	userStore store.UserStore
	postStore store.PostStore
}

func New(user store.UserStore, post store.PostStore) *Weibo {
	return &Weibo{
		userStore: user,
		postStore: post,
	}
}

func (w *Weibo) Sync() error {
	//TODO implement me
	panic("implement me")
}

func (w *Weibo) Close() error {
	//TODO implement me
	panic("implement me")
}
