// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: feeds.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id,name,url,user_id,created_at,updated_at)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING id, name, url, user_id, created_at, updated_at, last_fetched_at
`

type CreateFeedParams struct {
	ID        uuid.UUID
	Name      string
	Url       string
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.Name,
		arg.Url,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT id, name, url, user_id, created_at, updated_at, last_fetched_at FROM feeds
`

func (q *Queries) GetFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.LastFetchedAt,
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

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :many
select id, name, url, user_id, created_at, updated_at, last_fetched_at from feeds order by last_fetched_at asc nulls first limit $1
`

func (q *Queries) GetNextFeedToFetch(ctx context.Context, limit int32) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getNextFeedToFetch, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.LastFetchedAt,
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

const markFeedAsFetched = `-- name: MarkFeedAsFetched :one
update feeds 
set last_fetched_at = now(),
updated_at = now()
where id=$1 returning id, name, url, user_id, created_at, updated_at, last_fetched_at
`

func (q *Queries) MarkFeedAsFetched(ctx context.Context, id uuid.UUID) (Feed, error) {
	row := q.db.QueryRowContext(ctx, markFeedAsFetched, id)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}
