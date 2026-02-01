package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/domain"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/repository"
	"github.com/redis/go-redis/v9"
)

type PostService struct {
	repo        repository.PostRepository
	commentRepo repository.CommentRepository
	redis       *redis.Client
	log         *slog.Logger
}

func NewPostService(repo repository.PostRepository, commentRepo repository.CommentRepository, redis *redis.Client, log *slog.Logger) *PostService {
	return &PostService{repo: repo, commentRepo: commentRepo, redis: redis, log: log}
}

func (p *PostService) Create(ctx context.Context, post *domain.Post) error {
	if post == nil {
		err := errors.New("post is nil")
		p.log.Error("Create post failed", "error", err)
		return err
	}

	if post.Title == "" || post.Content == "" || post.Author == "" {
		err := errors.New("title, content and author are required")
		p.log.Error("Create post failed", "error", err)
		return nil
	}

	post.Flag = true

	if post.ID == "" {
		post.ID = uuid.NewString()
	}

	if err := p.repo.Create(ctx, post); err != nil {
		p.log.Error("failed create post in repo", "error", err)
		return err
	}

	if p.redis != nil {
		if err := p.redis.Del(ctx, "posts:list").Err(); err != nil {
			p.log.Error("failed del from redis", "error", err)
		}
	}

	return nil
}

func (p *PostService) Get(ctx context.Context, postId string) (*domain.Post, error) {
	if postId == "" {
		err := errors.New("postId is required")
		p.log.Error("Get post failed", "error", err)
		return nil, err
	}

	cacheKey := "post:" + postId

	if p.redis != nil {
		cached, err := p.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var post domain.Post
			if err := json.Unmarshal([]byte(cached), &post); err == nil {
				return &post, nil
			}
		}
	}

	post, err := p.repo.Get(ctx, postId)
	if err != nil {
		p.log.Error("failed Get post repo", "error", err)
		return nil, err
	}

	comments, err := p.commentRepo.GetByPostID(ctx, post.ID)
	if err != nil {
		p.log.Error("failed get cooments of post repo", "error", err)
		return nil, err
	}
	post.Comments = buildCommentTree(comments)

	if p.redis != nil {
		bytes, _ := json.Marshal(post)
		if err = p.redis.Set(ctx, cacheKey, bytes, time.Minute*5).Err(); err != nil {
			p.log.Error("failed set to redis", "error", err)
		}
	}

	return post, nil
}

func (p *PostService) GetList(ctx context.Context) ([]*domain.Post, error) {
	cacheKey := "posts:list"

	if p.redis != nil {
		cached, err := p.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var posts []*domain.Post
			if err := json.Unmarshal([]byte(cached), &posts); err == nil {
				return posts, nil
			}
		}
	}

	posts, err := p.repo.GetList(ctx)
	if err != nil {
		p.log.Error("failed get list repo", "error", err)
		return nil, err
	}

	if len(posts) == 0 {
		return posts, nil
	}

	postIDs := make([]string, 0, len(posts))
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}

	comments, err := p.commentRepo.GetByPostIDs(ctx, postIDs)
	if err != nil {
		p.log.Error("failed get comments for posts repo", "error", err)
		return nil, err
	}

	commentsByPost := make(map[string][]*domain.Comment)
	for _, c := range comments {
		commentsByPost[c.PostID] = append(commentsByPost[c.PostID], c)
	}

	for _, post := range posts {
		post.Comments = buildCommentTree(commentsByPost[post.ID])
	}

	if p.redis != nil {
		bytes, _ := json.Marshal(posts)
		if err := p.redis.Set(ctx, cacheKey, bytes, time.Minute*5).Err(); err != nil {
			p.log.Error("failed set туту епт redis", "error", err)
		}
	}

	return posts, nil
}

func (p *PostService) SetFlag(ctx context.Context, postId string, flag bool) error {
	if postId == "" {
		err := errors.New("postId is required")
		p.log.Error("failed switch flag", "error", err)
		return err
	}

	if err := p.repo.SetFlag(ctx, postId, flag); err != nil {
		p.log.Error("failed update flag repo", "error", err)
		return err
	}

	if p.redis != nil {
		if err := p.redis.Del(ctx, "post:"+postId).Err(); err != nil {
			p.log.Warn("failed del in redis func:SetFlag", "error", err)
		}
		if err := p.redis.Del(ctx, "posts:list").Err(); err != nil {
			p.log.Warn("failed del in redis func:SetFlag", "error", err)
		}
	}

	return nil
}

func buildCommentTree(comments []*domain.Comment) []*domain.Comment {
	byParent := make(map[string][]*domain.Comment)

	for _, c := range comments {
		parentID := ""
		if c.ParentID != nil {
			parentID = *c.ParentID
		}
		byParent[parentID] = append(byParent[parentID], c)
	}

	var attach func(parentID string) []*domain.Comment
	attach = func(parentID string) []*domain.Comment {
		children := byParent[parentID]
		for _, c := range children {
			c.Children = attach(c.ID)
		}
		return children
	}

	return attach("")
}
