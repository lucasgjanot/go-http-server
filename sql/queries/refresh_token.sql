-- name: CreateRefreshToken :one

INSERT INTO
    refresh_tokens (token, user_id, expires_at, created_at, updated_at)
VALUES
    ($1,$2,$3,$4,$5)
RETURNING 
    *
;

-- name: GetRefreshToken :one
SELECT 
    *
FROM
    refresh_tokens
WHERE
    token = $1
;

-- name: RevokeToken :one
UPDATE 
    refresh_tokens
SET 
    revoked_at = timezone('utc', now())
WHERE
    token = $1
RETURNING
    *
;
