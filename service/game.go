package service

import "database/sql"

type GameService struct {
	db *sql.DB
}

func NewGameService(db *sql.DB) *GameService {
	return &GameService{
		db: db,
	}
}
