package repositories

import (
	"context"
	"fmt"

	"lqstudio-backend/internal/database/sqlc"
	"lqstudio-backend/internal/models"
)

// PackageRepository handles package data access
type PackageRepository struct {
	queries *sqlc.Queries
}

// NewPackageRepository creates a new package repository
func NewPackageRepository(queries *sqlc.Queries) *PackageRepository {
	return &PackageRepository{queries: queries}
}

// GetByID retrieves a package by ID
func (r *PackageRepository) GetByID(ctx context.Context, id string) (*models.Package, error) {
	pkg, err := r.queries.GetPackageByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toModel(pkg)
}

// GetActive retrieves all active packages
func (r *PackageRepository) GetActive(ctx context.Context) ([]*models.Package, error) {
	rows, err := r.queries.GetActivePackages(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(rows)
}

// ListAll retrieves all packages
func (r *PackageRepository) ListAll(ctx context.Context) ([]*models.Package, error) {
	rows, err := r.queries.ListAllPackages(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(rows)
}

// Create creates a new package
func (r *PackageRepository) Create(ctx context.Context, pkg *models.Package) error {
	// Convert offers to JSONB
	offersJSON, err := StringsToJSONB(pkg.Offers)
	if err != nil {
		return fmt.Errorf("failed to marshal offers: %w", err)
	}

	params := sqlc.CreatePackageParams{
		ID:              pkg.ID,
		Name:            pkg.Name,
		Description:     StringPtr(pkg.Description),
		DurationMinutes: pkg.DurationMinutes,
		Price:           DecimalToNumeric(pkg.Price),
		Discount:        DecimalToNumeric(pkg.Discount),
		Offers:          offersJSON,
		ImageUrl:        StringPtr(pkg.ImageURL),
		IsActive:        BoolPtr(pkg.IsActive),
	}

	created, err := r.queries.CreatePackage(ctx, params)
	if err != nil {
		return err
	}

	// Update pkg with returned values
	result, err := r.toModel(created)
	if err != nil {
		return err
	}
	*pkg = *result
	return nil
}

// Update updates an existing package
func (r *PackageRepository) Update(ctx context.Context, pkg *models.Package) error {
	// Convert offers to JSONB
	offersJSON, err := StringsToJSONB(pkg.Offers)
	if err != nil {
		return fmt.Errorf("failed to marshal offers: %w", err)
	}

	params := sqlc.UpdatePackageParams{
		Column1:         pkg.ID,
		Name:            pkg.Name,
		Description:     StringPtr(pkg.Description),
		DurationMinutes: pkg.DurationMinutes,
		Price:           DecimalToNumeric(pkg.Price),
		Discount:        DecimalToNumeric(pkg.Discount),
		Offers:          offersJSON,
		ImageUrl:        StringPtr(pkg.ImageURL),
		IsActive:        BoolPtr(pkg.IsActive),
	}

	updated, err := r.queries.UpdatePackage(ctx, params)
	if err != nil {
		return err
	}

	// Update pkg with returned values
	result, err := r.toModel(updated)
	if err != nil {
		return err
	}
	*pkg = *result
	return nil
}

// Delete deletes a package
func (r *PackageRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeletePackage(ctx, id)
}

// ToggleActive toggles package active status
func (r *PackageRepository) ToggleActive(ctx context.Context, id string) (*models.Package, error) {
	pkg, err := r.queries.TogglePackageActive(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toModel(pkg)
}

// Helper methods
func (r *PackageRepository) toModel(row sqlc.Package) (*models.Package, error) {
	// Convert JSONB offers to []string
	offers, err := JSONBToStrings(row.Offers)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal offers: %w", err)
	}

	return &models.Package{
		ID:              row.ID,
		Name:            row.Name,
		Description:     StringVal(row.Description),
		DurationMinutes: row.DurationMinutes,
		Offers:          offers,
		Price:           NumericToDecimal(row.Price),
		Discount:        NumericToDecimal(row.Discount),
		ImageURL:        StringVal(row.ImageUrl),
		IsActive:        BoolVal(row.IsActive),
		CreatedAt:       TimestamptzToTime(row.CreatedAt),
		UpdatedAt:       TimestamptzToTime(row.UpdatedAt),
	}, nil
}

func (r *PackageRepository) toModels(rows []sqlc.Package) ([]*models.Package, error) {
	packages := make([]*models.Package, len(rows))
	for i, row := range rows {
		pkg, err := r.toModel(row)
		if err != nil {
			return nil, err
		}
		packages[i] = pkg
	}
	return packages, nil
}
