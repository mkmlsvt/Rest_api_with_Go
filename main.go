package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"myProject/internal/data"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	env  string
	dsn  string
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Printf("could not load env")
		os.Exit(1)
	}
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("FOOTBALLPLAYER_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// cfg.dsn = "postgresql://postgres:123456@localhost/goFootballPlayer?sslmode=disable"
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("database connection success")
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	addr := fmt.Sprintf(":%d", cfg.port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting %s server on %s", cfg.env, addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
