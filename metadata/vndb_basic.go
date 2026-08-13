package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type VNDBBasicInfo struct {
	Titles      []string
	Companies   []string
	Year        *int
	Description *string
}

// 发给 VNDB /vn 的请求。
type vndbBasicRequest struct {
	Filters []any  `json:"filters"`
	Fields  string `json:"fields"`
	Results int    `json:"results"`
}

// VNDB /vn 返回的单条原始数据。
type vndbBasicResult struct {
	Title       string          `json:"title"`
	AltTitle    *string         `json:"alttitle"`
	Released    *string         `json:"released"`
	Description *string         `json:"description"`
	Developers  []vndbDeveloper `json:"developers"`
}

type vndbBasicResponse struct {
	Results []vndbBasicResult `json:"results"`
}

// GetVNDBBasicInfo 根据已经确认的 VNDB ID 获取全部基础资料。
func GetVNDBBasicInfo(vndbID string) (VNDBBasicInfo, error) {
	result, err := queryVNDBByID(
		vndbID,
		"title,alttitle,released,description,developers.name",
	)
	if err != nil {
		return VNDBBasicInfo{}, err
	}

	return VNDBBasicInfo{
		Titles:      buildTitleCandidates(result),
		Companies:   buildCompanyCandidates(result),
		Year:        parseVNDBYear(result.Released),
		Description: result.Description,
	}, nil
}

// queryVNDBByID 根据 VNDB ID 查询指定字段，并返回一条 VN 原始数据。
func queryVNDBByID(vndbID string, fields string) (vndbBasicResult, error) {
	requestData := vndbBasicRequest{
		Filters: []any{
			"id",
			"=",
			vndbID,
		},
		Fields:  fields,
		Results: 1,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		return vndbBasicResult{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		vndbVNURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return vndbBasicResult{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := vndbHTTPClient.Do(req)
	if err != nil {
		return vndbBasicResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return vndbBasicResult{}, fmt.Errorf(
			"VNDB 请求失败，状态码: %d",
			resp.StatusCode,
		)
	}

	var response vndbBasicResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return vndbBasicResult{}, err
	}

	if len(response.Results) == 0 {
		return vndbBasicResult{}, fmt.Errorf(
			"未找到 VNDB 条目: %s",
			vndbID,
		)
	}

	return response.Results[0], nil
}

// buildTitleCandidates 把 title / alttitle 整理成标题候选。
func buildTitleCandidates(result vndbBasicResult) []string {
	titles := make([]string, 0, 2)

	title := strings.TrimSpace(result.Title)

	if title != "" {
		titles = append(titles, title)
	}

	if result.AltTitle != nil {
		altTitle := strings.TrimSpace(*result.AltTitle)

		if altTitle != "" && altTitle != title {
			titles = append(titles, altTitle)
		}
	}

	return titles
}

// buildCompanyCandidates 把 developers 整理成公司候选。
func buildCompanyCandidates(result vndbBasicResult) []string {
	companies := make([]string, 0, len(result.Developers))
	seen := make(map[string]bool)

	for _, developer := range result.Developers {
		name := strings.TrimSpace(developer.Name)

		if name == "" || seen[name] {
			continue
		}

		seen[name] = true
		companies = append(companies, name)
	}

	return companies
}

// parseVNDBYear 从 VNDB release date 中提取年份。
func parseVNDBYear(released *string) *int {
	if released == nil || len(*released) < 4 {
		return nil
	}

	year, err := strconv.Atoi((*released)[:4])
	if err != nil {
		return nil
	}

	return &year
}

// GetVNDBTitleCandidates 根据 VNDB ID 获取标题候选。
// 用于后续单独编辑游戏标题。
func GetVNDBTitleCandidates(vndbID string) ([]string, error) {
	result, err := queryVNDBByID(
		vndbID,
		"title,alttitle",
	)
	if err != nil {
		return nil, err
	}

	return buildTitleCandidates(result), nil
}

// GetVNDBCompanyCandidates 根据 VNDB ID 获取开发商候选。
// 用于后续单独编辑游戏公司。
func GetVNDBCompanyCandidates(vndbID string) ([]string, error) {
	result, err := queryVNDBByID(
		vndbID,
		"developers.name",
	)
	if err != nil {
		return nil, err
	}

	return buildCompanyCandidates(result), nil
}

// GetVNDBYearCandidate 根据 VNDB ID 获取发行年份候选。
// VNDB 无法确定年份时返回 nil。
func GetVNDBYearCandidate(vndbID string) (*int, error) {
	result, err := queryVNDBByID(
		vndbID,
		"released",
	)
	if err != nil {
		return nil, err
	}

	return parseVNDBYear(result.Released), nil
}

// GetVNDBDescription 根据 VNDB ID 获取简介。
// VNDB 没有简介时返回 nil。
func GetVNDBDescription(vndbID string) (*string, error) {
	result, err := queryVNDBByID(
		vndbID,
		"description",
	)
	if err != nil {
		return nil, err
	}

	return result.Description, nil
}
