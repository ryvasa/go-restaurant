package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db}
}

func (r *UserRepositoryImpl) GetAll(ctx context.Context) ([]domain.User, error) {
	users := []domain.User{}
	rows, err := r.db.QueryContext(ctx, "SELECT id,name,email,phone,role,created_at,updated_at FROM users WHERE deleted = false AND deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user domain.User) (domain.User, error) {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (id,name,email,password,phone,role) VALUES (?, ?, ?, ?, ?, ?)",
		user.Id, user.Name, user.Email, user.Password, user.Phone, user.Role)
	if err != nil {
		return domain.User{}, err
	}

	createdUser, err := r.Get(ctx, user.Id.String())
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (r *UserRepositoryImpl) Get(ctx context.Context, id string) (domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRowContext(ctx, "SELECT id,name,email,phone,role,created_at,updated_at FROM users WHERE id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user domain.User) (domain.User, error) {
	// Ambil data user yang ada
	existingUser, err := r.Get(ctx, user.Id.String())
	if err != nil {
		return domain.User{}, err
	}

	// Update hanya field yang tidak kosong
	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Password != "" {
		existingUser.Password = user.Password
	}
	if user.Phone != "" {
		existingUser.Phone = user.Phone
	}
	if user.Role != "" {
		existingUser.Role = user.Role
	}

	// Eksekusi query update
	_, err = r.db.ExecContext(ctx, "UPDATE users SET name = ?, email = ?, password = ?, phone = ?, role = ? WHERE id = ?",
		existingUser.Name, existingUser.Email, existingUser.Password, existingUser.Phone, existingUser.Role, existingUser.Id)
	if err != nil {
		return domain.User{}, err
	}

	return r.Get(ctx, user.Id.String())
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRowContext(ctx, "SELECT id,name,email,phone,role,created_at,updated_at FROM users WHERE email = ?", email).Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET deleted = ?, deleted_at = ? WHERE id = ?", true, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) Restore(ctx context.Context, id string) (domain.User, error) {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET deleted = ?, deleted_at = ? WHERE id = ?", false, nil, id)
	if err != nil {
		return domain.User{}, err
	}
	user, _ := r.Get(ctx, id)
	return user, nil
}

func (r *UserRepositoryImpl) GetDeletedUserById(ctx context.Context, id string) ([]domain.User, error) {
	users := []domain.User{}
	rows, err := r.db.QueryContext(ctx, "SELECT id,name,email,phone,role,created_at,updated_at FROM users WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
