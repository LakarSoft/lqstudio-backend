package sqlc

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

// NewNullString creates a sql.NullString
func NewNullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// NewNullBool creates a nullable bool pointer
func NewNullBool(b bool) *bool {
	return &b
}

// NewNullInt32 creates a nullable int32 pointer
func NewNullInt32(i int32) *int32 {
	return &i
}

// NumericToFloat64 converts pgtype.Numeric to float64
func NumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}

	// Convert using big.Float for precision
	bigFloat := new(big.Float)
	bigFloat.SetInt(n.Int)

	// Apply the scale (exponent)
	if n.Exp != 0 {
		scale := new(big.Float).SetFloat64(1)
		for i := int32(0); i < -n.Exp; i++ {
			scale.Mul(scale, big.NewFloat(10))
		}
		bigFloat.Quo(bigFloat, scale)
	}

	result, _ := bigFloat.Float64()
	return result
}

// Float64ToNumeric converts float64 to pgtype.Numeric
func Float64ToNumeric(f float64) pgtype.Numeric {
	// Convert float to string with full precision
	str := big.NewFloat(f).Text('f', 2) // 2 decimal places for currency

	var num pgtype.Numeric
	_ = num.Scan(str)
	return num
}

// PtrToString safely dereferences a string pointer
func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// PtrToBool safely dereferences a bool pointer
func PtrToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// PtrToInt32 safely dereferences an int32 pointer
func PtrToInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}
