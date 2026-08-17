package metadata

type bangumiTag struct {
	Name string `json:"name"`
}

// GetBangumiTags 获取指定 Bangumi 条目的 Tag 候选。
func GetBangumiTags(subjectID string) ([]string, error) {
	return nil, errNotImplemented
}
