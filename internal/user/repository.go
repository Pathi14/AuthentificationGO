package user

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user User) error {
	fmt.Println("Attempting to create user:", user.Email)

	existingUser, err := r.GetByEmail(user.Email)
	if err != nil {
		fmt.Println("Error checking existing user:", err)
		return err
	}
	if existingUser != nil {
		fmt.Println("User already exists:", user.Email)
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	_, err = r.db.Exec(
		"INSERT INTO users (name, age, mobile_number, email, password) VALUES ($1, $2, $3, $4, $5)",
		user.Name, user.Age, user.MobileNumber, user.Email, user.Password)

	if err != nil {
		fmt.Println("Error inserting user:", err)
		return fmt.Errorf("error inserting user: %w", err)
	}

	fmt.Println("User created successfully:", user.Email)
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var u User
	err := r.db.QueryRow("SELECT id, name, email FROM users WHERE email = $1", email).
		Scan(&u.ID, &u.Name, &u.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Login(email, password string) (*User, error) {
	var u User
	var hashedPassword string
	err := r.db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", email).
		Scan(&u.ID, &u.Name, &u.Email, &hashedPassword)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid email or password")
	}
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	return &u, nil
}

func (r *UserRepository) FindByID(id int) (*User, error) {
	var user User
	query := "SELECT id, name, age, mobile_number, email FROM users WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Age, &user.MobileNumber, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ResetPassword(email, hashedPassword string) error {

	_, err := r.db.Exec("UPDATE users SET password = $1 WHERE email = $2", hashedPassword, email)

	return err

}
