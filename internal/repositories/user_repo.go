package repositories

import (
	"context"

	"lqstudio-backend/internal/database/sqlc"
)

// User represents a user model for repository layer
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	Phone        string
	Role         string
	Notes        string
}

// UserRepository handles user data access
type UserRepository struct {
	queries *sqlc.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(queries *sqlc.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	result, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toModel(result), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	result, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return r.toModel(result), nil
}

// GetAdminByEmail retrieves an admin user by email
func (r *UserRepository) GetAdminByEmail(ctx context.Context, email string) (*User, error) {
	result, err := r.queries.GetAdminByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return r.toModel(result), nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *User) error {
	params := sqlc.CreateUserParams{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Phone:        StringPtr(user.Phone),
		PasswordHash: StringPtr(user.PasswordHash),
		Role:         user.Role,
	}

	created, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	*user = *r.toModel(created)
	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *User) error {
	params := sqlc.UpdateUserParams{
		Column1: user.ID,
		Name:    user.Name,
		Phone:   StringPtr(user.Phone),
	}

	updated, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return err
	}

	*user = *r.toModel(updated)
	return nil
}

// ListAll retrieves all users
func (r *UserRepository) ListAll(ctx context.Context) ([]*User, error) {
	results, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(results), nil
}

// Helper methods
func (r *UserRepository) toModel(row sqlc.User) *User {
	return &User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: StringVal(row.PasswordHash),
		Name:         row.Name,
		Phone:        StringVal(row.Phone),
		Role:         row.Role,
		Notes:        StringVal(row.Notes),
	}
}

func (r *UserRepository) toModels(rows []sqlc.User) []*User {
	users := make([]*User, len(rows))
	for i, row := range rows {
		users[i] = r.toModel(row)
	}
	return users
}
