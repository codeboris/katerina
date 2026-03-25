package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	dsn := flag.String("dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	dir := flag.String("dir", "migrations", "Migrations directory")
	flag.Parse()

	cmd := "up"
	if flag.NArg() > 0 {
		cmd = flag.Arg(0)
	}

	if *dsn == "" {
		log.Fatal("DB_DSN env var or -dsn flag is required")
	}

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("set dialect: %v", err)
	}

	switch cmd {
	case "up", "down", "reset", "status", "version":
		if err := goose.Run(cmd, db, *dir); err != nil {
			log.Fatalf("goose %s: %v", cmd, err)
		}
	case "create":
		if flag.NArg() < 2 {
			log.Fatal("usage: migrate create <name>")
		}
		name := flag.Arg(1)
		if err := goose.Create(db, *dir, name, "sql"); err != nil {
			log.Fatalf("goose create: %v", err)
		}
		fmt.Printf("Created migration: %s\n", name)
	default:
		log.Fatalf("unknown command: %s (up|down|reset|status|version|create)", cmd)
	}
}
