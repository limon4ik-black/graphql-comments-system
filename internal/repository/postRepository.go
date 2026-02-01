package repository

import (
	"context"

	"github.com/limon4ik-black/graphql-comments-system.git/internal/domain"
)

type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	Get(ctx context.Context, postId string) (*domain.Post, error)
	GetList(ctx context.Context) ([]*domain.Post, error)
	SetFlag(ctx context.Context, postId string, flag bool) error
}
