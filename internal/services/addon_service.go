package services

import (
	"context"

	"lqstudio-backend/internal/dto"
	"lqstudio-backend/internal/repositories"
)

// AddonService handles add-on operations
type AddonService struct {
	addonRepo *repositories.AddonRepository
}

// NewAddonService creates a new add-on service
func NewAddonService(addonRepo *repositories.AddonRepository) *AddonService {
	return &AddonService{addonRepo: addonRepo}
}

// GetByID retrieves an add-on by ID
func (s *AddonService) GetByID(ctx context.Context, id string) (*dto.AddonResponse, error) {
	addon, err := s.addonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToAddonResponse(addon), nil
}

// GetActive retrieves active add-ons for customers
func (s *AddonService) GetActive(ctx context.Context) ([]*dto.AddonResponse, error) {
	addons, err := s.addonRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToAddonsResponse(addons), nil
}

// ListAll retrieves all add-ons (admin)
func (s *AddonService) ListAll(ctx context.Context) ([]*dto.AddonResponse, error) {
	addons, err := s.addonRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToAddonsResponse(addons), nil
}

// Create creates a new add-on (admin)
func (s *AddonService) Create(ctx context.Context, req *dto.CreateAddonRequest) (*dto.AddonResponse, error) {
	// Convert DTO to domain model
	addon := req.ToAddonModel()

	// Create addon
	if err := s.addonRepo.Create(ctx, addon); err != nil {
		return nil, err
	}

	return dto.ToAddonResponse(addon), nil
}

// Update updates an add-on (admin)
func (s *AddonService) Update(ctx context.Context, id string, req *dto.UpdateAddonRequest) (*dto.AddonResponse, error) {
	// Ensure ID matches
	if id != req.ID {
		return nil, ErrIDMismatch
	}

	// Convert DTO to domain model
	addon := req.ToAddonModel()

	// Update addon
	if err := s.addonRepo.Update(ctx, addon); err != nil {
		return nil, err
	}

	return dto.ToAddonResponse(addon), nil
}

// Delete deletes an add-on (admin)
func (s *AddonService) Delete(ctx context.Context, id string) error {
	return s.addonRepo.Delete(ctx, id)
}

// ToggleActive toggles add-on active status (admin)
func (s *AddonService) ToggleActive(ctx context.Context, id string) (*dto.AddonResponse, error) {
	addon, err := s.addonRepo.ToggleActive(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToAddonResponse(addon), nil
}
