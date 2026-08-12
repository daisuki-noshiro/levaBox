package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 官方vn查询入口
const vndbVNURL = "https://api.vndb.org/kana/vn"

type vndbSearchRequest struct {
	Filters []any  `json:"filters"`
	Fields  string `json:"fields"`
	Sort    string `json:"sort"`
	Results int    `json:"results"`
}

type vndbDeveloper struct {
	Name string `json:"name"`
}

type vndbImage struct {
	URL string `json:"url"`
}

type vndbSearchResult struct {
	ID         string          `json:"id"`
	Title      string          `json:"title"`
	AltTitle   string          `json:"alttitle"`
	Released   string          `json:"released"`
	Image      *vndbImage      `json:"image"`
	Developers []vndbDeveloper `json:"developers"`
}

type vndbSearchResponse struct {
	Results []vndbSearchResult `json:"results"`
	More    bool               `json:"more"`
}

var vndbHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

func SearchVNDB(keyword string) ([]vndbSearchResult, error) {
	requestData := vndbSearchRequest{
		Filters: []any{
			"search",
			"=",
			keyword,
		},
		Fields:  "title,alttitle,released,image.url,developers.name",
		Sort:    "searchrank",
		Results: 10,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, vndbVNURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := vndbHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"VNDB 请求失败，状态码: %d",
			resp.StatusCode,
		)
	}

	var result vndbSearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

func convertVNDBResult(result vndbSearchResult) MetadataCandidate {
	developers := make([]string, 0, len(result.Developers))

	for _, developer := range result.Developers {
		developers = append(
			developers,
			developer.Name,
		)
	}

	var coverURL string

	if result.Image != nil {
		coverURL = result.Image.URL
	}

	return MetadataCandidate{
		Source:     "vndb",
		SourceID:   result.ID,
		Title:      result.Title,
		AltTitle:   result.AltTitle,
		Released:   result.Released,
		Developers: developers,
		CoverURL:   coverURL,
	}
}
