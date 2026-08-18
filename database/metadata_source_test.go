package database

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"GalgameBox/model"
)

func insertTransactionFixture(
	t *testing.T,
	db DBTX,
	game model.Game,
	tag model.Tag,
	source model.GameMetadataSource,
) {
	t.Helper()

	if err := InsertGame(db, game); err != nil {
		t.Fatal(err)
	}
	if err := InsertTag(db, tag); err != nil {
		t.Fatal(err)
	}
	if err := AddGameTag(db, game.ID, tag.ID); err != nil {
		t.Fatal(err)
	}
	if err := InsertGameMetadataSource(db, source); err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseHelpersCommitInTransaction(t *testing.T) {
	db := openTestDatabase(t)
	game := newTestGame(t, "commit-game")
	tag := model.Tag{ID: "commit-tag", Name: "Key"}
	source := model.GameMetadataSource{
		GameID:     game.ID,
		Source:     "vndb",
		ExternalID: "v20424",
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	insertTransactionFixture(t, tx, game, tag, source)

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	if _, err := GetGameByID(db, game.ID); err != nil {
		t.Fatalf("committed game was not found: %v", err)
	}

	savedTag, exists, err := GetTagByName(db, tag.Name)
	if err != nil {
		t.Fatal(err)
	}
	if !exists || savedTag != tag {
		t.Fatalf("saved tag = %#v, exists = %v", savedTag, exists)
	}

	tags, err := GetGameTags(db, game.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(tags, []model.Tag{tag}) {
		t.Fatalf("game tags = %#v, want %#v", tags, []model.Tag{tag})
	}

	sources, err := GetGameMetadataSources(db, game.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sources, []model.GameMetadataSource{source}) {
		t.Fatalf("metadata sources = %#v, want %#v", sources, []model.GameMetadataSource{source})
	}
}

func TestDatabaseHelpersRollbackTransaction(t *testing.T) {
	db := openTestDatabase(t)
	game := newTestGame(t, "rollback-game")
	tag := model.Tag{ID: "rollback-tag", Name: "Rollback Tag"}
	source := model.GameMetadataSource{
		GameID:     game.ID,
		Source:     "vndb",
		ExternalID: "v20424",
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	insertTransactionFixture(t, tx, game, tag, source)

	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}

	if _, err := GetGameByID(db, game.ID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetGameByID error = %v, want sql.ErrNoRows", err)
	}

	_, exists, err := GetTagByName(db, tag.Name)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("rolled back tag still exists")
	}

	tags, err := GetGameTags(db, game.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 0 {
		t.Fatalf("rolled back game_tags rows = %#v", tags)
	}

	sources, err := GetGameMetadataSources(db, game.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(sources) != 0 {
		t.Fatalf("rolled back metadata sources = %#v", sources)
	}
}

func TestGameMetadataSourceUniquePerGameAndSource(t *testing.T) {
	db := openTestDatabase(t)
	game := newTestGame(t, "unique-source-game")
	if err := InsertGame(db, game); err != nil {
		t.Fatal(err)
	}

	first := model.GameMetadataSource{
		GameID:     game.ID,
		Source:     "vndb",
		ExternalID: "v20424",
	}
	if err := InsertGameMetadataSource(db, first); err != nil {
		t.Fatal(err)
	}

	duplicate := first
	duplicate.ExternalID = "v99999"
	if err := InsertGameMetadataSource(db, duplicate); err == nil {
		t.Fatal("duplicate game/source pair should fail")
	}
}

func TestGameMetadataSourceAllowsSharedExternalID(t *testing.T) {
	db := openTestDatabase(t)
	game1 := newTestGame(t, "shared-source-game-1")
	game2 := newTestGame(t, "shared-source-game-2")

	if err := InsertGame(db, game1); err != nil {
		t.Fatal(err)
	}
	if err := InsertGame(db, game2); err != nil {
		t.Fatal(err)
	}

	for _, game := range []model.Game{game1, game2} {
		source := model.GameMetadataSource{
			GameID:     game.ID,
			Source:     "vndb",
			ExternalID: "v20424",
		}
		if err := InsertGameMetadataSource(db, source); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetGameMetadataSourcesOrdersBySource(t *testing.T) {
	db := openTestDatabase(t)
	game := newTestGame(t, "ordered-sources-game")
	if err := InsertGame(db, game); err != nil {
		t.Fatal(err)
	}

	sources := []model.GameMetadataSource{
		{GameID: game.ID, Source: "vndb", ExternalID: "v20424"},
		{GameID: game.ID, Source: "bangumi", ExternalID: "200763"},
	}
	for _, source := range sources {
		if err := InsertGameMetadataSource(db, source); err != nil {
			t.Fatal(err)
		}
	}

	got, err := GetGameMetadataSources(db, game.ID)
	if err != nil {
		t.Fatal(err)
	}
	want := []model.GameMetadataSource{sources[1], sources[0]}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetGameMetadataSources() = %#v, want %#v", got, want)
	}
}
