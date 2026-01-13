-- name: CreateUser :one
INSERT INTO
    users (email)
VALUES
    ($1)
RETURNING
    *
;

-- name: DeleteAllUsers :many
DELETE FROM
    users
RETURNING *
;