package database

import (
	"database/sql"
	"path/filepath"
	"testing"

	"GalgameBox/model"

	_ "modernc.org/sqlite"
)

func openTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	db.SetMaxOpenConns(1)
	if err := initDatabase(db); err != nil {
		t.Fatal(err)
	}

	return db
}

func newTestGame(t *testing.T, id string) model.Game {
	t.Helper()

	return model.Game{
		ID:       id,
		Title:    "测试游戏",
		Progress: model.ProgressNotStarted,
		Launch: model.LaunchConfig{
			ExecutablePath:   filepath.Join(t.TempDir(), id+".exe"),
			WorkingDirectory: "",
		},
	}
}

func TestUpdateGameTitle(t *testing.T) {
	db := openTestDatabase(t)

	game := newTestGame(t, "G001")
	game.Title = "旧标题"

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

func TestGameYearRoundTrip(t *testing.T) {
	knownYear := 2018
	tests := []struct {
		name string
		year *int
	}{
		{name: "unknown year", year: nil},
		{name: "known year", year: &knownYear},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := openTestDatabase(t)
			game := newTestGame(t, "year-test")
			game.Year = test.year
			game.Company = ""
			game.Description = ""

			if err := InsertGame(db, game); err != nil {
				t.Fatal(err)
			}

			savedGame, err := GetGameByID(db, game.ID)
			if err != nil {
				t.Fatal(err)
			}

			if test.year == nil {
				if savedGame.Year != nil {
					t.Fatalf("Year = %v, want nil", savedGame.Year)
				}
			} else if savedGame.Year == nil || *savedGame.Year != *test.year {
				t.Fatalf("Year = %v, want %d", savedGame.Year, *test.year)
			}

			if savedGame.Company != "" || savedGame.Description != "" {
				t.Fatalf(
					"Company/Description = %q/%q, want empty strings",
					savedGame.Company,
					savedGame.Description,
				)
			}
		})
	}
}

func TestNullableGameColumnsReadAsModelEmptyValues(t *testing.T) {
	db := openTestDatabase(t)
	executablePath := filepath.Join(t.TempDir(), "nullable.exe")

	_, err := db.Exec(`
		INSERT INTO games (
			id,
			title,
			company,
			year,
			description,
			cover_path,
			background_type,
			background_path,
			bgm_path,
			executable_path,
			working_directory,
			last_played_at
		)
		VALUES (?, ?, NULL, NULL, NULL, NULL, NULL, NULL, NULL, ?, NULL, NULL);
	`, "nullable", "NULL 测试", executablePath)
	if err != nil {
		t.Fatal(err)
	}

	game, err := GetGameByID(db, "nullable")
	if err != nil {
		t.Fatal(err)
	}
	assertGameNullFields(t, game)

	games, err := ListGames(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("len(ListGames()) = %d, want 1", len(games))
	}
	assertGameNullFields(t, games[0])
}

func assertGameNullFields(t *testing.T, game model.Game) {
	t.Helper()

	if game.Company != "" {
		t.Fatalf("Company = %q, want empty", game.Company)
	}
	if game.Year != nil {
		t.Fatalf("Year = %v, want nil", game.Year)
	}
	if game.Description != "" {
		t.Fatalf("Description = %q, want empty", game.Description)
	}
	if game.CoverPath != "" {
		t.Fatalf("CoverPath = %q, want empty", game.CoverPath)
	}
	if game.Background.Type != "" || game.Background.Path != "" {
		t.Fatalf("Background = %#v, want empty", game.Background)
	}
	if game.BGMPath != "" {
		t.Fatalf("BGMPath = %q, want empty", game.BGMPath)
	}
	if game.Launch.WorkingDirectory != "" {
		t.Fatalf("WorkingDirectory = %q, want empty", game.Launch.WorkingDirectory)
	}
	if game.LastPlayedAt != nil {
		t.Fatalf("LastPlayedAt = %v, want nil", game.LastPlayedAt)
	}
}

func TestUpdateGameYearCanSetAndClear(t *testing.T) {
	db := openTestDatabase(t)
	game := newTestGame(t, "update-year")

	if err := InsertGame(db, game); err != nil {
		t.Fatal(err)
	}

	year := 2018
	if err := UpdateGameYear(db, game.ID, &year); err != nil {
		t.Fatal(err)
	}

	savedGame, err := GetGameByID(db, game.ID)
	if err != nil {
		t.Fatal(err)
	}
	if savedGame.Year == nil || *savedGame.Year != year {
		t.Fatalf("Year = %v, want %d", savedGame.Year, year)
	}

	if err := UpdateGameYear(db, game.ID, nil); err != nil {
		t.Fatal(err)
	}

	savedGame, err = GetGameByID(db, game.ID)
	if err != nil {
		t.Fatal(err)
	}
	if savedGame.Year != nil {
		t.Fatalf("Year = %v after clearing, want nil", savedGame.Year)
	}
}
