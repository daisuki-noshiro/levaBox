package metadata

import "errors"

var errNotImplemented = errors.New("暂未实现")

// Source 表示元数据信息源。
type Source string

const (
	SourceVNDB    Source = "vndb"
	SourceBangumi Source = "bangumi"
)

// ImageCandidate 表示可供用户选择的远程图片。
type ImageCandidate struct {
	Source    Source
	URL       string
	Thumbnail *string
	Width     *int
	Height    *int
}

// Result 表示单个信息源返回的统一元数据。
type Result struct {
	Source     Source
	ExternalID string

	CompanyCandidates []string
	Year              *int
	Description       *string

	Tags []string

	Covers      []ImageCandidate
	Backgrounds []ImageCandidate
}
