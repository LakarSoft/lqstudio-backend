package handlers

import "testing"

func TestNormalizeBookingModuleFilter(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{name: "empty", input: "", expect: ""},
		{name: "trim and lowercase", input: "  Convocation  ", expect: "convocation"},
		{name: "already normalized", input: "raya", expect: "raya"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeBookingModuleFilter(tt.input)
			if got != tt.expect {
				t.Fatalf("expected %q, got %q", tt.expect, got)
			}
		})
	}
}
