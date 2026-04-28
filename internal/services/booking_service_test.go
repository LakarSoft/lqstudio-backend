package services

import (
	"reflect"
	"testing"

	"lqstudio-backend/internal/dto"
	"lqstudio-backend/internal/models"
)

func TestGenerateTimeSlotsUsesRayaInterval(t *testing.T) {
	got := generateTimeSlots(10, 11, models.RayaSlotDurationMinutes)
	want := []string{"10:00 AM", "10:20 AM", "10:40 AM"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestGenerateTimeSlotsUsesConvocationInterval(t *testing.T) {
	got := generateTimeSlots(10, 11, models.ConvocationSlotDurationMinutes)
	want := []string{"10:00 AM", "10:30 AM"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestGenerateTimeSlotsRejectsInvalidInterval(t *testing.T) {
	got := generateTimeSlots(10, 11, 0)

	if len(got) != 0 {
		t.Fatalf("expected no slots for invalid interval, got %v", got)
	}
}

func TestValidateBookingSlotTimesAcceptsConvocationInterval(t *testing.T) {
	pkg := &models.Package{
		ID:              "pkg-convo",
		Module:          models.ModuleConvocation,
		DurationMinutes: 30,
	}
	slots := []dto.SlotRequest{
		{Date: "2026-05-01", Time: "10:30 AM", ThemeID: "theme-convo"},
	}

	if err := validateBookingSlotTimes(slots, pkg, 10, 18); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateBookingSlotTimesRejectsMisalignedConvocationTime(t *testing.T) {
	pkg := &models.Package{
		ID:              "pkg-convo",
		Module:          models.ModuleConvocation,
		DurationMinutes: 30,
	}
	slots := []dto.SlotRequest{
		{Date: "2026-05-01", Time: "10:20 AM", ThemeID: "theme-convo"},
	}

	if err := validateBookingSlotTimes(slots, pkg, 10, 18); err == nil {
		t.Fatal("expected error for misaligned convocation time")
	}
}

func TestValidateBookingSlotTimesAcceptsRayaInterval(t *testing.T) {
	pkg := &models.Package{
		ID:              "pkg-raya",
		Module:          models.ModuleRaya,
		DurationMinutes: 20,
	}
	slots := []dto.SlotRequest{
		{Date: "2026-05-01", Time: "10:20 AM", ThemeID: "theme-raya"},
	}

	if err := validateBookingSlotTimes(slots, pkg, 10, 18); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateBookingSlotTimesRejectsInvalidTimeFormat(t *testing.T) {
	pkg := &models.Package{
		ID:              "pkg-convo",
		Module:          models.ModuleConvocation,
		DurationMinutes: 30,
	}
	slots := []dto.SlotRequest{
		{Date: "2026-05-01", Time: "10:30", ThemeID: "theme-convo"},
	}

	if err := validateBookingSlotTimes(slots, pkg, 10, 18); err == nil {
		t.Fatal("expected error for invalid time format")
	}
}

func TestValidateBookingSlotTimesRejectsOutOfHoursTime(t *testing.T) {
	pkg := &models.Package{
		ID:              "pkg-convo",
		Module:          models.ModuleConvocation,
		DurationMinutes: 30,
	}
	slots := []dto.SlotRequest{
		{Date: "2026-05-01", Time: "6:00 PM", ThemeID: "theme-convo"},
	}

	if err := validateBookingSlotTimes(slots, pkg, 10, 18); err == nil {
		t.Fatal("expected error for out-of-hours time")
	}
}

func TestValidateBookingSlotTimesNormalizesWhitespace(t *testing.T) {
	pkg := &models.Package{
		ID:              "pkg-convo",
		Module:          models.ModuleConvocation,
		DurationMinutes: 30,
	}
	slots := []dto.SlotRequest{
		{Date: " 2026-05-01 ", Time: " 10:30 AM ", ThemeID: " theme-convo "},
	}

	if err := validateBookingSlotTimes(slots, pkg, 10, 18); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if slots[0].Date != "2026-05-01" || slots[0].Time != "10:30 AM" || slots[0].ThemeID != "theme-convo" {
		t.Fatalf("expected normalized slot, got %+v", slots[0])
	}
}
