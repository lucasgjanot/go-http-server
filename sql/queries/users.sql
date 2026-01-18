-- name: CreateUser :one
INSERT INTO
    users (email, hashed_password)
VALUES
    ($1, $2)
RETURNING
    *
;

-- name: DeleteAllUsers :many
DELETE FROM
    users
RETURNING *
;

-- name: GetUserByEmail :one
SELECT * FROM
    users
WHERE
    email = $1
;

-- name: GetUserById :one
SELECT * FROM
    users
WHERE
    id = $1
;

-- name: UpdateUser :one
UPDATE 
    users
SET
    email = $2, hashed_password = $3, updated_at = timezone('utc', now())
WHERE
    id = $1
RETURNING
    *
;

-- name: UpgradeUser :one
UPDATE 
    users
SET
    is_chirpy_red = true, updated_at = timezone('utc', now())
WHERE
    id = $1
RETURNING
    *
;