package repositories

import (
	"database/sql"
	"errors"
	"sterling-hms-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(firstName, lastName, email, hashedPassword string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		`INSERT INTO users (first_name, last_name, email, password, created_at, updated_at, is_active)
		 VALUES ($1, $2, $3, $4, NOW(), NOW(), true)
		 RETURNING id, first_name, last_name, email, created_at, updated_at, is_active`,
		firstName, lastName, email, hashedPassword,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, first_name, last_name, email, password, created_at, updated_at, last_login, is_active
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, first_name, last_name, email, created_at, updated_at, last_login, is_active
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateLastLogin(id int) error {
	_, err := r.db.Exec(
		`UPDATE users SET last_login = NOW() WHERE id = $1`,
		id,
	)
	return err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = $1`, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) GetAll() ([]*models.User, error) {
	rows, err := r.db.Query(
		`SELECT id, first_name, last_name, email, created_at, updated_at, last_login, is_active
		 FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
