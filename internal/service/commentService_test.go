package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/limon4ik-black/graphql-comments-system.git/internal/domain"
)

func TestCommentService_Create(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("success", func(t *testing.T) {
		mockCommentRepo := &mockCommentRepo{
			createFunc: func(ctx context.Context, comment *domain.Comment) error {
				return nil
			},
		}

		mockPostRepo := &mockPostRepo{
			getFunc: func(ctx context.Context, postID string) (*domain.Post, error) {
				return &domain.Post{ID: postID, Flag: true}, nil
			},
		}

		s := NewCommentService(mockCommentRepo, nil, mockPostRepo, logger)

		comment := &domain.Comment{
			PostID: "post-1",
			Text:   "hello",
			Author: "user",
		}

		err := s.Create(context.Background(), comment)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if comment.ID == "" {
			t.Fatal("expected comment ID to be set")
		}
	})

	t.Run("nil comment", func(t *testing.T) {
		s := NewCommentService(nil, nil, nil, logger)
		err := s.Create(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error for nil comment")
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		s := NewCommentService(nil, nil, nil, logger)
		err := s.Create(context.Background(), &domain.Comment{})
		if err == nil {
			t.Fatal("expected error for missing fields")
		}
	})

	t.Run("comments forbidden", func(t *testing.T) {
		mockCommentRepo := &mockCommentRepo{
			createFunc: func(ctx context.Context, comment *domain.Comment) error {
				return nil
			},
		}

		mockPostRepo := &mockPostRepo{
			getFunc: func(ctx context.Context, postID string) (*domain.Post, error) {
				return &domain.Post{ID: postID, Flag: false}, nil
			},
		}

		s := NewCommentService(mockCommentRepo, nil, mockPostRepo, logger)

		comment := &domain.Comment{
			PostID: "post-1",
			Text:   "hello",
			Author: "user",
		}

		err := s.Create(context.Background(), comment)
		if err == nil {
			t.Fatal("expected error when comments are forbidden")
		}
	})

	t.Run("post repo error", func(t *testing.T) {
		mockCommentRepo := &mockCommentRepo{
			createFunc: func(ctx context.Context, comment *domain.Comment) error {
				return nil
			},
		}

		mockPostRepo := &mockPostRepo{
			getFunc: func(ctx context.Context, postID string) (*domain.Post, error) {
				return nil, errors.New("db error")
			},
		}

		s := NewCommentService(mockCommentRepo, nil, mockPostRepo, logger)

		comment := &domain.Comment{
			PostID: "post-1",
			Text:   "hello",
			Author: "user",
		}

		err := s.Create(context.Background(), comment)
		if err == nil {
			t.Fatal("expected error from post repo")
		}
	})
}

func TestCommentService_GetByPostID(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockCommentRepo{
			getByPostIDFunc: func(ctx context.Context, postID string) ([]*domain.Comment, error) {
				return []*domain.Comment{
					{ID: "1", PostID: postID, Text: "a"},
					{ID: "2", PostID: postID, Text: "b"},
				}, nil
			},
		}

		s := NewCommentService(mockRepo, nil, nil, log)

		comments, err := s.GetByPostID(context.Background(), "post-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(comments) != 2 {
			t.Fatalf("expected 2 comments, got %d", len(comments))
		}
	})

	t.Run("empty postID", func(t *testing.T) {
		s := NewCommentService(nil, nil, nil, log)

		_, err := s.GetByPostID(context.Background(), "")
		if err == nil {
			t.Fatal("expected error for empty postID")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := &mockCommentRepo{
			getByPostIDFunc: func(ctx context.Context, postID string) ([]*domain.Comment, error) {
				return nil, errors.New("db error")
			},
		}

		s := NewCommentService(mockRepo, nil, nil, log)

		_, err := s.GetByPostID(context.Background(), "post-1")
		if err == nil {
			t.Fatal("expected error from repo")
		}
	})
}
