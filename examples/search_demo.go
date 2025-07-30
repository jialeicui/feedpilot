package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jialeicui/feedpilot/pkg/meta"
	"github.com/jialeicui/feedpilot/pkg/store"
	"github.com/jialeicui/feedpilot/pkg/store/kv"
)

func main() {
	// Create a temporary directory for the demo
	dir, err := os.MkdirTemp("", "feedpilot-search-demo")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create a Badger store
	badgerStore, err := kv.NewBadger(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Create user and post stores with separate buckets
	userStore := store.NewUserStore(badgerStore.WithBucket("users"))
	postStore := store.NewPostStore(badgerStore.WithBucket("posts"))

	// Insert some sample users
	users := []*meta.User{
		{ID: "1", Username: "alice", DisplayName: "Alice Smith", Bio: "Software developer from NYC, loves Go programming"},
		{ID: "2", Username: "bob", DisplayName: "Bob Jones", Bio: "UI/UX designer and digital artist"},
		{ID: "3", Username: "charlie", DisplayName: "Charlie Brown", Bio: "Product manager and startup entrepreneur"},
		{ID: "4", Username: "diana", DisplayName: "Diana Prince", Bio: "Data scientist working on machine learning"},
	}

	fmt.Println("🚀 FeedPilot Search Demo")
	fmt.Println("========================")
	fmt.Println()

	fmt.Println("📝 Inserting sample users...")
	for _, user := range users {
		if err := userStore.Put(user.ID, user); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  ✓ %s (@%s): %s\n", user.DisplayName, user.Username, user.Bio)
	}
	fmt.Println()

	// Insert some sample posts
	posts := []*meta.Post{
		{ID: "1", Text: "Hello world! This is my first post about programming with Go"},
		{ID: "2", Text: "Beautiful sunset today. Nature is absolutely amazing!"},
		{ID: "3", Text: "Working on a new feature for our mobile app. Excited to ship it next week!"},
		{ID: "4", Text: "Coffee and code make the perfect combination for productive mornings"},
		{ID: "5", Text: "Just finished reading an excellent book about machine learning algorithms"},
	}

	fmt.Println("📝 Inserting sample posts...")
	for _, post := range posts {
		if err := postStore.Put(post.ID, post); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  ✓ Post %s: %s\n", post.ID, truncateText(post.Text, 50))
	}
	fmt.Println()

	// Demonstrate user search
	fmt.Println("🔍 User Search Examples")
	fmt.Println("------------------------")

	// Search by username
	fmt.Printf("Search for 'alice':\n")
	searchResults, err := userStore.Search("alice", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	printUserResults(searchResults)

	// Search by profession
	fmt.Printf("Search for 'developer':\n")
	searchResults, err = userStore.Search("developer", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	printUserResults(searchResults)

	// Search by location
	fmt.Printf("Search for 'NYC':\n")
	searchResults, err = userStore.Search("NYC", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	printUserResults(searchResults)

	// Case insensitive search
	fmt.Printf("Search for 'DESIGNER' (case insensitive):\n")
	searchResults, err = userStore.Search("DESIGNER", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	printUserResults(searchResults)

	// Demonstrate post search
	fmt.Println("🔍 Post Search Examples")
	fmt.Println("------------------------")

	// Search for programming-related posts
	fmt.Printf("Search for 'programming':\n")
	postResults, err := postStore.Search("programming", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	printPostResults(postResults)

	// Search for app-related posts
	fmt.Printf("Search for 'app':\n")
	postResults, err = postStore.Search("app", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	printPostResults(postResults)

	// Search for nature-related posts
	fmt.Printf("Search for 'nature':\n")
	postResults, err = postStore.Search("nature", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	printPostResults(postResults)

	// Search with no results
	fmt.Printf("Search for 'nonexistent':\n")
	postResults, err = postStore.Search("nonexistent", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	printPostResults(postResults)

	fmt.Println("✅ Demo completed successfully!")
}

func printUserResults(users []*meta.User) {
	if len(users) == 0 {
		fmt.Println("  No users found.")
	} else {
		for _, user := range users {
			fmt.Printf("  ✓ %s (@%s): %s\n", user.DisplayName, user.Username, user.Bio)
		}
	}
	fmt.Println()
}

func printPostResults(posts []*meta.Post) {
	if len(posts) == 0 {
		fmt.Println("  No posts found.")
	} else {
		for _, post := range posts {
			fmt.Printf("  ✓ Post %s: %s\n", post.ID, truncateText(post.Text, 60))
		}
	}
	fmt.Println()
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}