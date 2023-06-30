-- name: CreateFeed :one
INSERT INTO feeds (id,name,url,user_id,created_at,updated_at)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetNextFeedToFetch :many
select * from feeds order by last_fetched_at asc nulls first limit $1;

-- name: MarkFeedAsFetched :one
update feeds 
set last_fetched_at = now(),
updated_at = now()
where id=$1 returning *;