package services

import (
	"reflect"
	"testing"

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
