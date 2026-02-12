package services

import (
	"lqstudio-backend/pkg/errors"
)

var (
	// ErrIDMismatch occurs when URL ID doesn't match request body ID
	ErrIDMismatch = errors.NewBadRequestError("id in URL does not match id in request body")
)
