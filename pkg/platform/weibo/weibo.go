package weibo

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/jialeicui/feedpilot/pkg/platform"
	"github.com/jialeicui/feedpilot/pkg/store"
)

var _ platform.Platform = (*Weibo)(nil)

type Weibo struct {
	ctx    context.Context
	cancel context.CancelFunc
	store  store.Store
}

func (w *Weibo) Name() string {
	return "weibo"
}

func (w *Weibo) Sync() error {
	return nil
}

func (w *Weibo) RegisterRoutes(router gin.IRouter) {
}

func (w *Weibo) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}

func New(store store.Store) *Weibo {
	ctx, cancel := context.WithCancel(context.Background())
	return &Weibo{
		ctx:    ctx,
		cancel: cancel,
		store:  store,
	}
}
