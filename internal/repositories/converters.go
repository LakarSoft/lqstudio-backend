package repositories

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// =============================================================================
// JSONB Conversions
// =============================================================================

// StringsToJSONB converts []string to JSONB bytes
func StringsToJSONB(strings []string) ([]byte, error) {
	if strings == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(strings)
}

// JSONBToStrings converts JSONB bytes to []string
func JSONBToStrings(data []byte) ([]string, error) {
	if data == nil {
		return []string{}, nil
	}

	var result []string
	if err := json.Unmarshal(data, &result); err != nil {
		return []string{}, err
	}

	return result, nil
}

// =============================================================================
// Numeric Conversions (pgtype.Numeric <-> decimal.Decimal)
// =============================================================================

// DecimalToNumeric converts decimal.Decimal to pgtype.Numeric
func DecimalToNumeric(d decimal.Decimal) pgtype.Numeric {
	var num pgtype.Numeric
	_ = num.Scan(d.String())
	return num
}

// NumericToDecimal converts pgtype.Numeric to decimal.Decimal
func NumericToDecimal(n pgtype.Numeric) decimal.Decimal {
	if !n.Valid {
		return decimal.Zero
	}

	// Convert pgtype.Numeric to string and parse as decimal
	str := n.Int.String()

	// Apply the scale (exponent)
	d, _ := decimal.NewFromString(str)

	if n.Exp != 0 {
		// Exp is negative for decimal places
		// e.g., Exp = -2 means divide by 100
		scale := decimal.NewFromInt(1)
		for i := int32(0); i < -n.Exp; i++ {
			scale = scale.Mul(decimal.NewFromInt(10))
		}
		d = d.Div(scale)
	}

	return d
}

// =============================================================================
// Timestamp Conversions
// =============================================================================

// TimestamptzToTime converts pgtype.Timestamptz to time.Time
func TimestamptzToTime(ts pgtype.Timestamptz) time.Time {
	if !ts.Valid {
		return time.Time{}
	}
	return ts.Time
}

// TimeToTimestamptz converts time.Time to pgtype.Timestamptz
func TimeToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// =============================================================================
// Date Conversions
// =============================================================================

// DateToTime converts pgtype.Date to time.Time
func DateToTime(d pgtype.Date) time.Time {
	if !d.Valid {
		return time.Time{}
	}
	return d.Time
}

// TimeToDate converts time.Time to pgtype.Date
func TimeToDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// =============================================================================
// Time of Day Conversions
// =============================================================================

// TimeToString converts pgtype.Time to 12-hour format string (e.g., "10:00 AM", "2:00 PM")
func TimeToString(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}

	// pgtype.Time stores microseconds since midnight
	hours := t.Microseconds / (60 * 60 * 1000000)
	minutes := (t.Microseconds % (60 * 60 * 1000000)) / (60 * 1000000)

	return time.Date(0, 1, 1, int(hours), int(minutes), 0, 0, time.UTC).Format("3:04 PM")
}

// StringToTime converts a 12-hour format string (e.g., "10:00 AM", "2:00 PM") to pgtype.Time
func StringToTime(s string) (pgtype.Time, error) {
	if s == "" {
		return pgtype.Time{Valid: false}, nil
	}

	t, err := time.Parse("3:04 PM", s)
	if err != nil {
		return pgtype.Time{Valid: false}, err
	}

	// Convert to microseconds since midnight
	microseconds := int64(t.Hour())*60*60*1000000 + int64(t.Minute())*60*1000000

	return pgtype.Time{
		Microseconds: microseconds,
		Valid:        true,
	}, nil
}

// =============================================================================
// Pointer Helpers
// =============================================================================

// StringPtr converts string to *string (nil if empty)
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// StringVal safely dereferences *string
func StringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BoolPtr converts bool to *bool
func BoolPtr(b bool) *bool {
	return &b
}

// BoolVal safely dereferences *bool
func BoolVal(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Int32Ptr converts int32 to *int32
func Int32Ptr(i int32) *int32 {
	return &i
}

// Int32Val safely dereferences *int32
func Int32Val(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}
