package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const (
	version = "1.0.0"
)

type config struct {
	User 	 string
	Password string
	Database string
	Host     string
	Port     int
	Environment		 string
}


type application struct {
	config config
	logger *log.Logger
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg config

	cfg.User     = os.Getenv("POSTGRES_USER")
	cfg.Password = os.Getenv("POSTGRES_PW")
	cfg.Database = os.Getenv("DATABASE_SQL")
	cfg.Host 	 = os.Getenv("POSTGRES_HOST")
	cfg.Environment		 = os.Getenv("ENVIRONMENT")

	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
    if err != nil {
        panic(err)
    }
	cfg.Port     = port

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Environment, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

	app := &application{
		config: cfg,
		logger:logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	srv := &http.Server{
		Addr: 		  fmt.Sprintf(":%d", cfg.Port),
		Handler: 	  app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	ctx := context.Background()

	connString := fmt.Sprintf( `user=%s password=%s database=%s host=%s port=%d`,
								cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port)
	fmt.Print(connString)
	conn, err := pgx.Connect(ctx, connString)

	fmt.Print(conn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	// queries := db.New(conn)


	logger.Printf("starting %s server on %s", cfg.Environment, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}