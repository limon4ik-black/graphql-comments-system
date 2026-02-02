package service

import (
	"context"
	"errors"
	"io"
	"testing"

	"log/slog"

	"github.com/limon4ik-black/graphql-comments-system.git/internal/domain"
)

type mockPostRepo struct {
	createFunc  func(ctx context.Context, post *domain.Post) error
	getFunc     func(ctx context.Context, postId string) (*domain.Post, error)
	getListFunc func(ctx context.Context) ([]*domain.Post, error)
	setFlagFunc func(ctx context.Context, postId string, flag bool) error
}

func (m *mockPostRepo) Create(ctx context.Context, post *domain.Post) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, post)
	}
	return nil
}
func (m *mockPostRepo) Get(ctx context.Context, postId string) (*domain.Post, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, postId)
	}
	return nil, nil
}
func (m *mockPostRepo) GetList(ctx context.Context) ([]*domain.Post, error) {
	if m.getListFunc != nil {
		return m.getListFunc(ctx)
	}
	return nil, nil
}
func (m *mockPostRepo) SetFlag(ctx context.Context, postId string, flag bool) error {
	if m.setFlagFunc != nil {
		return m.setFlagFunc(ctx, postId, flag)
	}
	return nil
}

type mockCommentRepo struct {
	createFunc       func(ctx context.Context, comment *domain.Comment) error
	getByPostIDFunc  func(ctx context.Context, postID string) ([]*domain.Comment, error)
	getByPostIDsFunc func(ctx context.Context, postIDs []string) ([]*domain.Comment, error)
}

func (m *mockCommentRepo) Create(ctx context.Context, comment *domain.Comment) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, comment)
	}
	return nil
}
func (m *mockCommentRepo) GetByPostID(ctx context.Context, postID string) ([]*domain.Comment, error) {
	if m.getByPostIDFunc != nil {
		return m.getByPostIDFunc(ctx, postID)
	}
	return nil, nil
}

func (m *mockCommentRepo) GetByPostIDs(ctx context.Context, postIDs []string) ([]*domain.Comment, error) {
	if m.getByPostIDsFunc != nil {
		return m.getByPostIDsFunc(ctx, postIDs)
	}
	return nil, nil
}

func TestPostService_Create(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("success", func(t *testing.T) {
		mockRepo := &mockPostRepo{
			createFunc: func(ctx context.Context, post *domain.Post) error {
				return nil
			},
		}
		s := NewPostService(mockRepo, nil, nil, logger)
		p := &domain.Post{Title: "test", Content: "test", Author: "test"}
		if err := s.Create(context.Background(), p); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if p.ID == "" {
			t.Fatalf("expected ID to be set")
		}
	})

	t.Run("nil post", func(t *testing.T) {
		mockRepo := &mockPostRepo{}
		s := NewPostService(mockRepo, nil, nil, logger)
		err := s.Create(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error for nil post")
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		mockRepo := &mockPostRepo{}
		s := NewPostService(mockRepo, nil, nil, logger)
		err := s.Create(context.Background(), &domain.Post{Title: "", Content: "", Author: ""})
		if err != nil {
			t.Fatal("Create should return nil for empty fields in this implementation")
		}
	})
}

func TestPostService_Get(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()
	post := &domain.Post{ID: "1", Title: "test", Content: "test", Author: "test"}
	comments := []*domain.Comment{
		{ID: "c1", PostID: "1", Text: "hello"},
	}

	mockRepo := &mockPostRepo{
		getFunc: func(ctx context.Context, postId string) (*domain.Post, error) {
			return post, nil
		},
	}
	mockComments := &mockCommentRepo{
		getByPostIDFunc: func(ctx context.Context, postID string) ([]*domain.Comment, error) {
			return comments, nil
		},
	}

	s := NewPostService(mockRepo, mockComments, nil, logger)

	t.Run("success", func(t *testing.T) {
		got, err := s.Get(ctx, "1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != "1" {
			t.Errorf("expected ID '1', got %s", got.ID)
		}
		if len(got.Comments) != 1 {
			t.Errorf("expected 1 comment, got %d", len(got.Comments))
		}
	})

	t.Run("empty postId", func(t *testing.T) {
		_, err := s.Get(ctx, "")
		if err == nil {
			t.Fatal("expected error for empty postId")
		}
	})
}

func TestPostService_GetList(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()
	posts := []*domain.Post{
		{ID: "1", Title: "test", Content: "test", Author: "test"},
		{ID: "2", Title: "test2", Content: "test2", Author: "test2"},
	}
	comments := []*domain.Comment{
		{ID: "c1", PostID: "1", Text: "test"},
		{ID: "c2", PostID: "2", Text: "test2"},
	}

	mockRepo := &mockPostRepo{
		getListFunc: func(ctx context.Context) ([]*domain.Post, error) {
			return posts, nil
		},
	}
	mockComments := &mockCommentRepo{
		getByPostIDsFunc: func(ctx context.Context, postIDs []string) ([]*domain.Comment, error) {
			return comments, nil
		},
	}

	s := NewPostService(mockRepo, mockComments, nil, logger)

	got, err := s.GetList(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 posts, got %d", len(got))
	}
	if len(got[0].Comments) != 1 {
		t.Errorf("expected 1 comment for first post, got %d", len(got[0].Comments))
	}
}

func TestPostService_SetFlag(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	ctx := context.Background()
	mockRepo := &mockPostRepo{
		setFlagFunc: func(ctx context.Context, postId string, flag bool) error {
			if postId == "error" {
				return errors.New("repo error")
			}
			return nil
		},
	}

	s := NewPostService(mockRepo, nil, nil, logger)

	t.Run("success", func(t *testing.T) {
		err := s.SetFlag(ctx, "1", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty postId", func(t *testing.T) {
		err := s.SetFlag(ctx, "", true)
		if err == nil {
			t.Fatal("expected error for empty postId")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		err := s.SetFlag(ctx, "error", true)
		if err == nil {
			t.Fatal("expected error from repo")
		}
	})
}
