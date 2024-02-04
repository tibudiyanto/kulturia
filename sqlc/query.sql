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
    sqlc.embed(entry),
    sqlc.embed(a)
FROM
    entry
    JOIN asset a on a.entry_id = entry.id
ORDER BY
    entry.id;

-- name: CreateEntry :one
INSERT INTO
    entry (name, origin, desc)
VALUES
    (?, ?, ?) RETURNING *;

-- name: CreateAsset :one
INSERT INTO
    asset (entry_id, "location")
VALUES
    (?, ?) RETURNING *;