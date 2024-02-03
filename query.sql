-- name: GetEntry :one
SELECT
    sqlc.embed(entry)
FROM
    entry
WHERE
    id = ?
LIMIT
    1;

-- name: GetEntries :many
SELECT
    sqlc.embed(entry)
FROM
    entry
ORDER BY
    id;

-- name: CreateEntry :one
INSERT INTO
    entry (name, origin, desc)
VALUES
    (?, ?, ?) RETURNING *;