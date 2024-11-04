package main

import (
	"audio-go/internal/auth"
	"audio-go/internal/db"
	"audio-go/internal/env"
	"audio-go/internal/store"
	// "time"
	"go.uber.org/zap"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

func main() {
	addr:= "postgres://postgres:niki77@localhost:9595/postgres?sslmode=disable"

	cfg := config{
		addr: env.GetString("ADDR", ":9595"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", addr),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),

		// auth: authConfig{
		// 	basic: basicConfig{
		// 		user: env.GetString("AUTH_BASIC_USER", "admin"),
		// 		pass: env.GetString("AUTH_BASIC_PASS", "admin"),
		// 	},
		// 	token: tokenConfig{
		// 		secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
		// 		exp:    time.Hour * 24 * 3, // 3 days
		// 		iss:    "audio",
		// 	},
		// },
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	//Auth
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	store := store.NewStorage(db)

	app := &application{
		config:        cfg,
		store:         store,
		authenticator: jwtAuthenticator,
		logger: logger,
	}
	mux := app.mount()

	err = app.run(mux)

	if err != nil {
		logger.Fatal(err)
	}
}
