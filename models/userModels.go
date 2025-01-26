package models

import (
	"errors"
	"meeting-backend/config"
	"meeting-backend/token"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   uint   `json:"role_id"`
}

func VerifyPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(email string, password string) (string, error) {
	var err error
	var u User

	query := "SELECT s.Staff_ID, s.Email, u.User_Password, u.Role_ID FROM users u JOIN staff s ON u.Staff_ID = s.Staff_ID WHERE s.Email = ?"
	res, err := config.DB.Query(query, email)

	if err != nil {
		return "", err
	}

	if res.Next() {
		if err := res.Scan(&u.ID, &u.RoleID, &u.Email, &u.Password); err != nil {
			return "", err
		}
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	genToken, err := token.GenerateToken(u.ID, u.RoleID)

	if err != nil {
		return "", err
	}

	return genToken, nil
}

func GetUserByID(uid uint) (User, error) {
	var u User

	query := "SELECT s.Staff_ID, s.Email, u.User_Password FROM users u JOIN staff s ON u.Staff_ID = s.Staff_ID WHERE u.User_ID = ?"

	res, err := config.DB.Query(query, uid)

	if err != nil {
		return u, errors.New("internal server error")
	}

	if res.Next() {
		if err := res.Scan(&u.ID, &u.Email, &u.Password); err != nil {
			return u, errors.New("internal server error")
		}
	}

	u.Password = ""

	return u, nil
}
