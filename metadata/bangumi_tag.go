package metadata

import "strings"

type bangumiTag struct {
	Name string `json:"name"`
}

// GetBangumiTags 获取指定 Bangumi 条目的 Tag 候选。
func GetBangumiTags(subjectID string) ([]string, error) {
	subject, err := getBangumiSubject(subjectID)
	if err != nil {
		return nil, err
	}

	return buildBangumiTags(subject.Tags), nil
}

func buildBangumiTags(tags []bangumiTag) []string {
	result := make([]string, 0, len(tags))
	seen := make(map[string]bool)

	for _, tag := range tags {
		name := strings.TrimSpace(tag.Name)

		if name == "" || seen[name] {
			continue
		}

		seen[name] = true
		result = append(result, name)
	}

	return result
}
