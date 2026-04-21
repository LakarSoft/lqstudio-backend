-- name: GetThemeByID :one
SELECT * FROM themes
WHERE id = $1::varchar LIMIT 1;

-- name: GetActiveThemes :many
SELECT * FROM themes
WHERE is_active = true
  AND ($1::text = '' OR module = $1)
ORDER BY module ASC, name ASC;

-- name: ListAllThemes :many
SELECT * FROM themes
WHERE ($1::text = '' OR module = $1)
ORDER BY module ASC, created_at DESC;

-- name: CreateTheme :one
INSERT INTO themes (
    id, module, name, description, image_url, price, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdateTheme :one
UPDATE themes
SET
    module = $2,
    name = $3,
    description = $4,
    image_url = $5,
    price = $6,
    is_active = $7,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;

-- name: DeleteTheme :exec
DELETE FROM themes WHERE id = $1::varchar;

-- name: ToggleThemeActive :one
UPDATE themes
SET is_active = NOT is_active,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;

-- name: UpdateThemeImageURL :one
UPDATE themes
SET
    image_url = $2,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;
