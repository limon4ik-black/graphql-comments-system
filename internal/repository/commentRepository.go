package repository

import (
	"context"

	"github.com/limon4ik-black/graphql-comments-system.git/internal/domain"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	GetByPostIDs(ctx context.Context, postId []string) ([]*domain.Comment, error)
	GetByPostID(ctx context.Context, postId string) ([]*domain.Comment, error)
}
