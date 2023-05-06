// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: feature.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const getFeature = `-- name: GetFeature :one
SELECT id, title, slug, parent_id, description, created_by, updated_by, created_at, updated_at FROM features
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetFeature(ctx context.Context, id uuid.UUID) (Feature, error) {
	row := q.db.QueryRowContext(ctx, getFeature, id)
	var i Feature
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Slug,
		&i.ParentID,
		&i.Description,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listFeatures = `-- name: ListFeatures :many
SELECT id, title, slug, parent_id, description, created_by, updated_by, created_at, updated_at FROM features 
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListFeaturesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListFeatures(ctx context.Context, arg ListFeaturesParams) ([]Feature, error) {
	rows, err := q.db.QueryContext(ctx, listFeatures, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Feature{}
	for rows.Next() {
		var i Feature
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Slug,
			&i.ParentID,
			&i.Description,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
