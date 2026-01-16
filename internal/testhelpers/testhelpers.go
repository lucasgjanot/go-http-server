package testhelpers

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/joho/godotenv"
	"github.com/lucasgjanot/go-http-server/internal/auth"
	"github.com/lucasgjanot/go-http-server/internal/config"
	"github.com/lucasgjanot/go-http-server/internal/database"
	"github.com/lucasgjanot/go-http-server/internal/router"
)

var once sync.Once

type NewUser struct {
	Email string
	Password string
}

func InitTest(t *testing.T) string {
	cfg := config.NewConfig()
	r := router.New(cfg)

	initTestDB(t)

	ts := httptest.NewServer(r.Handler)
	t.Cleanup(ts.Close)

	if err := ResetDatabase(); err != nil {
		t.Fatalf("Error resetting database: %s", err)
	}
	t.Cleanup(func() {
		_ = ResetDatabase()
	})
	return ts.URL
}



func projectRoot() string {
	var projectRootDir string
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	dir := wd
	for {
		// Use go.mod as the marker; you can add more markers if needed.
		marker := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(marker); err == nil {
			projectRootDir = dir
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding marker
			log.Fatalf("project root not found from %q (no go.mod up the tree)", wd)
		}
		dir = parent
	}
    return projectRootDir
}


// InitTestDB loads the environment and initializes the database exactly once
func initTestDB(t *testing.T) {
	t.Helper()

	once.Do(func() {
		envPath := filepath.Join(projectRoot(), ".env.development")
		// Load .env.development
		if err := godotenv.Load(envPath); err != nil {
			log.Fatalf("failed to load .env.development: %v", err)
		}

		// Initialize the database using DATABASE_URL from env
		dsn := os.Getenv("DB_URL")
		database.Init(dsn)
	})
}

func CreateUser(args NewUser) (database.User, error) {
	if args.Email == "" {
		args.Email = "test@exemple.com"
	}
	if args.Password == "" {
		args.Password = "secure"
	}
	hashedPassword, err := auth.HashPassword(args.Password)
	if err != nil {
		return database.User{}, err
	}
	return database.Users.CreateUser(context.Background(), database.CreateUserParams{
		Email: args.Email,
		HashedPassword: hashedPassword,
	})
}

func ResetDatabase() error {
	_, err := database.Users.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func CreateChirp(args database.CreateChirpParams) (database.Chirp, error) {
	return database.Chirps.CreateChirp(
		context.Background(),
		args,
	)
}