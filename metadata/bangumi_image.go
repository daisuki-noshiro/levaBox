package metadata

type bangumiImages struct {
	Large  string `json:"large"`
	Common string `json:"common"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
	Grid   string `json:"grid"`
}

// GetBangumiImages 获取指定 Bangumi 条目的封面候选。
func GetBangumiImages(subjectID string) ([]ImageCandidate, error) {
	return nil, errNotImplemented
}
