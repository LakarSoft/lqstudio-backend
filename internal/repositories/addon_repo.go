package repositories

import (
	"context"

	"lqstudio-backend/internal/database/sqlc"
	"lqstudio-backend/internal/models"
)

// AddonRepository handles addon data access
type AddonRepository struct {
	queries *sqlc.Queries
}

// NewAddonRepository creates a new addon repository
func NewAddonRepository(queries *sqlc.Queries) *AddonRepository {
	return &AddonRepository{queries: queries}
}

// GetByID retrieves an addon by ID
func (r *AddonRepository) GetByID(ctx context.Context, id string) (*models.AddOn, error) {
	result, err := r.queries.GetAddonByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toModel(result), nil
}

// GetByIDs retrieves multiple addons by IDs
func (r *AddonRepository) GetByIDs(ctx context.Context, ids []string) ([]*models.AddOn, error) {
	results, err := r.queries.GetAddonsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return r.toModels(results), nil
}

// GetActive retrieves all active addons
func (r *AddonRepository) GetActive(ctx context.Context) ([]*models.AddOn, error) {
	results, err := r.queries.GetActiveAddons(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(results), nil
}

// ListAll retrieves all addons
func (r *AddonRepository) ListAll(ctx context.Context) ([]*models.AddOn, error) {
	results, err := r.queries.ListAllAddons(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(results), nil
}

// Create creates a new addon
func (r *AddonRepository) Create(ctx context.Context, addon *models.AddOn) error {
	params := sqlc.CreateAddonParams{
		ID:          addon.ID,
		Name:        addon.Name,
		Description: StringPtr(addon.Description),
		Price:       DecimalToNumeric(addon.Price),
		Unit:        StringPtr(addon.Unit),
		IsActive:    BoolPtr(addon.IsActive),
	}

	created, err := r.queries.CreateAddon(ctx, params)
	if err != nil {
		return err
	}

	*addon = *r.toModel(created)
	return nil
}

// Update updates an existing addon
func (r *AddonRepository) Update(ctx context.Context, addon *models.AddOn) error {
	params := sqlc.UpdateAddonParams{
		Column1:     addon.ID,
		Name:        addon.Name,
		Description: StringPtr(addon.Description),
		Price:       DecimalToNumeric(addon.Price),
		Unit:        StringPtr(addon.Unit),
		IsActive:    BoolPtr(addon.IsActive),
	}

	updated, err := r.queries.UpdateAddon(ctx, params)
	if err != nil {
		return err
	}

	*addon = *r.toModel(updated)
	return nil
}

// Delete deletes an addon
func (r *AddonRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteAddon(ctx, id)
}

// ToggleActive toggles addon active status
func (r *AddonRepository) ToggleActive(ctx context.Context, id string) (*models.AddOn, error) {
	result, err := r.queries.ToggleAddonActive(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toModel(result), nil
}

// Helper methods
func (r *AddonRepository) toModel(row sqlc.Addon) *models.AddOn {
	return &models.AddOn{
		ID:          row.ID,
		Name:        row.Name,
		Description: StringVal(row.Description),
		Price:       NumericToDecimal(row.Price),
		Unit:        StringVal(row.Unit),
		IsActive:    BoolVal(row.IsActive),
		CreatedAt:   TimestamptzToTime(row.CreatedAt),
		UpdatedAt:   TimestamptzToTime(row.UpdatedAt),
	}
}

func (r *AddonRepository) toModels(rows []sqlc.Addon) []*models.AddOn {
	addons := make([]*models.AddOn, len(rows))
	for i, row := range rows {
		addons[i] = r.toModel(row)
	}
	return addons
}
