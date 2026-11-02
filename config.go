package main

import (
	"encoding/json"
	"log"
	"os"
)

type AuthConfig struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type DatabaseConfig struct {
	DSN string `json:"dsn"`
}

type Config struct {
	Port      string         `json:"port"`
	AdminPath string         `json:"admin_path"`
	Auth      AuthConfig     `json:"auth"`
	Database  DatabaseConfig `json:"database"`
}

var cfg Config

func LoadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&cfg)
}

func SaveConfig(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(&cfg)
}

func MustLoadConfig(path string) {
	if err := LoadConfig(path); err != nil {
		log.Fatalf("gagal load config: %v", err)
	}
}
