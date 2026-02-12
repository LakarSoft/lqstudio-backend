package services

import (
	"context"

	"lqstudio-backend/internal/dto"
	"lqstudio-backend/internal/repositories"
)

// ThemeService handles theme operations
type ThemeService struct {
	themeRepo *repositories.ThemeRepository
}

// NewThemeService creates a new theme service
func NewThemeService(themeRepo *repositories.ThemeRepository) *ThemeService {
	return &ThemeService{themeRepo: themeRepo}
}

// GetByID retrieves a theme by ID
func (s *ThemeService) GetByID(ctx context.Context, id string) (*dto.ThemeResponse, error) {
	theme, err := s.themeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToThemeResponse(theme), nil
}

// GetActive retrieves active themes for customers
func (s *ThemeService) GetActive(ctx context.Context) ([]*dto.ThemeResponse, error) {
	themes, err := s.themeRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToThemesResponse(themes), nil
}

// ListAll retrieves all themes (admin)
func (s *ThemeService) ListAll(ctx context.Context) ([]*dto.ThemeResponse, error) {
	themes, err := s.themeRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToThemesResponse(themes), nil
}

// Create creates a new theme (admin)
func (s *ThemeService) Create(ctx context.Context, req *dto.CreateThemeRequest) (*dto.ThemeResponse, error) {
	// Convert DTO to domain model
	theme := req.ToThemeModel()

	// Create theme
	if err := s.themeRepo.Create(ctx, theme); err != nil {
		return nil, err
	}

	return dto.ToThemeResponse(theme), nil
}

// Update updates a theme (admin)
func (s *ThemeService) Update(ctx context.Context, id string, req *dto.UpdateThemeRequest) (*dto.ThemeResponse, error) {
	// Ensure ID matches
	if id != req.ID {
		return nil, ErrIDMismatch
	}

	// Convert DTO to domain model
	theme := req.ToThemeModel()

	// Update theme
	if err := s.themeRepo.Update(ctx, theme); err != nil {
		return nil, err
	}

	return dto.ToThemeResponse(theme), nil
}

// Delete deletes a theme (admin)
func (s *ThemeService) Delete(ctx context.Context, id string) error {
	return s.themeRepo.Delete(ctx, id)
}

// ToggleActive toggles theme active status (admin)
func (s *ThemeService) ToggleActive(ctx context.Context, id string) (*dto.ThemeResponse, error) {
	theme, err := s.themeRepo.ToggleActive(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToThemeResponse(theme), nil
}
