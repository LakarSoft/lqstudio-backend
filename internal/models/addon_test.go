package models

import "testing"

func TestAddonValidateAddonNormalizesModule(t *testing.T) {
	addon := &AddOn{
		ID:     "addon-1",
		Module: "  Convocation  ",
		Name:   "Frame",
	}

	if err := addon.ValidateAddon(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if addon.Module != "convocation" {
		t.Fatalf("expected normalized module convocation, got %q", addon.Module)
	}
}

func TestAddonValidateAddonRejectsInvalidModule(t *testing.T) {
	addon := &AddOn{
		ID:     "addon-1",
		Module: "wedding",
		Name:   "Frame",
	}

	err := addon.ValidateAddon()
	if err != ErrInvalidAddonModule {
		t.Fatalf("expected ErrInvalidAddonModule, got %v", err)
	}
}
