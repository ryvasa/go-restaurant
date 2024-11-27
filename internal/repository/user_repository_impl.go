package repository

import (
	"database/sql"
	"time"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) GetAll(tx *sql.Tx) ([]domain.User, error) {
	users := []domain.User{}
	rows, err := tx.Query("SELECT id,name,email,phone,role,created_at,updated_at FROM users WHERE deleted = false AND deleted_at IS NULL")
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

func (r *UserRepositoryImpl) Create(tx *sql.Tx, user domain.User) (domain.User, error) {
	_, err := tx.Exec("INSERT INTO users (id,name,email,password,role) VALUES (?, ?,  ?, ?, ?)",
		user.Id, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return domain.User{}, err
	}

	createdUser, err := r.Get(tx, user.Id.String())
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (r *UserRepositoryImpl) Get(tx *sql.Tx, id string) (domain.User, error) {
	user := domain.User{}
	var phone sql.NullString
	err := tx.QueryRow("SELECT id, name, email, phone, role, created_at, updated_at FROM users WHERE id = ? AND deleted = false AND deleted_at IS NULL",
		id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

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

func (r *UserRepositoryImpl) Update(tx *sql.Tx, user domain.User) (domain.User, error) {
	// Ambil data user yang ada
	existingUser, err := r.Get(tx, user.Id.String())
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
	if user.Phone != nil {
		existingUser.Phone = user.Phone
	}
	if user.Role != "" {
		existingUser.Role = user.Role
	}

	// Eksekusi query update
	_, err = tx.Exec("UPDATE users SET name = ?, email = ?, password = ?, phone = ?, role = ? WHERE id = ?",
		existingUser.Name, existingUser.Email, existingUser.Password, existingUser.Phone, existingUser.Role, existingUser.Id)
	if err != nil {
		return domain.User{}, err
	}

	return r.Get(tx, user.Id.String())
}

func (r *UserRepositoryImpl) GetByEmail(tx *sql.Tx, email string) (domain.User, error) {
	user := domain.User{}
	err := tx.QueryRow("SELECT id,name,password,email,phone,role,created_at,updated_at FROM users WHERE email = ?", email).Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) Delete(tx *sql.Tx, id string) error {
	_, err := tx.Exec("UPDATE users SET deleted = ?, deleted_at = ? WHERE id = ?", true, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) Restore(tx *sql.Tx, id string) (domain.User, error) {
	_, err := tx.Exec("UPDATE users SET deleted = ?, deleted_at = ? WHERE id = ?", false, nil, id)
	if err != nil {
		return domain.User{}, err
	}
	user, _ := r.Get(tx, id)
	return user, nil
}

func (r *UserRepositoryImpl) GetDeletedUserById(tx *sql.Tx, id string) (domain.User, error) {
	users := domain.User{}
	err := tx.QueryRow("SELECT id,name,email,phone,role,created_at,updated_at FROM users WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?", id).Scan(&users)
	if err != nil {
		return domain.User{}, err
	}
	return users, nil
}
