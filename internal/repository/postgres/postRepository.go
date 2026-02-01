package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/limon4ik-black/graphql-comments-system.git/internal/domain"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/repository"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) repository.PostRepository {
	return &PostRepo{db: db}
}

func (r *PostRepo) Create(ctx context.Context, post *domain.Post) error {
	query := `
		INSERT INTO posts (id, title, content, author, comments_allowed)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		post.ID,
		post.Title,
		post.Content,
		post.Author,
		post.Flag,
	)
	return err
}

func (r *PostRepo) Get(ctx context.Context, postId string) (*domain.Post, error) {
	query := `
		SELECT id, title, content, author, comments_allowed
		FROM posts
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, postId)

	post := &domain.Post{}
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Author,
		&post.Flag,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	return post, nil
}

func (r *PostRepo) GetList(ctx context.Context) ([]*domain.Post, error) {
	query := `
		SELECT id, title, content, author, comments_allowed
		FROM posts
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Author,
			&post.Flag,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepo) SetFlag(ctx context.Context, postId string, flag bool) error {
	query := `
		UPDATE posts
		SET comments_allowed = $1
		WHERE id = $2
	`
	res, err := r.db.ExecContext(ctx, query, flag, postId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("post not found")
	}

	return nil
}
