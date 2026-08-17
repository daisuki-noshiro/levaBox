package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	bangumiSearchURL  = "https://api.bgm.tv/v0/search/subjects"
	bangumiSubjectURL = "https://api.bgm.tv/v0/subjects"
)

type bangumiSearchRequest struct {
	Keyword string              `json:"keyword"`
	Sort    string              `json:"sort"`
	Filter  bangumiSearchFilter `json:"filter"`
}

type bangumiSearchFilter struct {
	Type []int `json:"type"`
}

type bangumiSearchItem struct {
	ID     int            `json:"id"`
	Name   string         `json:"name"`
	NameCN string         `json:"name_cn"`
	Date   string         `json:"date"`
	Images *bangumiImages `json:"images"`
}

type bangumiSearchResponse struct {
	Data []bangumiSearchItem `json:"data"`
}

// BangumiSearchResult 表示按名称搜索得到的 Bangumi 条目。
type BangumiSearchResult struct {
	ID       string
	Name     string
	NameCN   string
	Date     *string
	ImageURL *string
}

type bangumiSubject struct {
	ID      int              `json:"id"`
	Date    string           `json:"date"`
	Summary string           `json:"summary"`
	Infobox []bangumiInfobox `json:"infobox"`
	Tags    []bangumiTag     `json:"tags"`
	Images  *bangumiImages   `json:"images"`
}

var bangumiHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// SearchBangumi 根据游戏名称搜索 Bangumi 游戏条目。
func SearchBangumi(keyword string) ([]BangumiSearchResult, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, fmt.Errorf("Bangumi 搜索关键词不能为空")
	}

	requestData := bangumiSearchRequest{
		Keyword: keyword,
		Sort:    "match",
		Filter: bangumiSearchFilter{
			Type: []int{4},
		},
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		bangumiSearchURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "daisuki-noshiro/levaBox")

	resp, err := bangumiHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"Bangumi 请求失败，状态码: %d",
			resp.StatusCode,
		)
	}

	var response bangumiSearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	results := make([]BangumiSearchResult, 0, len(response.Data))

	for _, item := range response.Data {
		results = append(results, buildBangumiSearchResult(item))
	}

	return results, nil
}

func buildBangumiSearchResult(item bangumiSearchItem) BangumiSearchResult {
	var date *string
	if strings.TrimSpace(item.Date) != "" {
		value := strings.TrimSpace(item.Date)
		date = &value
	}

	var imageURL *string
	if item.Images != nil && strings.TrimSpace(item.Images.Common) != "" {
		value := strings.TrimSpace(item.Images.Common)
		imageURL = &value
	}

	return BangumiSearchResult{
		ID:       strconv.Itoa(item.ID),
		Name:     strings.TrimSpace(item.Name),
		NameCN:   strings.TrimSpace(item.NameCN),
		Date:     date,
		ImageURL: imageURL,
	}
}

// ResolveBangumiID 根据搜索结果确定默认 Bangumi 条目。
func ResolveBangumiID(keyword string) (string, error) {
	results, err := SearchBangumi(keyword)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", fmt.Errorf("未找到 Bangumi 条目: %s", keyword)
	}

	return results[0].ID, nil
}

// getBangumiSubject 获取指定 Bangumi 条目的完整原始数据。
func getBangumiSubject(subjectID string) (bangumiSubject, error) {
	subjectID = strings.TrimSpace(subjectID)
	if subjectID == "" {
		return bangumiSubject{}, fmt.Errorf("Bangumi ID 不能为空")
	}

	url := fmt.Sprintf(
		"%s/%s",
		bangumiSubjectURL,
		subjectID,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return bangumiSubject{}, err
	}

	req.Header.Set("User-Agent", "daisuki-noshiro/levaBox")

	resp, err := bangumiHTTPClient.Do(req)
	if err != nil {
		return bangumiSubject{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return bangumiSubject{}, fmt.Errorf(
			"Bangumi 请求失败，状态码: %d",
			resp.StatusCode,
		)
	}

	var subject bangumiSubject

	if err := json.NewDecoder(resp.Body).Decode(&subject); err != nil {
		return bangumiSubject{}, err
	}

	return subject, nil
}

// GetBangumiMetadata 获取指定 Bangumi 条目的完整元数据。
func GetBangumiMetadata(subjectID string) (Result, error) {
	subject, err := getBangumiSubject(subjectID)
	if err != nil {
		return Result{}, err
	}

	basic := buildBangumiBasicInfo(subject)

	return Result{
		Source:            SourceBangumi,
		ExternalID:        strconv.Itoa(subject.ID),
		CompanyCandidates: basic.Companies,
		Year:              basic.Year,
		Description:       basic.Description,
		Tags:              buildBangumiTags(subject.Tags),
		Covers:            buildBangumiImages(subject.Images),
		Backgrounds:       nil,
	}, nil
}
