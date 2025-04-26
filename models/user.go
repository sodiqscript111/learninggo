package models

import (
	"errors"
	"fmt"
	"learninggo/db"
	"learninggo/utils"
	"log"
)

type User struct {
	Id       int
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

var users []User = []User{}

func (u *User) InsertUser() {
	query := `INSERT INTO users(email, password) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Fatalf("❌ Failed to prepare statement: %v", err)
	}
	defer stmt.Close()
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {

	}
	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		log.Fatalf("❌ Failed to insert user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("❌ Failed to get inserted user ID: %v", err)
	}
	u.Id = int(id)
	log.Printf("✅ User inserted successfully with ID %d\n", u.Id)
}

func (u *User) ValidateUser() error {
	query := `SELECT password, id FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&retrievedPassword, &u.Id)
	if err != nil {
		fmt.Println("DB Error:", err)
		return errors.New("Invalid email or password")
	}

	fmt.Println("Retrieved Hash:", retrievedPassword)
	fmt.Println("Entered Password:", u.Password)
	fmt.Println("User ID:", u.Id)

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	fmt.Println("Password Match Result:", passwordIsValid)

	if !passwordIsValid {
		return errors.New("Invalid email or password")
	}

	return nil
}
