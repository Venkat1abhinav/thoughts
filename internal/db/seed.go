package db

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/owned_dragon/thoughts/internal/store"
)

const (
	userCount          = 100
	minPostsPerUser    = 3
	maxPostsPerUser    = 100
	minCommentsPerPost = 3
	maxCommentsPerPost = 100
)

func Seed(storage store.Store) error {
	ctx := context.Background()

	// Users
	users := generateUsers(userCount)

	if err := storage.Users.CreateMany(ctx, users); err != nil {
		return fmt.Errorf("create users: %w", err)
	}

	fmt.Printf("created %d users\n", len(users))

	// Generate posts
	var posts []*store.Post

	for _, user := range users {
		postCount := gofakeit.Number(
			minPostsPerUser,
			maxPostsPerUser,
		)

		for range postCount {
			posts = append(posts, generatePost(user.ID))
		}
	}

	// Posts
	if err := storage.Posts.CreateMany(ctx, posts); err != nil {
		return fmt.Errorf("create posts: %w", err)
	}

	fmt.Printf("created %d posts\n", len(posts))

	// Generate comments
	var comments []*store.Comment

	for _, post := range posts {
		commentCount := gofakeit.Number(
			minCommentsPerPost,
			maxCommentsPerPost,
		)

		comments = append(
			comments,
			generateComments(users, post.ID, commentCount)...,
		)
	}

	// Comments
	if err := storage.Comments.CreateMany(ctx, comments); err != nil {
		return fmt.Errorf("create comments: %w", err)
	}

	fmt.Printf("created %d comments\n", len(comments))
	fmt.Println("seeding completed")

	return nil
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := range users {
		firstName := gofakeit.FirstName()
		lastName := gofakeit.LastName()

		username := strings.ToLower(
			firstName + "." +
				lastName +
				strconv.Itoa(1000+i),
		)

		users[i] = &store.User{
			Username:  username,
			FirstName: firstName,
			LastName:  lastName,
			Email:     username + "@example.com",
			Password:  []byte(gofakeit.Password(true, false, true, false, false, 12)),
		}
	}

	return users
}

func generatePost(userID int64) *store.Post {
	return &store.Post{
		UserID:  userID,
		Title:   gofakeit.Sentence(5),
		Content: gofakeit.Paragraph(2, 4, 10, " "),
		Tags:    gofakeit.Product().Categories,
	}
}

func generateComments(
	users []*store.User,
	postID int64,
	num int,
) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := range comments {
		user := users[gofakeit.Number(0, len(users)-1)]

		comments[i] = &store.Comment{
			PostID: postID,
			UserID: user.ID,

			Content: gofakeit.Paragraph(1, 3, 8, " "),
		}
	}

	return comments
}
