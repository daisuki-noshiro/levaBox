package service

import (
	"GalgameBox/database"
	"database/sql"
	"errors"
	"path/filepath"
)

var ErrGameAlreadyExists = errors.New("game already exists")

type ImportService struct {
	db *sql.DB
}

type StartImportResult struct {
	Status         string       `json:"status"`
	Draft          *ImportDraft `json:"draft,omitempty"`
	ExistingGameID string       `json:"existingGameId,omitempty"`
}

func NewImportService(db *sql.DB) *ImportService {
	return &ImportService{
		db: db,
	}
}

func (s *ImportService) StartImport(executablePath string) (StartImportResult, error) {
	existingGame, exists, err :=
		database.FindGameByExecutablePath(
			s.db,
			executablePath,
		)
	if err != nil {
		return StartImportResult{}, err
	}

	if exists {
		return StartImportResult{
			Status:         "already_exists",
			ExistingGameID: existingGame.ID,
		}, nil
	}

	absolutePath, err := filepath.Abs(executablePath)
	if err != nil {
		return StartImportResult{}, err
	}

	absolutePath = filepath.Clean(absolutePath)

	workingDirectory := filepath.Dir(absolutePath)
	searchKeyword := filepath.Base(workingDirectory)

	draft := ImportDraft{
		ExecutablePath:   absolutePath,
		WorkingDirectory: workingDirectory,
		SearchKeyword:    searchKeyword,
	}

	return StartImportResult{
		Status: "ready",
		Draft:  &draft,
	}, nil
}
