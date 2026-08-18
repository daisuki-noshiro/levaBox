package service

import (
	"reflect"
	"testing"

	"GalgameBox/metadata"
)

func TestBuildImportDraft(t *testing.T) {
	year := 2018
	vndbDescription := "VNDB description"
	bangumiDescription := "Bangumi description"

	results := []metadata.Result{
		{
			Source:            metadata.SourceVNDB,
			ExternalID:        "v20424",
			CompanyCandidates: []string{"Key"},
			Year:              &year,
			Description:       &vndbDescription,
			Covers: []metadata.ImageCandidate{
				{URL: "https://example.com/vndb-cover-1.jpg"},
				{URL: "https://example.com/vndb-cover-2.jpg"},
				{URL: "https://example.com/vndb-cover-3.jpg"},
			},
			Backgrounds: []metadata.ImageCandidate{
				{URL: "https://example.com/vndb-background-1.jpg"},
				{URL: "https://example.com/vndb-background-2.jpg"},
			},
		},
		{
			Source:            metadata.SourceBangumi,
			ExternalID:        "200763",
			CompanyCandidates: []string{"VISUAL ARTS", "Key"},
			Year:              &year,
			Description:       &bangumiDescription,
			Tags:              []string{"Galgame", "Key"},
			Covers: []metadata.ImageCandidate{
				{URL: "https://example.com/bangumi-cover-1.jpg"},
			},
		},
	}

	draft := BuildImportDraft(
		ImportDraft{
			ExecutablePath:   `C:\games\Summer Pockets\game.exe`,
			WorkingDirectory: `C:\games\Summer Pockets`,
			SearchKeyword:    "  Summer Pockets  ",
			Title:            "old title",
		},
		results,
	)

	if draft.SearchKeyword != "Summer Pockets" {
		t.Fatalf("SearchKeyword = %q, want %q", draft.SearchKeyword, "Summer Pockets")
	}
	if draft.Title != "Summer Pockets" {
		t.Fatalf("Title = %q, want %q", draft.Title, "Summer Pockets")
	}
	if draft.Company != "Key" {
		t.Fatalf("Company = %q, want %q", draft.Company, "Key")
	}
	if draft.Year == nil || *draft.Year != 2018 {
		t.Fatalf("Year = %v, want 2018", draft.Year)
	}
	if draft.Year == results[0].Year {
		t.Fatal("Year should be copied instead of sharing the metadata pointer")
	}
	if draft.Description == nil || *draft.Description != "Bangumi description" {
		t.Fatalf("Description = %v, want Bangumi description", draft.Description)
	}
	if !reflect.DeepEqual(draft.TagCandidates, []string{"Galgame", "Key"}) {
		t.Fatalf("TagCandidates = %v", draft.TagCandidates)
	}
	if len(draft.CoverCandidates) != 4 {
		t.Fatalf("len(CoverCandidates) = %d, want 4", len(draft.CoverCandidates))
	}
	if len(draft.BackgroundCandidates) != 2 {
		t.Fatalf("len(BackgroundCandidates) = %d, want 2", len(draft.BackgroundCandidates))
	}

	wantSources := []ResolvedSource{
		{Source: metadata.SourceVNDB, ExternalID: "v20424"},
		{Source: metadata.SourceBangumi, ExternalID: "200763"},
	}
	if !reflect.DeepEqual(draft.Sources, wantSources) {
		t.Fatalf("Sources = %#v, want %#v", draft.Sources, wantSources)
	}
}

func TestBuildImportDraftFallbacksAndDeduplication(t *testing.T) {
	t.Run("year falls back when VNDB has no year", func(t *testing.T) {
		fallbackYear := 2019
		draft := BuildImportDraft(ImportDraft{SearchKeyword: "game"}, []metadata.Result{
			{Source: metadata.SourceVNDB},
			{Source: metadata.SourceBangumi, Year: &fallbackYear},
		})

		if draft.Year == nil || *draft.Year != fallbackYear {
			t.Fatalf("Year = %v, want %d", draft.Year, fallbackYear)
		}
	})

	t.Run("description falls back when Bangumi has no description", func(t *testing.T) {
		fallbackDescription := "VNDB description"
		draft := BuildImportDraft(ImportDraft{SearchKeyword: "game"}, []metadata.Result{
			{Source: metadata.SourceVNDB, Description: &fallbackDescription},
			{Source: metadata.SourceBangumi},
		})

		if draft.Description == nil || *draft.Description != fallbackDescription {
			t.Fatalf("Description = %v, want %q", draft.Description, fallbackDescription)
		}
	})

	t.Run("company stays empty without a common candidate", func(t *testing.T) {
		draft := BuildImportDraft(ImportDraft{SearchKeyword: "game"}, []metadata.Result{
			{Source: metadata.SourceVNDB, CompanyCandidates: []string{"Key"}},
			{Source: metadata.SourceBangumi, CompanyCandidates: []string{"VisualArt's"}},
		})

		if draft.Company != "" {
			t.Fatalf("Company = %q, want empty", draft.Company)
		}
	})

	t.Run("company comparison ignores case and surrounding spaces", func(t *testing.T) {
		draft := BuildImportDraft(ImportDraft{SearchKeyword: "game"}, []metadata.Result{
			{CompanyCandidates: []string{"Key"}},
			{CompanyCandidates: []string{" key "}},
		})

		if draft.Company != "Key" {
			t.Fatalf("Company = %q, want %q", draft.Company, "Key")
		}
	})

	t.Run("tags are trimmed and deduplicated case-insensitively", func(t *testing.T) {
		draft := BuildImportDraft(ImportDraft{SearchKeyword: "game"}, []metadata.Result{
			{Tags: []string{" Galgame ", "Key", ""}},
			{Tags: []string{"galgame", " KEY ", "Drama"}},
		})

		want := []string{"Galgame", "Key", "Drama"}
		if !reflect.DeepEqual(draft.TagCandidates, want) {
			t.Fatalf("TagCandidates = %v, want %v", draft.TagCandidates, want)
		}
	})

	t.Run("image URLs are deduplicated", func(t *testing.T) {
		draft := BuildImportDraft(ImportDraft{SearchKeyword: "game"}, []metadata.Result{
			{
				Covers: []metadata.ImageCandidate{
					{URL: "https://example.com/cover.jpg"},
					{URL: ""},
				},
				Backgrounds: []metadata.ImageCandidate{
					{URL: "https://example.com/background.jpg"},
				},
			},
			{
				Covers: []metadata.ImageCandidate{
					{URL: "https://example.com/cover.jpg"},
					{URL: "https://example.com/cover-2.jpg"},
				},
				Backgrounds: []metadata.ImageCandidate{
					{URL: "https://example.com/background.jpg"},
				},
			},
		})

		if len(draft.CoverCandidates) != 2 {
			t.Fatalf("len(CoverCandidates) = %d, want 2", len(draft.CoverCandidates))
		}
		if len(draft.BackgroundCandidates) != 1 {
			t.Fatalf("len(BackgroundCandidates) = %d, want 1", len(draft.BackgroundCandidates))
		}
	})
}
