package service

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"GalgameBox/database"
	"GalgameBox/metadata"
	"GalgameBox/model"
)

func newSaveImportTestService(t *testing.T) (*ImportService, *sql.DB, string) {
	t.Helper()

	configRoot := t.TempDir()
	t.Setenv("AppData", configRoot)
	t.Setenv("XDG_CONFIG_HOME", configRoot)
	t.Setenv("HOME", configRoot)

	db, err := database.Open()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	mediaRoot := filepath.Join(t.TempDir(), "media")
	return &ImportService{
		db:        db,
		mediaRoot: mediaRoot,
	}, db, mediaRoot
}

func newSaveImportRequest(t *testing.T) SaveImportRequest {
	t.Helper()

	workingDirectory := t.TempDir()
	return SaveImportRequest{
		ExecutablePath:   filepath.Join(workingDirectory, "game.exe"),
		WorkingDirectory: workingDirectory,
		Title:            "测试游戏",
	}
}

func TestSaveImportCompleteSuccess(t *testing.T) {
	service, db, mediaRoot := newSaveImportTestService(t)

	coverBytes := []byte("cover bytes")
	backgroundBytes := []byte("background bytes")
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/cover.jpg":
			response.Header().Set("Content-Type", "image/jpeg")
			_, _ = response.Write(coverBytes)
		case "/background":
			response.Header().Set("Content-Type", "image/png")
			_, _ = response.Write(backgroundBytes)
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	existingTag := model.Tag{ID: "existing-key-tag", Name: "Key"}
	if err := database.InsertTag(db, existingTag); err != nil {
		t.Fatal(err)
	}

	year := 2018
	description := "  Bangumi description  "
	request := newSaveImportRequest(t)
	request.Title = "  Summer Pockets  "
	request.Company = "  Key  "
	request.Year = &year
	request.Description = &description
	request.Tags = []string{" Key ", "key", "", "恋爱"}
	request.Sources = []ResolvedSource{
		{Source: metadata.SourceVNDB, ExternalID: " v20424 "},
		{Source: metadata.SourceVNDB, ExternalID: "v99999"},
		{Source: "", ExternalID: "ignored"},
		{Source: metadata.SourceBangumi, ExternalID: " 200763 "},
	}
	request.Cover = &metadata.ImageCandidate{
		URL: server.URL + "/cover.jpg?v=1",
	}
	request.Background = &metadata.ImageCandidate{
		URL: server.URL + "/background",
	}

	game, err := service.SaveImport(request)
	if err != nil {
		t.Fatal(err)
	}

	if game.ID == "" {
		t.Fatal("Game.ID should not be empty")
	}
	if game.Title != "Summer Pockets" || game.Company != "Key" {
		t.Fatalf("Title/Company = %q/%q", game.Title, game.Company)
	}
	if game.Description != "Bangumi description" {
		t.Fatalf("Description = %q", game.Description)
	}
	if game.Year == nil || *game.Year != year {
		t.Fatalf("Year = %v, want %d", game.Year, year)
	}
	if game.Progress != model.ProgressNotStarted || game.Favorite {
		t.Fatalf("Progress/Favorite = %q/%v", game.Progress, game.Favorite)
	}
	if len(game.TagIDs) != 2 || game.TagIDs[0] != existingTag.ID {
		t.Fatalf("TagIDs = %#v", game.TagIDs)
	}

	mediaDirectory := filepath.Join(mediaRoot, game.ID)
	if filepath.Dir(game.CoverPath) != mediaDirectory {
		t.Fatalf("CoverPath = %q, want directory %q", game.CoverPath, mediaDirectory)
	}
	if filepath.Base(game.CoverPath) != "cover.jpg" {
		t.Fatalf("CoverPath filename = %q, want cover.jpg", filepath.Base(game.CoverPath))
	}
	if game.Background.Type != model.BackgroundImage {
		t.Fatalf("Background.Type = %q", game.Background.Type)
	}
	if filepath.Base(game.Background.Path) != "background.png" {
		t.Fatalf("Background filename = %q, want background.png", filepath.Base(game.Background.Path))
	}

	gotCover, err := os.ReadFile(game.CoverPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotCover, coverBytes) {
		t.Fatalf("cover contents = %q", gotCover)
	}
	gotBackground, err := os.ReadFile(game.Background.Path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotBackground, backgroundBytes) {
		t.Fatalf("background contents = %q", gotBackground)
	}

	savedGame, err := database.GetGameByID(db, game.ID)
	if err != nil {
		t.Fatal(err)
	}
	if savedGame.Year == nil || *savedGame.Year != year {
		t.Fatalf("saved Year = %v, want %d", savedGame.Year, year)
	}
	if savedGame.CoverPath != game.CoverPath || savedGame.Background.Path != game.Background.Path {
		t.Fatalf("saved media paths do not match returned game")
	}

	tags, err := database.GetGameTags(db, game.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 {
		t.Fatalf("len(tags) = %d, want 2", len(tags))
	}
	tagNames := map[string]bool{tags[0].Name: true, tags[1].Name: true}
	if !tagNames["Key"] || !tagNames["恋爱"] {
		t.Fatalf("tags = %#v", tags)
	}

	sources, err := database.GetGameMetadataSources(db, game.ID)
	if err != nil {
		t.Fatal(err)
	}
	wantSources := []model.GameMetadataSource{
		{GameID: game.ID, Source: "bangumi", ExternalID: "200763"},
		{GameID: game.ID, Source: "vndb", ExternalID: "v20424"},
	}
	if !reflect.DeepEqual(sources, wantSources) {
		t.Fatalf("sources = %#v, want %#v", sources, wantSources)
	}
}

func TestSaveImportWithoutMedia(t *testing.T) {
	service, _, mediaRoot := newSaveImportTestService(t)
	request := newSaveImportRequest(t)

	game, err := service.SaveImport(request)
	if err != nil {
		t.Fatal(err)
	}
	if game.CoverPath != "" {
		t.Fatalf("CoverPath = %q, want empty", game.CoverPath)
	}
	if game.Background.Type != "" || game.Background.Path != "" {
		t.Fatalf("Background = %#v, want empty", game.Background)
	}
	if game.Year != nil || game.Description != "" {
		t.Fatalf("Year/Description = %v/%q", game.Year, game.Description)
	}

	if _, err := os.Stat(filepath.Join(mediaRoot, game.ID)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("empty media directory should not exist; error = %v", err)
	}
}

func TestSaveImportRejectsDuplicateExecutable(t *testing.T) {
	service, db, _ := newSaveImportTestService(t)
	request := newSaveImportRequest(t)

	if _, err := service.SaveImport(request); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SaveImport(request); !errors.Is(err, ErrGameAlreadyExists) {
		t.Fatalf("second SaveImport error = %v, want ErrGameAlreadyExists", err)
	}

	games, err := database.ListGames(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("len(games) = %d, want 1", len(games))
	}
}

func TestSaveImportCleansMediaAfterDownloadFailure(t *testing.T) {
	service, db, mediaRoot := newSaveImportTestService(t)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/cover.jpg" {
			response.Header().Set("Content-Type", "image/jpeg")
			_, _ = response.Write([]byte("cover"))
			return
		}

		http.Error(response, "failed", http.StatusInternalServerError)
	}))
	defer server.Close()

	request := newSaveImportRequest(t)
	request.Cover = &metadata.ImageCandidate{URL: server.URL + "/cover.jpg"}
	request.Background = &metadata.ImageCandidate{URL: server.URL + "/background.jpg"}

	if _, err := service.SaveImport(request); err == nil {
		t.Fatal("SaveImport should fail when the background download fails")
	}

	games, err := database.ListGames(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 0 {
		t.Fatalf("len(games) = %d, want 0", len(games))
	}
	assertDirectoryEmptyOrMissing(t, mediaRoot)
}

func TestSaveImportRollsBackDatabaseAndCleansMedia(t *testing.T) {
	service, db, mediaRoot := newSaveImportTestService(t)

	_, err := db.Exec(`
		CREATE TRIGGER fail_import_metadata_source
		BEFORE INSERT ON game_metadata_sources
		BEGIN
			SELECT RAISE(ABORT, 'forced metadata source failure');
		END;
	`)
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "image/jpeg")
		_, _ = response.Write([]byte("cover"))
	}))
	defer server.Close()

	request := newSaveImportRequest(t)
	request.Tags = []string{"Key", "恋爱"}
	request.Sources = []ResolvedSource{
		{Source: metadata.SourceVNDB, ExternalID: "v20424"},
	}
	request.Cover = &metadata.ImageCandidate{URL: server.URL + "/cover.jpg"}

	if _, err := service.SaveImport(request); err == nil {
		t.Fatal("SaveImport should fail on the injected metadata source error")
	}

	games, err := database.ListGames(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 0 {
		t.Fatalf("rolled back games = %#v", games)
	}

	tags, err := database.ListTags(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 0 {
		t.Fatalf("rolled back tags = %#v", tags)
	}

	assertTableRowCount(t, db, "game_tags", 0)
	assertTableRowCount(t, db, "game_metadata_sources", 0)
	assertDirectoryEmptyOrMissing(t, mediaRoot)
}

func TestSaveImportValidatesRequiredFieldsAndYear(t *testing.T) {
	service, db, _ := newSaveImportTestService(t)

	tests := []struct {
		name    string
		prepare func(*SaveImportRequest)
	}{
		{
			name: "empty executable path",
			prepare: func(request *SaveImportRequest) {
				request.ExecutablePath = "  "
			},
		},
		{
			name: "empty title",
			prepare: func(request *SaveImportRequest) {
				request.Title = "  "
			},
		},
		{
			name: "invalid year",
			prepare: func(request *SaveImportRequest) {
				year := 0
				request.Year = &year
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := newSaveImportRequest(t)
			test.prepare(&request)
			if _, err := service.SaveImport(request); err == nil {
				t.Fatal("SaveImport should reject invalid input")
			}
		})
	}

	games, err := database.ListGames(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 0 {
		t.Fatalf("invalid requests saved games: %#v", games)
	}
}

func assertDirectoryEmptyOrMissing(t *testing.T, directory string) {
	t.Helper()

	entries, err := os.ReadDir(directory)
	if errors.Is(err, os.ErrNotExist) {
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("directory %q is not empty: %#v", directory, entries)
	}
}

func assertTableRowCount(t *testing.T, db *sql.DB, table string, want int) {
	t.Helper()

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != want {
		t.Fatalf("%s row count = %d, want %d", table, count, want)
	}
}
