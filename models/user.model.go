package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func CreateUser(user User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users(email, password, username) VALUES($1, $2, $3)`

	_, err = db.Exec(stmt, user.Email, string(hashedPassword), user.Username)

	return err
}

func GetUserById(id string) (User, error) {
	query := `SELECT * FROM users WHERE id=$1`

	stmt, err := db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	var user User
	err = stmt.QueryRow(id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CheckEmail(email string) (User, error) {
	query := `SELECT * FROM users WHERE email=$1`

	stmt, err := db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	var user User
	err = stmt.QueryRow(email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
