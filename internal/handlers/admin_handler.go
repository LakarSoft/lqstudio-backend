package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"lqstudio-backend/internal/config"
	"lqstudio-backend/internal/dto"
	"lqstudio-backend/internal/services"
	apperrors "lqstudio-backend/pkg/errors"
	"lqstudio-backend/pkg/upload"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	packageService *services.PackageService
	themeService   *services.ThemeService
	addonService   *services.AddonService
	bookingService *services.BookingService
	uploadConfig   *config.UploadConfig
}

func NewAdminHandler(
	packageService *services.PackageService,
	themeService *services.ThemeService,
	addonService *services.AddonService,
	bookingService *services.BookingService,
	uploadConfig *config.UploadConfig,
) *AdminHandler {
	return &AdminHandler{
		packageService: packageService,
		themeService:   themeService,
		addonService:   addonService,
		bookingService: bookingService,
		uploadConfig:   uploadConfig,
	}
}

// ================================================================================
// Public Handlers (Customer-facing)
// ================================================================================

// GetActivePackages godoc
// @Summary Get active packages
// @Description Get list of all active packages available for booking
// @Tags packages
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=[]dto.PackageResponse}
// @Failure 500 {object} dto.ApiResponse
// @Router /api/packages [get]
func (h *AdminHandler) GetActivePackages(c echo.Context) error {
	packages, err := h.packageService.GetActive(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, packages, "Active packages retrieved successfully")
}

// GetActiveThemes godoc
// @Summary Get active themes
// @Description Get list of all active themes available for booking
// @Tags themes
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=[]dto.ThemeResponse}
// @Failure 500 {object} dto.ApiResponse
// @Router /api/themes [get]
func (h *AdminHandler) GetActiveThemes(c echo.Context) error {
	themes, err := h.themeService.GetActive(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, themes, "Active themes retrieved successfully")
}

// GetActiveAddons godoc
// @Summary Get active add-ons
// @Description Get list of all active add-ons available for booking
// @Tags addons
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=[]dto.AddonResponse}
// @Failure 500 {object} dto.ApiResponse
// @Router /api/addons [get]
func (h *AdminHandler) GetActiveAddons(c echo.Context) error {
	addons, err := h.addonService.GetActive(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, addons, "Active add-ons retrieved successfully")
}

// ================================================================================
// Package Management Handlers (Admin)
// ================================================================================

// GetAllPackages godoc
// @Summary Get all packages (Admin)
// @Description Get list of all packages including inactive ones (admin only)
// @Tags admin-packages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.ApiResponse{data=[]dto.PackageResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/packages [get]
func (h *AdminHandler) GetAllPackages(c echo.Context) error {
	packages, err := h.packageService.ListAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, packages, "Packages retrieved successfully")
}

// GetPackage godoc
// @Summary Get package by ID (Admin)
// @Description Get a single package by ID (admin only)
// @Tags admin-packages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Package ID"
// @Success 200 {object} dto.ApiResponse{data=dto.PackageResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/packages/{id} [get]
func (h *AdminHandler) GetPackage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid package ID")
	}

	pkg, err := h.packageService.GetByID(c.Request().Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, pkg, "Package retrieved successfully")
}

// CreatePackage godoc
// @Summary Create new package (Admin)
// @Description Create a new package (admin only)
// @Tags admin-packages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body dto.CreatePackageRequest true "Package data"
// @Success 201 {object} dto.ApiResponse{data=dto.PackageResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Router /api/admin/packages [post]
func (h *AdminHandler) CreatePackage(c echo.Context) error {
	var req dto.CreatePackageRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	pkg, err := h.packageService.Create(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return SendCreated(c, pkg, "Package created successfully")
}

// UpdatePackage godoc
// @Summary Update package (Admin)
// @Description Update an existing package (admin only)
// @Tags admin-packages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Package ID"
// @Param request body dto.UpdatePackageRequest true "Package data"
// @Success 200 {object} dto.ApiResponse{data=dto.PackageResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/admin/packages/{id} [put]
func (h *AdminHandler) UpdatePackage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid package ID")
	}

	var req dto.UpdatePackageRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	pkg, err := h.packageService.Update(c.Request().Context(), id, &req)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return SendOK(c, pkg, "Package updated successfully")
}

// DeletePackage godoc
// @Summary Delete package (Admin)
// @Description Delete a package (admin only)
// @Tags admin-packages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Package ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/packages/{id} [delete]
func (h *AdminHandler) DeletePackage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid package ID")
	}

	err := h.packageService.Delete(c.Request().Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendDeleted(c, "Package deleted successfully")
}

// TogglePackageActive godoc
// @Summary Toggle package active status (Admin)
// @Description Toggle package active/inactive status (admin only)
// @Tags admin-packages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Package ID"
// @Success 200 {object} dto.ApiResponse{data=dto.PackageResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/packages/{id}/toggle-active [patch]
func (h *AdminHandler) TogglePackageActive(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid package ID")
	}

	pkg, err := h.packageService.ToggleActive(c.Request().Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, pkg, "Package status toggled successfully")
}

// UploadPackageImage godoc
// @Summary Upload package image (Admin)
// @Description Upload an image for a package (admin only)
// @Tags admin-packages
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Package ID"
// @Param file formData file true "Package image"
// @Success 200 {object} dto.ApiResponse{data=dto.PackageResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/packages/{id}/image [post]
func (h *AdminHandler) UploadPackageImage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		appErr := apperrors.NewBadRequestError("Package ID is required")
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		appErr := apperrors.NewBadRequestError("File is required")
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	if err := upload.ValidateFile(fileHeader, h.uploadConfig.MaxFileSize, h.uploadConfig.AllowedTypes); err != nil {
		if strings.Contains(err.Error(), "exceeds") {
			appErr := apperrors.NewFileTooLargeError(h.uploadConfig.MaxFileSize)
			return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
		}
		appErr := apperrors.NewInvalidFileTypeError(h.uploadConfig.AllowedTypes)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	file, err := fileHeader.Open()
	if err != nil {
		appErr := apperrors.NewUploadFailedError(err)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}
	defer file.Close()

	uniqueFilename := upload.GenerateUniqueFilename(fileHeader.Filename, "package")
	urlPath, err := upload.SaveFile(file, uniqueFilename, h.uploadConfig.StoragePath, "packages")
	if err != nil {
		appErr := apperrors.NewUploadFailedError(err)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	pkg, err := h.packageService.UpdateImageURL(c.Request().Context(), id, urlPath)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
		}
		appErr := apperrors.NewInternalError("Failed to update package image", err)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	return SendOK(c, pkg, "Package image uploaded successfully")
}

// ================================================================================
// Theme Management Handlers (Admin)
// ================================================================================

// GetAllThemes godoc
// @Summary Get all themes (Admin)
// @Description Get list of all themes including inactive ones (admin only)
// @Tags admin-themes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.ApiResponse{data=[]dto.ThemeResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/themes [get]
func (h *AdminHandler) GetAllThemes(c echo.Context) error {
	themes, err := h.themeService.ListAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, themes, "Themes retrieved successfully")
}

// GetTheme godoc
// @Summary Get theme by ID (Admin)
// @Description Get a single theme by ID (admin only)
// @Tags admin-themes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Theme ID"
// @Success 200 {object} dto.ApiResponse{data=dto.ThemeResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/themes/{id} [get]
func (h *AdminHandler) GetTheme(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid theme ID")
	}

	theme, err := h.themeService.GetByID(c.Request().Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, theme, "Theme retrieved successfully")
}

// CreateTheme godoc
// @Summary Create new theme (Admin)
// @Description Create a new theme (admin only)
// @Tags admin-themes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body dto.CreateThemeRequest true "Theme data"
// @Success 201 {object} dto.ApiResponse{data=dto.ThemeResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Router /api/admin/themes [post]
func (h *AdminHandler) CreateTheme(c echo.Context) error {
	var req dto.CreateThemeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	theme, err := h.themeService.Create(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return SendCreated(c, theme, "Theme created successfully")
}

// UpdateTheme godoc
// @Summary Update theme (Admin)
// @Description Update an existing theme (admin only)
// @Tags admin-themes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Theme ID"
// @Param request body dto.UpdateThemeRequest true "Theme data"
// @Success 200 {object} dto.ApiResponse{data=dto.ThemeResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/admin/themes/{id} [put]
func (h *AdminHandler) UpdateTheme(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid theme ID")
	}

	var req dto.UpdateThemeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	theme, err := h.themeService.Update(c.Request().Context(), id, &req)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return SendOK(c, theme, "Theme retrieved successfully")
}

// DeleteTheme godoc
// @Summary Delete theme (Admin)
// @Description Delete a theme (admin only)
// @Tags admin-themes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Theme ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/themes/{id} [delete]
func (h *AdminHandler) DeleteTheme(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid theme ID")
	}

	err := h.themeService.Delete(c.Request().Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendDeleted(c, "Theme deleted successfully")
}

// ToggleThemeActive godoc
// @Summary Toggle theme active status (Admin)
// @Description Toggle theme active/inactive status (admin only)
// @Tags admin-themes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Theme ID"
// @Success 200 {object} dto.ApiResponse{data=dto.ThemeResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/themes/{id}/toggle-active [patch]
func (h *AdminHandler) ToggleThemeActive(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid theme ID")
	}

	theme, err := h.themeService.ToggleActive(c.Request().Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, theme, "Theme retrieved successfully")
}

// UploadThemeImage godoc
// @Summary Upload theme image (Admin)
// @Description Upload an image for a theme (admin only)
// @Tags admin-themes
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Theme ID"
// @Param file formData file true "Theme image"
// @Success 200 {object} dto.ApiResponse{data=dto.ThemeResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/themes/{id}/image [post]
func (h *AdminHandler) UploadThemeImage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		appErr := apperrors.NewBadRequestError("Theme ID is required")
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		appErr := apperrors.NewBadRequestError("File is required")
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	if err := upload.ValidateFile(fileHeader, h.uploadConfig.MaxFileSize, h.uploadConfig.AllowedTypes); err != nil {
		if strings.Contains(err.Error(), "exceeds") {
			appErr := apperrors.NewFileTooLargeError(h.uploadConfig.MaxFileSize)
			return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
		}
		appErr := apperrors.NewInvalidFileTypeError(h.uploadConfig.AllowedTypes)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	file, err := fileHeader.Open()
	if err != nil {
		appErr := apperrors.NewUploadFailedError(err)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}
	defer file.Close()

	uniqueFilename := upload.GenerateUniqueFilename(fileHeader.Filename, "theme")
	urlPath, err := upload.SaveFile(file, uniqueFilename, h.uploadConfig.StoragePath, "themes")
	if err != nil {
		appErr := apperrors.NewUploadFailedError(err)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	theme, err := h.themeService.UpdateImageURL(c.Request().Context(), id, urlPath)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
		}
		appErr := apperrors.NewInternalError("Failed to update theme image", err)
		return c.JSON(appErr.StatusCode, dto.NewErrorResponseWithCode(appErr.Message, appErr.Code))
	}

	return SendOK(c, theme, "Theme image uploaded successfully")
}

// ================================================================================
// Add-On Management Handlers (Admin)
// ================================================================================

// GetAllAddons godoc
// @Summary Get all add-ons (Admin)
// @Description Get list of all add-ons including inactive ones (admin only)
// @Tags admin-addons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.ApiResponse{data=[]dto.AddonResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/addons [get]
func (h *AdminHandler) GetAllAddons(c echo.Context) error {
	addons, err := h.addonService.ListAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, addons, "Add-ons retrieved successfully")
}

// GetAddon godoc
// @Summary Get add-on by ID (Admin)
// @Description Get a single add-on by ID (admin only)
// @Tags admin-addons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Add-on ID"
// @Success 200 {object} dto.ApiResponse{data=dto.AddonResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/addons/{id} [get]
func (h *AdminHandler) GetAddon(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid addon ID")
	}

	addon, err := h.addonService.GetByID(c.Request().Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, addon, "Add-on retrieved successfully")
}

// CreateAddon godoc
// @Summary Create new add-on (Admin)
// @Description Create a new add-on (admin only)
// @Tags admin-addons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body dto.CreateAddonRequest true "Add-on data"
// @Success 201 {object} dto.ApiResponse{data=dto.AddonResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Router /api/admin/addons [post]
func (h *AdminHandler) CreateAddon(c echo.Context) error {
	var req dto.CreateAddonRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	addon, err := h.addonService.Create(c.Request().Context(), &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return SendCreated(c, addon, "Add-on created successfully")
}

// UpdateAddon godoc
// @Summary Update add-on (Admin)
// @Description Update an existing add-on (admin only)
// @Tags admin-addons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Add-on ID"
// @Param request body dto.UpdateAddonRequest true "Add-on data"
// @Success 200 {object} dto.ApiResponse{data=dto.AddonResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/admin/addons/{id} [put]
func (h *AdminHandler) UpdateAddon(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid addon ID")
	}

	var req dto.UpdateAddonRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	addon, err := h.addonService.Update(c.Request().Context(), id, &req)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return SendOK(c, addon, "Add-on retrieved successfully")
}

// DeleteAddon godoc
// @Summary Delete add-on (Admin)
// @Description Delete an add-on (admin only)
// @Tags admin-addons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Add-on ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/addons/{id} [delete]
func (h *AdminHandler) DeleteAddon(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid addon ID")
	}

	err := h.addonService.Delete(c.Request().Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendDeleted(c, "Add-on deleted successfully")
}

// ToggleAddonActive godoc
// @Summary Toggle add-on active status (Admin)
// @Description Toggle add-on active/inactive status (admin only)
// @Tags admin-addons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Add-on ID"
// @Success 200 {object} dto.ApiResponse{data=dto.AddonResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/addons/{id}/toggle-active [patch]
func (h *AdminHandler) ToggleAddonActive(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid addon ID")
	}

	addon, err := h.addonService.ToggleActive(c.Request().Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, addon, "Add-on retrieved successfully")
}

// ================================================================================
// Booking Management Handlers (Admin)
// ================================================================================

// ListBookings godoc
// @Summary List all bookings (Admin)
// @Description Get list of all bookings with filters, search, sorting, and pagination (admin only)
// @Tags admin-bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param status query string false "Filter by status" Enums(PENDING, APPROVED, REJECTED, COMPLETED)
// @Param email query string false "Filter by customer email (partial match)"
// @Param packageId query string false "Filter by package ID"
// @Param themeId query string false "Filter by theme ID"
// @Param slotDate query string false "Filter by slot date (YYYY-MM-DD)"
// @Param dateFrom query string false "Filter bookings created from this date (YYYY-MM-DD)"
// @Param dateTo query string false "Filter bookings created up to this date (YYYY-MM-DD)"
// @Param search query string false "Search customer name, email, or phone"
// @Param sortBy query string false "Sort field" Enums(createdAt, updatedAt, totalPrice, status) default(createdAt)
// @Param order query string false "Sort order" Enums(asc, desc) default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.ApiResponse{data=dto.PaginatedBookingsResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/bookings [get]
func (h *AdminHandler) ListBookings(c echo.Context) error {
	var filters dto.BookingFilters
	if err := c.Bind(&filters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	if filters.Status != "" {
		valid := map[string]bool{"PENDING": true, "APPROVED": true, "REJECTED": true, "COMPLETED": true}
		if !valid[filters.Status] {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid status parameter. must be PENDING, APPROVED, REJECTED, or COMPLETED")
		}
	}

	if filters.SortBy != "" {
		valid := map[string]bool{"createdAt": true, "updatedAt": true, "totalPrice": true, "status": true}
		if !valid[filters.SortBy] {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid sortBy parameter. must be createdAt, updatedAt, totalPrice, or status")
		}
	}

	if filters.Order != "" && filters.Order != "asc" && filters.Order != "desc" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order parameter. must be asc or desc")
	}

	result, err := h.bookingService.ListBookings(c.Request().Context(), &filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, result, "Bookings retrieved successfully")
}

// GetBooking godoc
// @Summary Get booking by ID (Admin)
// @Description Get a single booking by ID (admin only)
// @Tags admin-bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Booking ID"
// @Success 200 {object} dto.ApiResponse{data=dto.BookingResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/bookings/{id} [get]
func (h *AdminHandler) GetBooking(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid booking ID")
	}

	booking, err := h.bookingService.GetBookingByID(c.Request().Context(), id)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return SendOK(c, booking, "Booking retrieved successfully")
}

// UpdateBookingStatus godoc
// @Summary Update booking status (Admin)
// @Description Update booking status and send customer notification email (admin only)
// @Tags admin-bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Booking ID"
// @Param request body dto.UpdateBookingStatusRequest true "Status update data"
// @Success 200 {object} dto.ApiResponse{data=dto.UpdateBookingStatusResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/admin/bookings/{id}/status [patch]
func (h *AdminHandler) UpdateBookingStatus(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid booking ID")
	}

	var req dto.UpdateBookingStatusRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	response, err := h.bookingService.UpdateBookingStatus(c.Request().Context(), id, &req)
	if err != nil {
		if contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if contains(err.Error(), "invalid") {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Build response message based on email notification status
	message := "Booking status updated successfully"
	if req.Status == "APPROVED" || req.Status == "REJECTED" {
		if response.EmailNotificationSent {
			message = "Booking status updated successfully and customer has been notified via email"
		} else if response.EmailError != "" {
			message = fmt.Sprintf("Booking status updated successfully, but failed to send customer notification: %s. Please contact the customer directly.", response.EmailError)
		}
	}

	return SendOK(c, response, message)
}

// ================================================================================
// Helper Functions
// ================================================================================

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr)+1 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
