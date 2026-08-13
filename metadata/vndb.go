package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

var vndbHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// 查询游戏对应
func SearchVNDB(keyword string) ([]vndbAPIResult, error) {
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

	var result vndbSearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}
