package metadata

import (
	"strconv"
	"strings"
)

// BangumiBasicInfo 表示 Bangumi 提供的基础资料。
type BangumiBasicInfo struct {
	Companies   []string
	Year        *int
	Description *string
}

type bangumiInfobox struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// GetBangumiBasicInfo 获取指定 Bangumi 条目的基础资料。
func GetBangumiBasicInfo(subjectID string) (BangumiBasicInfo, error) {
	subject, err := getBangumiSubject(subjectID)
	if err != nil {
		return BangumiBasicInfo{}, err
	}

	return buildBangumiBasicInfo(subject), nil
}

// GetBangumiDescription 获取指定 Bangumi 条目的简介。
func GetBangumiDescription(subjectID string) (*string, error) {
	subject, err := getBangumiSubject(subjectID)
	if err != nil {
		return nil, err
	}

	return buildBangumiDescription(subject.Summary), nil
}

func buildBangumiBasicInfo(subject bangumiSubject) BangumiBasicInfo {
	return BangumiBasicInfo{
		Companies:   buildBangumiCompanies(subject.Infobox),
		Year:        parseBangumiYear(subject.Date),
		Description: buildBangumiDescription(subject.Summary),
	}
}

func buildBangumiCompanies(infobox []bangumiInfobox) []string {
	developerKeys := map[string]bool{
		"开发":   true,
		"开发商":  true,
		"游戏开发": true,
	}

	companies := make([]string, 0)
	seen := make(map[string]bool)

	for _, item := range infobox {
		key := strings.TrimSpace(item.Key)

		if !developerKeys[key] {
			continue
		}

		for _, value := range bangumiInfoboxStrings(item.Value) {
			for _, company := range splitBangumiCompanies(value) {
				if seen[company] {
					continue
				}

				seen[company] = true
				companies = append(companies, company)
			}
		}
	}

	return companies
}

// splitBangumiCompanies 将 Bangumi 中组合展示的开发商拆成独立候选。
func splitBangumiCompanies(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == '/' || r == '／'
	})

	companies := make([]string, 0, len(parts))

	for _, part := range parts {
		company := strings.TrimSpace(part)
		if company == "" {
			continue
		}

		companies = append(companies, company)
	}

	return companies
}

func bangumiInfoboxStrings(value any) []string {
	switch value := value.(type) {
	case string:
		value = strings.TrimSpace(value)

		if value == "" {
			return nil
		}

		return []string{value}

	case []any:
		var values []string

		for _, item := range value {
			values = append(
				values,
				bangumiInfoboxStrings(item)...,
			)
		}

		return values

	case map[string]any:
		if nestedValue, exists := value["v"]; exists {
			return bangumiInfoboxStrings(nestedValue)
		}
	}

	return nil
}

func parseBangumiYear(date string) *int {
	date = strings.TrimSpace(date)

	if len(date) < 4 {
		return nil
	}

	year, err := strconv.Atoi(date[:4])
	if err != nil {
		return nil
	}

	return &year
}

func buildBangumiDescription(summary string) *string {
	summary = strings.TrimSpace(summary)

	if summary == "" {
		return nil
	}

	return &summary
}
