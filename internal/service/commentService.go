package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/domain"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/repository"
	"github.com/redis/go-redis/v9"
)

type CommentService struct {
	repo     repository.CommentRepository
	redis    *redis.Client
	postRepo repository.PostRepository
	log      *slog.Logger
}

func NewCommentService(repo repository.CommentRepository, redis *redis.Client, postRepo repository.PostRepository, log *slog.Logger) *CommentService {
	return &CommentService{repo: repo, redis: redis, postRepo: postRepo, log: log}
}

func (s *CommentService) Create(ctx context.Context, comment *domain.Comment) error {
	if comment == nil {
		err := errors.New("comment is nil")
		s.log.Error("failed create comment", "error", err)
		return err
	}

	if comment.PostID == "" || comment.Author == "" || comment.Text == "" {
		err := errors.New("postID, author and text are required")
		s.log.Error("failed create comment", "error", err)
		return err
	}

	post, err := s.postRepo.Get(ctx, comment.PostID)
	if err != nil {
		return err
	}
	if !post.Flag {
		err = errors.New("comments are disabled for this post")
		s.log.Warn("flag is false", "error", err)
		return err
	}

	if comment.ID == "" {
		comment.ID = uuid.NewString()
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		s.log.Error("failed create comment repo", "error", err)
		return err
	}

	if s.redis != nil {
		if err := s.redis.Del(ctx, "post:"+comment.PostID).Err(); err != nil {
			s.log.Warn("failed del from redis", "error", err)
		}
		if err := s.redis.Del(ctx, "posts:list").Err(); err != nil {
			s.log.Warn("failed del from redis", "error", err)
		}
	}

	return nil
}

func (s *CommentService) GetByPostID(ctx context.Context, postID string) ([]*domain.Comment, error) {
	if postID == "" {
		err := errors.New("postID is required")
		s.log.Error("failed get post", "error", err)
		return nil, err
	}
	return s.repo.GetByPostID(ctx, postID)
}
