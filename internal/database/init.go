package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	DB    *sql.DB
	Users *Queries
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
	
		// create a queries instance
		Users = New(DB)
	})
}


