// сюр
package main

import (
	"database/sql"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"github.com/limon4ik-black/graphql-comments-system.git/graph"
	"github.com/limon4ik-black/graphql-comments-system.git/graph/model"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/config"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/logger"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/repository/postgres"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	cfg := config.Load()

	log := logger.New()

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		log.Error("failed to open postgres connection", "error", err)
		os.Exit(1)
	}

	port := cfg.AppPort

	postRepo := postgres.NewPostRepo(db)
	commentRepo := postgres.NewCommentRepo(db)

	postService := service.NewPostService(postRepo, commentRepo, redisClient, log)
	commentService := service.NewCommentService(commentRepo, redisClient, postRepo, log)

	resolver := &graph.Resolver{
		PostService:        postService,
		CommentService:     commentService,
		Redis:              redisClient,
		CommentSubscribers: make(map[string][]chan *model.Comment),
		Mu:                 sync.Mutex{},
		Log:                log,
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Info("connect to http://localhost:" + port + "/ for GraphQL playground")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Error("HTTP server failed", "error", err)
		os.Exit(1)
	}
}
