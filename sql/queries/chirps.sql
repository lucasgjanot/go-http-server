-- name: CreateChirp :one
INSERT INTO
    chirps (body, user_id)
VALUES
    ($1, $2)
RETURNING
    *
;

-- name: GetChirps :many
SELECT 
    *
FROM 
    chirps
ORDER BY
	CASE WHEN @is_desc::boolean THEN created_at END DESC,
	CASE WHEN NOT @is_desc::boolean THEN created_at END ASC;

-- name: GetChirpsByUserId :many
SELECT 
    *
FROM 
    chirps
WHERE   
    user_id = $1
ORDER BY
	CASE WHEN @is_desc::boolean THEN created_at END DESC,
	CASE WHEN NOT @is_desc::boolean THEN created_at END ASC;


-- name: GetChirp :one
SELECT 
    *
FROM 
    chirps
WHERE
    id = $1
;

-- name: DeleteChirp :one
DELETE FROM chirps
WHERE id = $1
RETURNING *
;