package service

import "GalgameBox/metadata"

// ResolvedSource 表示导入时已匹配的外部条目。
type ResolvedSource struct {
	Source     metadata.Source
	ExternalID string
}

// MetadataSourceIssue 表示单个元数据源查询失败。
// 某一个来源失败时，不应阻止其他来源继续工作。
type MetadataSourceIssue struct {
	Source  metadata.Source
	Message string
}

// MetadataCollection 表示一次多数据源查询的结果。
type MetadataCollection struct {
	Results []metadata.Result
	Issues  []MetadataSourceIssue
}

// ImportMetadataResult 表示最终交给调用方的元数据查询结果。
type ImportMetadataResult struct {
	Draft  ImportDraft
	Issues []MetadataSourceIssue
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
