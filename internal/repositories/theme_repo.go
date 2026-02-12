package repositories

import (
	"context"

	"lqstudio-backend/internal/database/sqlc"
	"lqstudio-backend/internal/models"
)

// ThemeRepository handles theme data access
type ThemeRepository struct {
	queries *sqlc.Queries
}

// NewThemeRepository creates a new theme repository
func NewThemeRepository(queries *sqlc.Queries) *ThemeRepository {
	return &ThemeRepository{queries: queries}
}

// GetByID retrieves a theme by ID
func (r *ThemeRepository) GetByID(ctx context.Context, id string) (*models.Theme, error) {
	result, err := r.queries.GetThemeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toModel(result), nil
}

// GetActive retrieves all active themes
func (r *ThemeRepository) GetActive(ctx context.Context) ([]*models.Theme, error) {
	results, err := r.queries.GetActiveThemes(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(results), nil
}

// ListAll retrieves all themes
func (r *ThemeRepository) ListAll(ctx context.Context) ([]*models.Theme, error) {
	results, err := r.queries.ListAllThemes(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(results), nil
}

// Create creates a new theme
func (r *ThemeRepository) Create(ctx context.Context, theme *models.Theme) error {
	params := sqlc.CreateThemeParams{
		ID:          theme.ID,
		Name:        theme.Name,
		Description: theme.Description,
		ImageUrl:    theme.ImageURL,
		Price:       DecimalToNumeric(theme.Price),
		IsActive:    BoolPtr(theme.IsActive),
	}

	created, err := r.queries.CreateTheme(ctx, params)
	if err != nil {
		return err
	}

	*theme = *r.toModel(created)
	return nil
}

// Update updates an existing theme
func (r *ThemeRepository) Update(ctx context.Context, theme *models.Theme) error {
	params := sqlc.UpdateThemeParams{
		Column1:     theme.ID,
		Name:        theme.Name,
		Description: theme.Description,
		ImageUrl:    theme.ImageURL,
		Price:       DecimalToNumeric(theme.Price),
		IsActive:    BoolPtr(theme.IsActive),
	}

	updated, err := r.queries.UpdateTheme(ctx, params)
	if err != nil {
		return err
	}

	*theme = *r.toModel(updated)
	return nil
}

// Delete deletes a theme
func (r *ThemeRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteTheme(ctx, id)
}

// ToggleActive toggles theme active status
func (r *ThemeRepository) ToggleActive(ctx context.Context, id string) (*models.Theme, error) {
	result, err := r.queries.ToggleThemeActive(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toModel(result), nil
}

// Helper methods
func (r *ThemeRepository) toModel(row sqlc.Theme) *models.Theme {
	return &models.Theme{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		ImageURL:    row.ImageUrl,
		Price:       NumericToDecimal(row.Price),
		IsActive:    BoolVal(row.IsActive),
		CreatedAt:   TimestamptzToTime(row.CreatedAt),
		UpdatedAt:   TimestamptzToTime(row.UpdatedAt),
	}
}

func (r *ThemeRepository) toModels(rows []sqlc.Theme) []*models.Theme {
	themes := make([]*models.Theme, len(rows))
	for i, row := range rows {
		themes[i] = r.toModel(row)
	}
	return themes
}
