package database

import (
	"context"
	"database/sql"
	"log"
	"sync"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type ChirpsInterface interface {
	CreateChirp(ctx context.Context, args CreateChirpParams) (Chirp, error)
	GetChirps(ctx context.Context) ([]Chirp, error)
	GetChirp(ctx context.Context, chirpId uuid.UUID) (Chirp, error)
}

type UsersInterface interface {
	CreateUser(ctx context.Context, args CreateUserParams) (User, error)
	DeleteAllUsers(ctx context.Context) ([]User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

var (
	DB    *sql.DB
	Chirps ChirpsInterface
	Users UsersInterface
	once sync.Once

)

// Init initializes the database without crashing the server
func Init(dsn string) {
	once.Do(func() {
		var err error
		DB, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
	
		// optional: ping to check connection
		if err := DB.Ping(); err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}
		
		queries := New(DB)
		// create a queries instance
		Users = queries
		Chirps = queries
	})
}


