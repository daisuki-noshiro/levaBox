package database

import (
	"database/sql"
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed schema.sql
var schemaSQL string

func Open() (*sql.DB, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dataDir := filepath.Join(configDir, "levaBox")

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "levabox.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	//第一版限制一个数据库连接
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if err := initDatabase(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func initDatabase(db *sql.DB) error {
	//外键检查
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	_, err = db.Exec(schemaSQL)
	if err != nil {
		return err
	}

	return nil
}
