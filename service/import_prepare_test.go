package service

import (
	"errors"
	"testing"

	"GalgameBox/metadata"
)

func TestPrepareImportMetadataDistinguishesNilAndEmptySources(t *testing.T) {
	callCount := 0
	originalHandlers := metadataSourceHandlers
	metadataSourceHandlers = map[metadata.Source]metadataSourceHandler{
		metadata.SourceVNDB: {
			Resolve: func(string) (string, error) {
				callCount++
				return "v1", nil
			},
			Fetch: func(externalID string) (metadata.Result, error) {
				return metadata.Result{Source: metadata.SourceVNDB, ExternalID: externalID}, nil
			},
		},
		metadata.SourceBangumi: {
			Resolve: func(string) (string, error) {
				callCount++
				return "1", nil
			},
			Fetch: func(externalID string) (metadata.Result, error) {
				return metadata.Result{Source: metadata.SourceBangumi, ExternalID: externalID}, nil
			},
		},
	}
	t.Cleanup(func() {
		metadataSourceHandlers = originalHandlers
	})

	service := &ImportService{}
	base := ImportDraft{SearchKeyword: "  game  ", Title: "old title"}

	withDefaults, err := service.PrepareImportMetadata(base, nil)
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 2 {
		t.Fatalf("default sources made %d calls, want 2", callCount)
	}
	if len(withDefaults.Draft.Sources) != 2 {
		t.Fatalf("len(Sources) = %d, want 2", len(withDefaults.Draft.Sources))
	}
	if withDefaults.Draft.Title != "game" || withDefaults.Draft.SearchKeyword != "game" {
		t.Fatalf("Draft keyword/title were not normalized: %#v", withDefaults.Draft)
	}

	withoutSources, err := service.PrepareImportMetadata(base, []metadata.Source{})
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 2 {
		t.Fatalf("empty sources unexpectedly made a request; call count = %d", callCount)
	}
	if len(withoutSources.Draft.Sources) != 0 {
		t.Fatalf("Sources = %#v, want none", withoutSources.Draft.Sources)
	}
}

func TestPrepareImportMetadataRejectsEmptyKeyword(t *testing.T) {
	service := &ImportService{}
	if _, err := service.PrepareImportMetadata(ImportDraft{SearchKeyword: "  "}, nil); err == nil {
		t.Fatal("PrepareImportMetadata should reject an empty keyword")
	}
}

func TestPrepareImportMetadataReturnsDraftWhenAllSourcesFail(t *testing.T) {
	originalHandlers := metadataSourceHandlers
	metadataSourceHandlers = map[metadata.Source]metadataSourceHandler{
		metadata.SourceVNDB: {
			Resolve: func(string) (string, error) {
				return "", errors.New("VNDB unavailable")
			},
		},
		metadata.SourceBangumi: {
			Resolve: func(string) (string, error) {
				return "", errors.New("Bangumi unavailable")
			},
		},
	}
	t.Cleanup(func() {
		metadataSourceHandlers = originalHandlers
	})

	result, err := (&ImportService{}).PrepareImportMetadata(
		ImportDraft{SearchKeyword: "  game  "},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.Draft.Title != "game" {
		t.Fatalf("Title = %q, want %q", result.Draft.Title, "game")
	}
	if len(result.Issues) != 2 {
		t.Fatalf("len(Issues) = %d, want 2", len(result.Issues))
	}
}
