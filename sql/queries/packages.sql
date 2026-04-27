-- name: GetPackageByID :one
SELECT * FROM packages
WHERE id = $1::varchar LIMIT 1;

-- name: GetActivePackages :many
SELECT * FROM packages
WHERE is_active = true
  AND ($1::text = '' OR module = $1)
ORDER BY module ASC, price ASC;

-- name: ListAllPackages :many
SELECT * FROM packages
WHERE ($1::text = '' OR module = $1)
ORDER BY module ASC, created_at DESC;

-- name: CreatePackage :one
INSERT INTO packages (
    id, module, name, description, duration_minutes, price, discount, offers, image_url, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdatePackage :one
UPDATE packages
SET
    module = $2,
    name = $3,
    description = $4,
    duration_minutes = $5,
    price = $6,
    discount = $7,
    offers = $8,
    image_url = $9,
    is_active = $10,
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

-- name: UpdatePackageImageURL :one
UPDATE packages
SET
    image_url = $2,
    updated_at = NOW()
WHERE id = $1::varchar
RETURNING *;
