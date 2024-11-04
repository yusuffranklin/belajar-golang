package models

import (
	"errors"
	"rest-api-practice/database"
	"rest-api-practice/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?, ?)"

	stmt, err := database.Db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPasswd, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPasswd)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId

	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := database.Db.QueryRow(query, u.Email)

	var retrievedPasswd string
	err := row.Scan(&u.ID, &retrievedPasswd)
	if err != nil {
		return errors.New("credentials invalid")
	}

	passwordValid := utils.CheckHashedPassword(u.Password, retrievedPasswd)
	if !passwordValid {
		return errors.New("credentials invalid")
	}

	return nil
}
