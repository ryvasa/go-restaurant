package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type UserRepositoryImpl struct {
	db DB
}

func NewUserRepository(db DB) UserRepository {
	return &UserRepositoryImpl{db}
}

func (r *UserRepositoryImpl) GetAll(ctx context.Context) ([]domain.User, error) {
	users := []domain.User{}
	query := `SELECT id,name,email,phone,role,created_at,updated_at FROM users WHERE deleted = false AND deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query)
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

func (r *UserRepositoryImpl) Create(ctx context.Context, user domain.User) error {
	query := `INSERT INTO users (id,name,email,password,role) VALUES (?, ?,  ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

func (r *UserRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user := domain.User{}
	var phone sql.NullString
	query := `SELECT id, name, email, phone, role, created_at, updated_at FROM users WHERE id = ? AND deleted = false AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Name, &user.Email, &phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return domain.User{}, err
	}

	if phone.Valid {
		user.Phone = &phone.String
	} else {
		user.Phone = nil
	}

	return user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, id uuid.UUID, user domain.User) error {
	query := `UPDATE users SET name = ?, email = ?, password = ?, phone = ?, role = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query,
		user.Name, user.Email, user.Password, user.Phone, user.Role, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}

	return nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user := domain.User{}
	query := `SELECT id,name,password,email,phone,role,created_at,updated_at FROM users WHERE email = ?`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET deleted = ?, deleted_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, true, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}

	return nil
}

func (r *UserRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET deleted = ?, deleted_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, false, nil, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}

	return nil
}

func (r *UserRepositoryImpl) GetDeletedUserById(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user := domain.User{}
	query := `SELECT id,name,email,phone,role,created_at,updated_at FROM users WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
