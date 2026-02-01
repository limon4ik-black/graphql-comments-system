package graph

import (
	"database/sql"
	"log/slog"
	"sync"

	"github.com/limon4ik-black/graphql-comments-system.git/graph/model"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/service"
	"github.com/redis/go-redis/v9"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	PostService        *service.PostService
	CommentService     *service.CommentService
	DB                 *sql.DB
	Redis              *redis.Client
	CommentSubscribers map[string][]chan *model.Comment
	Mu                 sync.Mutex
	Log                *slog.Logger
}
