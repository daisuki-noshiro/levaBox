package service

import (
	"errors"
	"reflect"
	"testing"

	"GalgameBox/metadata"
)

func TestCollectMetadataKeepsSourceOrderAndContinuesAfterFailure(t *testing.T) {
	sourceA := metadata.Source("source-a")
	sourceB := metadata.Source("source-b")
	unsupported := metadata.Source("unsupported")
	resolveCalls := make(map[metadata.Source]int)

	originalHandlers := metadataSourceHandlers
	metadataSourceHandlers = map[metadata.Source]metadataSourceHandler{
		sourceA: {
			Resolve: func(keyword string) (string, error) {
				resolveCalls[sourceA]++
				if keyword != "game" {
					t.Fatalf("Resolve keyword = %q, want %q", keyword, "game")
				}
				return "a-1", nil
			},
			Fetch: func(externalID string) (metadata.Result, error) {
				return metadata.Result{Source: sourceA, ExternalID: externalID}, nil
			},
		},
		sourceB: {
			Resolve: func(string) (string, error) {
				resolveCalls[sourceB]++
				return "", errors.New("temporary failure")
			},
			Fetch: func(string) (metadata.Result, error) {
				t.Fatal("Fetch should not run when Resolve fails")
				return metadata.Result{}, nil
			},
		},
	}
	t.Cleanup(func() {
		metadataSourceHandlers = originalHandlers
	})

	collection, err := CollectMetadata(
		"  game  ",
		[]metadata.Source{sourceA, sourceB, sourceA, unsupported},
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(collection.Results) != 1 || collection.Results[0].Source != sourceA {
		t.Fatalf("Results = %#v", collection.Results)
	}
	if len(collection.Issues) != 2 {
		t.Fatalf("len(Issues) = %d, want 2", len(collection.Issues))
	}
	if collection.Issues[0].Source != sourceB || collection.Issues[1].Source != unsupported {
		t.Fatalf("Issues = %#v", collection.Issues)
	}
	if resolveCalls[sourceA] != 1 || resolveCalls[sourceB] != 1 {
		t.Fatalf("resolveCalls = %v, want each source queried once", resolveCalls)
	}
}

func TestCollectMetadataRejectsEmptyKeyword(t *testing.T) {
	if _, err := CollectMetadata("  ", nil); err == nil {
		t.Fatal("CollectMetadata should reject an empty keyword")
	}
}

func TestDefaultMetadataSourcesReturnsCopy(t *testing.T) {
	first := DefaultMetadataSources()
	first[0] = metadata.Source("changed")

	want := []metadata.Source{metadata.SourceVNDB, metadata.SourceBangumi}
	if got := DefaultMetadataSources(); !reflect.DeepEqual(got, want) {
		t.Fatalf("DefaultMetadataSources() = %v, want %v", got, want)
	}
}
