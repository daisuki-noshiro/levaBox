package service

import "GalgameBox/metadata"

// ResolvedSource 表示导入时已匹配的外部条目。
type ResolvedSource struct {
	Source     metadata.Source
	ExternalID string
}

// ImportDraft 表示导入确认阶段使用的可编辑草稿。
type ImportDraft struct {
	ExecutablePath   string
	WorkingDirectory string
	SearchKeyword    string

	Title       string
	Company     string
	Year        *int
	Description *string

	TagCandidates        []string
	CoverCandidates      []metadata.ImageCandidate
	BackgroundCandidates []metadata.ImageCandidate

	Sources []ResolvedSource
}

// SaveImportRequest 表示用户最终确认后的导入数据。
type SaveImportRequest struct {
	ExecutablePath   string
	WorkingDirectory string

	Title       string
	Company     string
	Year        *int
	Description *string

	Tags []string

	Cover      *metadata.ImageCandidate
	Background *metadata.ImageCandidate

	Sources []ResolvedSource
}
