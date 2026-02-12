package services

import (
	"context"

	"lqstudio-backend/internal/dto"
	"lqstudio-backend/internal/repositories"
)

// PackageService handles package operations
type PackageService struct {
	packageRepo *repositories.PackageRepository
}

// NewPackageService creates a new package service
func NewPackageService(packageRepo *repositories.PackageRepository) *PackageService {
	return &PackageService{packageRepo: packageRepo}
}

// GetByID retrieves a package by ID
func (s *PackageService) GetByID(ctx context.Context, id string) (*dto.PackageResponse, error) {
	pkg, err := s.packageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToPackageResponse(pkg), nil
}

// GetActive retrieves active packages for customers
func (s *PackageService) GetActive(ctx context.Context) ([]*dto.PackageResponse, error) {
	packages, err := s.packageRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToPackagesResponse(packages), nil
}

// ListAll retrieves all packages (admin)
func (s *PackageService) ListAll(ctx context.Context) ([]*dto.PackageResponse, error) {
	packages, err := s.packageRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return dto.ToPackagesResponse(packages), nil
}

// Create creates a new package (admin)
func (s *PackageService) Create(ctx context.Context, req *dto.CreatePackageRequest) (*dto.PackageResponse, error) {
	// Convert DTO to domain model
	pkg := req.ToPackageModel()

	// Validate package
	if err := pkg.ValidatePackage(); err != nil {
		return nil, err
	}

	// Create package
	if err := s.packageRepo.Create(ctx, pkg); err != nil {
		return nil, err
	}

	return dto.ToPackageResponse(pkg), nil
}

// Update updates a package (admin)
func (s *PackageService) Update(ctx context.Context, id string, req *dto.UpdatePackageRequest) (*dto.PackageResponse, error) {
	// Ensure ID matches
	if id != req.ID {
		return nil, ErrIDMismatch
	}

	// Convert DTO to domain model
	pkg := req.ToPackageModel()

	// Validate package
	if err := pkg.ValidatePackage(); err != nil {
		return nil, err
	}

	// Update package
	if err := s.packageRepo.Update(ctx, pkg); err != nil {
		return nil, err
	}

	return dto.ToPackageResponse(pkg), nil
}

// Delete deletes a package (admin)
func (s *PackageService) Delete(ctx context.Context, id string) error {
	return s.packageRepo.Delete(ctx, id)
}

// ToggleActive toggles package active status (admin)
func (s *PackageService) ToggleActive(ctx context.Context, id string) (*dto.PackageResponse, error) {
	pkg, err := s.packageRepo.ToggleActive(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToPackageResponse(pkg), nil
}
