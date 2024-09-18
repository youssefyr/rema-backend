package Userc

import (
	"context"
	"fmt"
	"log"
	"xira/db"
	"xira/dbinit"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckEmail(ctx context.Context, email string, client db.PrismaClient) (bool, error) {

	exists, err := client.User.FindFirst(db.User.Email.Equals(email)).Exec(ctx)
	if err != nil {
		return false, err
	}
	return exists != nil, nil
}

func CreateUser(ctx context.Context, email, name, password string, rememberMe []string) (*db.UserModel, error) {

	dbinit.DbInit()
	client := dbinit.Client

	emailExists, err := CheckEmail(ctx, email, *client)
	if emailExists {
		return nil, fmt.Errorf("email already in use")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// Generate a new session ID
	sessionID := uuid.New().String()

	// Create the user
	user, err := client.User.CreateOne(
		db.User.Email.Set(email),
		db.User.Name.Set(name),
		db.User.Password.Set(string(hashedPassword)),
		db.User.SessionID.Set(sessionID),
		db.User.RememberMe.Set(rememberMe),
	).Exec(ctx)
	if err != nil {
		// Log the error and continue execution
		log.Println("Error creating user:", err)
	}

	return user, err
}
