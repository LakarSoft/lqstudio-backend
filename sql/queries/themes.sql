-- name: GetThemeByID :one
SELECT * FROM themes
WHERE id = $1::varchar LIMIT 1;

-- name: GetActiveThemes :many
SELECT * FROM themes
WHERE is_active = true
ORDER BY name ASC;

-- name: ListAllThemes :many
SELECT * FROM themes
ORDER BY name ASC;

-- name: CreateTheme :one
INSERT INTO themes (
    id, name, description, image_url, price, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateTheme :one
UPDATE themes
SET
    name = $2,
    description = $3,
    image_url = $4,
    price = $5,
    is_active = $6,
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
