package database

import (
	"database/sql"
	"testing"

	"GalgameBox/model"

	_ "modernc.org/sqlite"
)

func TestUpdateGameTitle(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	defer db.Close()

	db.SetMaxOpenConns(1)
	if err := initDatabase(db); err != nil {
		t.Fatal(err)
	}

	game := model.Game{
		ID:       "G001",
		Title:    "旧标题",
		Progress: model.ProgressNotStarted,
	}

	if err := InsertGame(db, game); err != nil {
		t.Fatal(err)
	}

	if err := UpdateGameTitle(db, "G001", "新标题"); err != nil {
		t.Fatal(err)
	}

	savedGame, err := GetGameByID(db, "G001")
	if err != nil {
		t.Fatal(err)
	}

	if savedGame.Title != "新标题" {
		t.Fatalf("标题修改失败")
	}
}
