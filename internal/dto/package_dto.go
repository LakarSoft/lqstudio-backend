package dto

import (
	"lqstudio-backend/internal/models"

	"github.com/shopspring/decimal"
)

// PackageResponse represents package information for API responses
// Matches the frontend Package entity structure with camelCase JSON tags
type PackageResponse struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description,omitempty"`
	DurationMinutes int32    `json:"durationMinutes"`
	Price           float64  `json:"price"`            // Base price in MYR
	Discount        float64  `json:"discount"`         // Discount percentage (0-100)
	Offers          []string `json:"offers,omitempty"` // List of features/benefits
	ImageURL        string   `json:"imageUrl,omitempty"`
	IsActive        bool     `json:"isActive"`
	BasePrice       float64  `json:"basePrice"`  // Same as price (for frontend compatibility)
	FinalPrice      float64  `json:"finalPrice"` // Calculated: price * (1 - discount/100)
}

// CreatePackageRequest for admin package creation
// Matches the frontend expectations with camelCase JSON tags
type CreatePackageRequest struct {
	ID              string   `json:"id" validate:"required,max=50"`
	Name            string   `json:"name" validate:"required"`
	Description     string   `json:"description"`
	DurationMinutes int32    `json:"durationMinutes" validate:"required,min=15"`
	Price           float64  `json:"price" validate:"required,min=0"`
	Discount        float64  `json:"discount" validate:"min=0,max=100"`
	Offers          []string `json:"offers"`
	ImageURL        string   `json:"imageUrl"`
}

// UpdatePackageRequest for admin package updates
// Matches the frontend expectations with camelCase JSON tags
type UpdatePackageRequest struct {
	ID              string   `json:"id" validate:"required,max=50"`
	Name            string   `json:"name" validate:"required"`
	Description     string   `json:"description"`
	DurationMinutes int32    `json:"durationMinutes" validate:"required,min=15"`
	Price           float64  `json:"price" validate:"required,min=0"`
	Discount        float64  `json:"discount" validate:"min=0,max=100"`
	Offers          []string `json:"offers"`
	ImageURL        string   `json:"imageUrl"`
	IsActive        bool     `json:"isActive"`
}

// ToPackageResponse converts a domain model Package to a PackageResponse DTO
// Handles decimal.Decimal to float64 conversion
func ToPackageResponse(pkg *models.Package) *PackageResponse {
	if pkg == nil {
		return nil
	}

	// Convert decimal.Decimal to float64
	price, _ := pkg.Price.Float64()
	discount, _ := pkg.Discount.Float64()
	finalPrice, _ := pkg.FinalPrice().Float64()

	return &PackageResponse{
		ID:              pkg.ID,
		Name:            pkg.Name,
		Description:     pkg.Description,
		DurationMinutes: pkg.DurationMinutes,
		Price:           price,
		Discount:        discount,
		Offers:          pkg.Offers,
		ImageURL:        pkg.ImageURL,
		IsActive:        pkg.IsActive,
		BasePrice:       price,
		FinalPrice:      finalPrice,
	}
}

// ToPackagesResponse converts a slice of domain model Packages to a slice of PackageResponse DTOs
func ToPackagesResponse(packages []*models.Package) []*PackageResponse {
	if packages == nil {
		return []*PackageResponse{}
	}

	responses := make([]*PackageResponse, len(packages))
	for i, pkg := range packages {
		responses[i] = ToPackageResponse(pkg)
	}
	return responses
}

// ToPackageModel converts CreatePackageRequest to domain model Package
// Handles float64 to decimal.Decimal conversion
func (r *CreatePackageRequest) ToPackageModel() *models.Package {
	return &models.Package{
		ID:              r.ID,
		Name:            r.Name,
		Description:     r.Description,
		DurationMinutes: r.DurationMinutes,
		Price:           decimal.NewFromFloat(r.Price),
		Discount:        decimal.NewFromFloat(r.Discount),
		Offers:          r.Offers,
		ImageURL:        r.ImageURL,
		IsActive:        true, // New packages are active by default
	}
}

// ToPackageModel converts UpdatePackageRequest to domain model Package
// Handles float64 to decimal.Decimal conversion
func (r *UpdatePackageRequest) ToPackageModel() *models.Package {
	return &models.Package{
		ID:              r.ID,
		Name:            r.Name,
		Description:     r.Description,
		DurationMinutes: r.DurationMinutes,
		Price:           decimal.NewFromFloat(r.Price),
		Discount:        decimal.NewFromFloat(r.Discount),
		Offers:          r.Offers,
		ImageURL:        r.ImageURL,
		IsActive:        r.IsActive,
	}
}
