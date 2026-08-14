package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// 官方 VN 查询入口
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
	Thumbnail string `json:"thumbnail"`
}

// VNDB API 返回的原始搜索结果。
// 只在 metadata 包内部使用。
type vndbAPIResult struct {
	ID         string          `json:"id"`
	Title      string          `json:"title"`
	AltTitle   string          `json:"alttitle"`
	Released   string          `json:"released"`
	Image      *vndbImage      `json:"image"`
	Developers []vndbDeveloper `json:"developers"`
}

type vndbSearchResponse struct {
	Results []vndbAPIResult `json:"results"`
	More    bool            `json:"more"`
}

// VNDBSearchResult 是 levaBox 使用的 VNDB 搜索结果。
// 作用只有一个：帮助用户判断“哪个条目才是自己的游戏”。
type VNDBSearchResult struct {
	ID        string
	Title     string
	AltTitle  string
	Released  string
	Companies []string
	Thumbnail string
}

var vndbHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// SearchVNDB 根据游戏名称搜索 VNDB 条目。
func SearchVNDB(keyword string) ([]VNDBSearchResult, error) {
	requestData := vndbSearchRequest{
		Filters: []any{
			"search",
			"=",
			keyword,
		},
		Fields:  "title,alttitle,released,image.thumbnail,developers.name",
		Sort:    "searchrank",
		Results: 10,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		vndbVNURL,
		bytes.NewReader(body),
	)
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

	var response vndbSearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	results := make([]VNDBSearchResult, 0, len(response.Results))

	for _, item := range response.Results {
		results = append(results, buildVNDBSearchResult(item))
	}

	return results, nil
}

func buildVNDBSearchResult(result vndbAPIResult) VNDBSearchResult {
	thumbnail := ""

	if result.Image != nil {
		thumbnail = result.Image.Thumbnail
	}

	return VNDBSearchResult{
		ID:        result.ID,
		Title:     strings.TrimSpace(result.Title),
		AltTitle:  strings.TrimSpace(result.AltTitle),
		Released:  result.Released,
		Companies: buildDeveloperNames(result.Developers),
		Thumbnail: thumbnail,
	}
}

// buildDeveloperNames 把 VNDB developers 整理成 levaBox 使用的名称列表。
func buildDeveloperNames(developers []vndbDeveloper) []string {
	companies := make([]string, 0, len(developers))
	seen := make(map[string]bool)

	for _, developer := range developers {
		name := strings.TrimSpace(developer.Name)

		if name == "" || seen[name] {
			continue
		}

		seen[name] = true
		companies = append(companies, name)
	}

	return companies
}
