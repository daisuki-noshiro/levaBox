package service

import (
	"strings"

	"GalgameBox/metadata"
)

// yearSourcePriority 表示 Year 字段特别信任的来源。
//
// 如果这些来源没有年份，才按照 Results 原有顺序寻找其他来源。
var yearSourcePriority = []metadata.Source{
	metadata.SourceVNDB,
}

// descriptionSourcePriority 表示 Description 字段特别信任的来源。
//
// 如果这些来源没有简介，才按照 Results 原有顺序寻找其他来源。
var descriptionSourcePriority = []metadata.Source{
	metadata.SourceBangumi,
}

// BuildImportDraft 将多个来源返回的 Result 合并到导入草稿。
//
// base 保存 StartImport 已经确定的本地信息，
// 例如启动路径、工作目录、搜索词和标题。
// 本函数只负责补充网络元数据。
func BuildImportDraft(base ImportDraft, results []metadata.Result) ImportDraft {
	// Title 始终以用户确认的搜索词为基础，
	// 不使用外部元数据源修改标题。
	base.SearchKeyword = strings.TrimSpace(base.SearchKeyword)
	base.Title = base.SearchKeyword

	base.Company = chooseCompany(results)
	base.Year = chooseYear(results)
	base.Description = chooseDescription(results)

	base.TagCandidates = mergeTags(results)
	base.CoverCandidates = mergeCoverCandidates(results)
	base.BackgroundCandidates = mergeBackgroundCandidates(results)

	base.Sources = buildResolvedSources(results)

	return base
}

// chooseCompany 从多个来源的公司候选中寻找共同项。
//
// 只有一个来源提供公司时，直接使用它的第一个候选。
// 多个来源都提供公司时，只有存在共同候选才自动选择。
// 找不到共同项则留空，交给用户确认。
func chooseCompany(results []metadata.Result) string {
	companyGroups := make([][]string, 0)

	for _, result := range results {
		companies := uniqueStrings(result.CompanyCandidates)

		if len(companies) == 0 {
			continue
		}

		companyGroups = append(
			companyGroups,
			companies,
		)
	}

	if len(companyGroups) == 0 {
		return ""
	}

	if len(companyGroups) == 1 {
		return companyGroups[0][0]
	}

	// common 保存仍然被所有来源共同支持的候选。
	common := make(map[string]bool)

	for _, company := range companyGroups[0] {
		common[normalizeTextKey(company)] = true
	}

	for i := 1; i < len(companyGroups); i++ {
		current := make(map[string]bool)

		for _, company := range companyGroups[i] {
			current[normalizeTextKey(company)] = true
		}

		for company := range common {
			if !current[company] {
				delete(common, company)
			}
		}
	}

	if len(common) == 0 {
		return ""
	}

	// 如果存在多个共同项，
	// 按第一个来源原本提供的顺序选择第一个。
	for _, company := range companyGroups[0] {
		if common[normalizeTextKey(company)] {
			return company
		}
	}

	return ""
}

// chooseYear 根据字段优先级选择默认年份。
func chooseYear(results []metadata.Result) *int {
	// 先查 Year 特别优先的来源。
	for _, source := range yearSourcePriority {
		for _, result := range results {
			if result.Source != source || result.Year == nil {
				continue
			}

			value := *result.Year
			return &value
		}
	}

	// 优先来源没有数据时，
	// 按 Results 原有顺序使用第一个有效值。
	for _, result := range results {
		if result.Year == nil {
			continue
		}

		value := *result.Year
		return &value
	}

	return nil
}

// chooseDescription 根据字段优先级选择默认简介。
func chooseDescription(results []metadata.Result) *string {
	// 先查 Description 特别优先的来源。
	for _, source := range descriptionSourcePriority {
		for _, result := range results {
			if result.Source != source {
				continue
			}

			if value := cleanOptionalString(result.Description); value != nil {
				return value
			}
		}
	}

	// 优先来源没有数据时，
	// 按 Results 原有顺序使用第一个有效值。
	for _, result := range results {
		if value := cleanOptionalString(result.Description); value != nil {
			return value
		}
	}

	return nil
}

// mergeTags 合并所有来源提供的 Tag 候选并去重。
// 顺序保持数据源查询顺序。
func mergeTags(results []metadata.Result) []string {
	tags := make([]string, 0)
	seen := make(map[string]bool)

	for _, result := range results {
		for _, tag := range result.Tags {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}

			key := normalizeTextKey(tag)
			if seen[key] {
				continue
			}

			seen[key] = true
			tags = append(tags, tag)
		}
	}

	return tags
}

// mergeCoverCandidates 合并所有来源提供的封面候选。
func mergeCoverCandidates(
	results []metadata.Result,
) []metadata.ImageCandidate {
	groups := make([][]metadata.ImageCandidate, 0, len(results))

	for _, result := range results {
		groups = append(groups, result.Covers)
	}

	return mergeImageGroups(groups)
}

// mergeBackgroundCandidates 合并所有来源提供的背景候选。
func mergeBackgroundCandidates(
	results []metadata.Result,
) []metadata.ImageCandidate {
	groups := make([][]metadata.ImageCandidate, 0, len(results))

	for _, result := range results {
		groups = append(groups, result.Backgrounds)
	}

	return mergeImageGroups(groups)
}

// mergeImageGroups 合并图片候选，并按照 URL 去重。
func mergeImageGroups(
	groups [][]metadata.ImageCandidate,
) []metadata.ImageCandidate {
	images := make([]metadata.ImageCandidate, 0)
	seen := make(map[string]bool)

	for _, group := range groups {
		for _, image := range group {
			url := strings.TrimSpace(image.URL)
			if url == "" {
				continue
			}

			if seen[url] {
				continue
			}

			seen[url] = true
			images = append(images, image)
		}
	}

	return images
}

// buildResolvedSources 记录成功获得元数据的外部条目。
//
// 同一个来源最多保留一个 ExternalID。
func buildResolvedSources(
	results []metadata.Result,
) []ResolvedSource {
	sources := make([]ResolvedSource, 0, len(results))
	seen := make(map[metadata.Source]bool)

	for _, result := range results {
		externalID := strings.TrimSpace(result.ExternalID)

		if result.Source == "" || externalID == "" {
			continue
		}

		if seen[result.Source] {
			continue
		}

		seen[result.Source] = true

		sources = append(
			sources,
			ResolvedSource{
				Source:     result.Source,
				ExternalID: externalID,
			},
		)
	}

	return sources
}

// uniqueStrings 清理字符串列表中的空值和重复值。
func uniqueStrings(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]bool)

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		key := normalizeTextKey(value)
		if seen[key] {
			continue
		}

		seen[key] = true
		result = append(result, value)
	}

	return result
}

// normalizeTextKey 用于字符串比较和去重。
// 保留最终展示内容，只把比较键统一成小写并去除首尾空格。
func normalizeTextKey(value string) string {
	return strings.ToLower(
		strings.TrimSpace(value),
	)
}

// cleanOptionalString 清理可选字符串。
// 空字符串统一视为没有数据。
func cleanOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	cleaned := strings.TrimSpace(*value)
	if cleaned == "" {
		return nil
	}

	return &cleaned
}
