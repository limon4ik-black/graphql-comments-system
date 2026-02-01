package postgres

import (
	"context"
	"database/sql"

	"github.com/limon4ik-black/graphql-comments-system.git/internal/domain"
	"github.com/limon4ik-black/graphql-comments-system.git/internal/repository"
)

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) repository.CommentRepository {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (id, post_id, parent_id, author, text)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		comment.ID,
		comment.PostID,
		comment.ParentID,
		comment.Author,
		comment.Text,
	)
	return err
}

func (r *CommentRepo) GetByPostID(ctx context.Context, postID string) ([]*domain.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, author, text
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		c := &domain.Comment{}
		if err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.ParentID,
			&c.Author,
			&c.Text,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *CommentRepo) GetByPostIDs(ctx context.Context, postIDs []string) ([]*domain.Comment, error) {
	if len(postIDs) == 0 {
		return nil, nil
	}

	query := `
		SELECT id, post_id, parent_id, author, text
		FROM comments
		WHERE post_id = ANY($1)
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, postIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		c := &domain.Comment{}
		if err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.ParentID,
			&c.Author,
			&c.Text,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
