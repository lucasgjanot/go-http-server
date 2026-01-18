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
	GetChirps(ctx context.Context, desc bool) ([]Chirp, error)
	GetChirpsByUserId(ctx context.Context, args GetChirpsByUserIdParams) ([]Chirp, error)
	GetChirp(ctx context.Context, chirpId uuid.UUID) (Chirp, error)
	DeleteChirp(ctx context.Context, chirpId uuid.UUID) (Chirp, error)
}

type UsersInterface interface {
	CreateUser(ctx context.Context, args CreateUserParams) (User, error)
	UpdateUser(ctx context.Context, args UpdateUserParams) (User, error)
	UpgradeUser(ctx context.Context, userID uuid.UUID ) (User, error)
	DeleteAllUsers(ctx context.Context) ([]User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, userID uuid.UUID) (User, error)
}

type AuthorizationInterface interface {
	CreateRefreshToken(ctx context.Context, args CreateRefreshTokenParams) (RefreshToken, error)
	GetRefreshToken(ctx context.Context, token string) (RefreshToken, error)
	RevokeToken(ctx context.Context, token string) (RefreshToken, error)
}

var (
	DB    *sql.DB
	Chirps ChirpsInterface
	Users UsersInterface
	Auth AuthorizationInterface
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
		Auth = queries
	})
}


