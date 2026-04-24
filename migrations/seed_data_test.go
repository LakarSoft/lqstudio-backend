package migrations

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestConvocationSeedDataCounts(t *testing.T) {
	sql := readMigration(t, "00004_add_modules_column.sql")

	assertCount(t, sql, "('pkg-convo-", 3, "convocation packages")
	assertCount(t, sql, "('theme-convo-", 3, "convocation themes")
	assertCount(t, sql, "('addon-convo-", 3, "convocation addons")
	assertCount(t, sql, "ALTER COLUMN module SET DEFAULT 'raya';", 3, "module defaults")
	assertCount(t, sql, "ON CONFLICT (id) DO UPDATE SET", 3, "idempotent seed upserts")
}

func TestConvocationPackageDurationsUseThirtyMinuteSlots(t *testing.T) {
	sql := readMigration(t, "00004_add_modules_column.sql")
	rowPattern := regexp.MustCompile(`\('pkg-convo-[^']+', 'convocation', '[^']+', '[^']+', ([0-9]+),`)
	matches := rowPattern.FindAllStringSubmatch(sql, -1)
	if len(matches) != 3 {
		t.Fatalf("expected 3 convocation package rows, got %d", len(matches))
	}

	for _, match := range matches {
		duration, err := strconv.Atoi(match[1])
		if err != nil {
			t.Fatalf("invalid package duration %q: %v", match[1], err)
		}
		if duration <= 0 || duration%30 != 0 {
			t.Fatalf("convocation package duration must be a positive multiple of 30, got %d", duration)
		}
	}
}

func readMigration(t *testing.T, name string) string {
	t.Helper()

	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read migration %s: %v", name, err)
	}
	return string(data)
}

func assertCount(t *testing.T, haystack, needle string, expected int, label string) {
	t.Helper()

	if actual := strings.Count(haystack, needle); actual != expected {
		t.Fatalf("expected %d %s, got %d", expected, label, actual)
	}
}
