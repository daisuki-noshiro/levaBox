package database

import (
	"database/sql"
	"testing"

	"GalgameBox/model"

	_ "modernc.org/sqlite"
)

func TestInsertGame(t *testing.T) {
	db, _ := sql.Open("sqlite", ":memory:")
	defer db.Close()

	db.SetMaxOpenConns(1)

	initDatabase(db)

	game := model.Game{
		ID:          "G001",
		Title:       "测试游戏",
		Company:     "测试公司",
		Year:        2026,
		Description: "这是测试数据",

		Favorite:         false,
		Progress:         model.ProgressNotStarted,
		TotalPlaySeconds: 0,
	}

	InsertGame(db, game)
	savedGame, _ := GetGameByID(db, "G001")
	if savedGame.Title != game.Title {
		t.Fatalf("Title mismatch")
	}
}
