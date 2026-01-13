package testhelpers

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/joho/godotenv"
	"github.com/lucasgjanot/go-http-server/internal/database"
)

var once sync.Once

// InitTestDB loads the environment and initializes the database exactly once
func InitTestDB(t *testing.T) {
	t.Helper()

	once.Do(func() {
		// Load .env.development
		if err := godotenv.Load("../../../../.env.development"); err != nil {
			log.Fatalf("failed to load .env.development: %v", err)
		}

		// Initialize the database using DATABASE_URL from env
		dsn := os.Getenv("DB_URL")
		database.Init(dsn)
	})
}

func CreateUser(email string) (database.User, error) {
	return database.Users.CreateUser(context.Background(), email)
}

func ResetDatabase() error {
	_, err := database.Users.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}