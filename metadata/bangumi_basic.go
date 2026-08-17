package metadata

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
	return BangumiBasicInfo{}, errNotImplemented
}

// GetBangumiDescription 获取指定 Bangumi 条目的简介。
func GetBangumiDescription(subjectID string) (*string, error) {
	return nil, errNotImplemented
}
