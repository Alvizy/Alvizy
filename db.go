package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dsn string) error {
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	return db.Ping()
}

func MustInitDB(dsn string) {
	if err := InitDB(dsn); err != nil {
		log.Fatalf("gagal koneksi database: %v", err)
	}
	log.Println("database terhubung")
}

func RunMigrations() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS profile (
			id         SERIAL PRIMARY KEY,
			name       TEXT NOT NULL DEFAULT 'Fanz',
			tagline    TEXT NOT NULL DEFAULT '',
			bio        TEXT NOT NULL DEFAULT '',
			location   TEXT NOT NULL DEFAULT '',
			email      TEXT NOT NULL DEFAULT '',
			github     TEXT NOT NULL DEFAULT '',
			instagram  TEXT NOT NULL DEFAULT '',
			linkedin   TEXT NOT NULL DEFAULT '',
			available  BOOLEAN NOT NULL DEFAULT true,
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		`INSERT INTO profile (id) VALUES (1) ON CONFLICT DO NOTHING`,
		`CREATE TABLE IF NOT EXISTS skills (
			id          SERIAL PRIMARY KEY,
			name        TEXT NOT NULL,
			category    TEXT NOT NULL DEFAULT 'Other',
			order_index INT NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id          SERIAL PRIMARY KEY,
			title       TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			tech_stack  TEXT NOT NULL DEFAULT '',
			github_url  TEXT NOT NULL DEFAULT '',
			live_url    TEXT NOT NULL DEFAULT '',
			featured    BOOLEAN NOT NULL DEFAULT false,
			order_index INT NOT NULL DEFAULT 0,
			created_at  TIMESTAMP DEFAULT NOW()
		)`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
