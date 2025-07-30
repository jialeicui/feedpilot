package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/jialeicui/feedpilot/pkg/meta"
)

func TestPostStore(t *testing.T) {
	var (
		ctl  = gomock.NewController(t)
		kv   = NewMockKvStore(ctl)
		must = require.New(t)
	)

	ps := NewPostStore(kv)
	must.NotNil(ps)

	// Test Put
	post := &meta.Post{
		ID:   "1",
		Text: "Hello world",
	}
	kv.EXPECT().Put("1", post).Return(nil)
	must.NoError(ps.Put("1", post))

	// Test Get
	kv.EXPECT().Get("1").Return(`{"id":"1","text":"Hello world"}`, nil)
	p, err := ps.Get("1")
	must.NoError(err)
	must.Equal(post, p)

	// Test Delete
	kv.EXPECT().Delete("1").Return(nil)
	must.NoError(ps.Delete("1"))

	// Test List
	kv.EXPECT().List(0, 0).Return([]string{"1"}, nil)
	kv.EXPECT().Get("1").Return(`{"id":"1","text":"Hello world"}`, nil)
	posts, err := ps.List(0, 0)
	must.NoError(err)
	must.Equal([]*meta.Post{post}, posts)

	// Test Search
	kv.EXPECT().Search("Hello", 0, 10).Return([]string{"1"}, nil)
	kv.EXPECT().Get("1").Return(`{"id":"1","text":"Hello world"}`, nil)
	posts, err = ps.Search("Hello", 0, 10)
	must.NoError(err)
	must.Equal([]*meta.Post{post}, posts)
}