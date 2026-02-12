-- name: GetPackageByID :one
SELECT * FROM packages
WHERE id = $1::varchar LIMIT 1;

-- name: GetActivePackages :many
SELECT * FROM packages
WHERE is_active = true
ORDER BY created_at DESC;

-- name: ListAllPackages :many
SELECT * FROM packages
ORDER BY created_at DESC;

-- name: CreatePackage :one
INSERT INTO packages (
    id, name, description, duration_minutes, price, discount, offers, image_url, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: UpdatePackage :one
UPDATE packages
SET
    name = $2,
    description = $3,
    duration_minutes = $4,
    price = $5,
    discount = $6,
    offers = $7,
    image_url = $8,
    is_active = $9,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;

-- name: DeletePackage :exec
DELETE FROM packages WHERE id = $1::varchar;

-- name: TogglePackageActive :one
UPDATE packages
SET is_active = NOT is_active,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;
