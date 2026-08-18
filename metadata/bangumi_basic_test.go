package metadata

import (
	"reflect"
	"testing"
)

func TestBuildBangumiCompaniesSplitsCombinedDevelopers(t *testing.T) {
	companies := buildBangumiCompanies([]bangumiInfobox{
		{Key: "开发商", Value: "VISUAL ARTS / Key"},
		{Key: "开发", Value: map[string]any{"v": "Key／Other"}},
	})

	want := []string{"VISUAL ARTS", "Key", "Other"}
	if !reflect.DeepEqual(companies, want) {
		t.Fatalf("buildBangumiCompanies() = %v, want %v", companies, want)
	}
}
