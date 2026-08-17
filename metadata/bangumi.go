package metadata

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

// SearchBangumi 根据游戏名称搜索 Bangumi 游戏条目。
func SearchBangumi(keyword string) ([]BangumiSearchResult, error) {
	return nil, errNotImplemented
}

// ResolveBangumiID 根据搜索结果确定默认 Bangumi 条目。
func ResolveBangumiID(keyword string) (string, error) {
	return "", errNotImplemented
}

// getBangumiSubject 获取指定 Bangumi 条目的完整原始数据。
func getBangumiSubject(subjectID string) (bangumiSubject, error) {
	return bangumiSubject{}, errNotImplemented
}

// GetBangumiMetadata 获取指定 Bangumi 条目的完整元数据。
func GetBangumiMetadata(subjectID string) (Result, error) {
	return Result{}, errNotImplemented
}
