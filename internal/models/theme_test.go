package models

import "testing"

func TestThemeValidateThemeNormalizesModule(t *testing.T) {
	theme := &Theme{
		ID:     "theme-1",
		Module: "  Convocation  ",
		Name:   "Convocation Theme",
	}

	if err := theme.ValidateTheme(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if theme.Module != "convocation" {
		t.Fatalf("expected normalized module convocation, got %q", theme.Module)
	}
}

func TestThemeValidateThemeRejectsInvalidModule(t *testing.T) {
	theme := &Theme{
		ID:     "theme-1",
		Module: "wedding",
		Name:   "Theme",
	}

	err := theme.ValidateTheme()
	if err != ErrInvalidThemeModule {
		t.Fatalf("expected ErrInvalidThemeModule, got %v", err)
	}
}
