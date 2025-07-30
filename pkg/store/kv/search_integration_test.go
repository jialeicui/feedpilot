package kv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jialeicui/feedpilot/pkg/meta"
	"github.com/jialeicui/feedpilot/pkg/store"
)

func TestSearchIntegration(t *testing.T) {
	var (
		must = require.New(t)
		dir  = t.TempDir()
	)

	// Create a badger store
	badgerStore, err := NewBadger(dir)
	must.NoError(err)

	// Create user and post stores
	userStore := store.NewUserStore(badgerStore.WithBucket("users"))
	postStore := store.NewPostStore(badgerStore.WithBucket("posts"))

	// Insert test users
	users := []*meta.User{
		{ID: "1", Username: "alice", DisplayName: "Alice Smith", Bio: "Software developer from NYC"},
		{ID: "2", Username: "bob", DisplayName: "Bob Jones", Bio: "UI/UX designer and artist"},
		{ID: "3", Username: "charlie", DisplayName: "Charlie Brown", Bio: "Product manager and entrepreneur"},
	}

	for _, user := range users {
		err := userStore.Put(user.ID, user)
		must.NoError(err)
	}

	// Insert test posts
	posts := []*meta.Post{
		{ID: "1", Text: "Hello world! This is my first post about programming"},
		{ID: "2", Text: "Beautiful sunset today. Nature is amazing!"},
		{ID: "3", Text: "Working on a new feature for our app. Excited to ship it!"},
		{ID: "4", Text: "Coffee and code make the perfect combination"},
	}

	for _, post := range posts {
		err := postStore.Put(post.ID, post)
		must.NoError(err)
	}

	// Test user search by username
	searchResults, err := userStore.Search("alice", 0, 10)
	must.NoError(err)
	must.Len(searchResults, 1)
	must.Equal("alice", searchResults[0].Username)

	// Test user search by display name
	searchResults, err = userStore.Search("Bob Jones", 0, 10)
	must.NoError(err)
	must.Len(searchResults, 1)
	must.Equal("bob", searchResults[0].Username)

	// Test user search by bio content
	searchResults, err = userStore.Search("developer", 0, 10)
	must.NoError(err)
	must.Len(searchResults, 1)
	must.Equal("alice", searchResults[0].Username)

	// Test user search case insensitive
	searchResults, err = userStore.Search("DESIGNER", 0, 10)
	must.NoError(err)
	must.Len(searchResults, 1)
	must.Equal("bob", searchResults[0].Username)

	// Test post search by content
	postResults, err := postStore.Search("programming", 0, 10)
	must.NoError(err)
	must.Len(postResults, 1)
	must.Equal("Hello world! This is my first post about programming", postResults[0].Text)

	// Test post search for multiple matches
	postResults, err = postStore.Search("the", 0, 10)
	must.NoError(err)
	must.GreaterOrEqual(len(postResults), 1) // Should find at least one post containing "the"

	// Test post search case insensitive
	postResults, err = postStore.Search("COFFEE", 0, 10)
	must.NoError(err)
	must.Len(postResults, 1)
	must.Equal("Coffee and code make the perfect combination", postResults[0].Text)

	// Test search with no results
	searchResults, err = userStore.Search("nonexistent", 0, 10)
	must.NoError(err)
	must.Empty(searchResults)

	postResults, err = postStore.Search("nonexistent", 0, 10)
	must.NoError(err)
	must.Empty(postResults)

	// Test search with pagination
	postResults, err = postStore.Search("", 0, 2) // Empty query should match all
	must.NoError(err)
	must.LessOrEqual(len(postResults), 2)
}