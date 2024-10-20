package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/jialeicui/feedpilot/pkg/meta"
)

func TestUserStore(t *testing.T) {
	var (
		ctl  = gomock.NewController(t)
		kv   = NewMockKvStore(ctl)
		must = require.New(t)
	)

	us := NewUserStore(kv)
	must.NotNil(us)

	// Test Put
	user := &meta.User{
		ID:       "1",
		Username: "user1",
	}
	kv.EXPECT().Put("1", user).Return(nil)
	must.NoError(us.Put("1", user))

	// Test Get
	kv.EXPECT().Get("1").Return(`{"id":"1","username":"user1"}`, nil)
	u, err := us.Get("1")
	must.NoError(err)
	must.Equal(user, u)

	// Test Delete
	kv.EXPECT().Delete("1").Return(nil)
	must.NoError(us.Delete("1"))

	// Test List
	kv.EXPECT().List(0, 0).Return([]string{"1"}, nil)
	kv.EXPECT().Get("1").Return(`{"id":"1","username":"user1"}`, nil)
	users, err := us.List(0, 0)
	must.NoError(err)
	must.Equal([]*meta.User{user}, users)
}
