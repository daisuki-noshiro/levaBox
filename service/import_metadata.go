package service

import (
	"fmt"
	"strings"

	"GalgameBox/metadata"
)

// metadataSourceHandler 描述 service 调用一个元数据源所需要的能力。
//
// Resolve：根据搜索词确定默认条目 ID。
// Fetch：根据条目 ID 获取完整元数据。
type metadataSourceHandler struct {
	Resolve func(string) (string, error)
	Fetch   func(string) (metadata.Result, error)
}

// metadataSourceHandlers 注册当前已经支持的元数据源。
//
// service 不需要知道 VNDB、Bangumi 内部如何请求 API，
// 只需要知道它们分别提供 Resolve 和 Fetch 两个入口。
var metadataSourceHandlers = map[metadata.Source]metadataSourceHandler{
	metadata.SourceVNDB: {
		Resolve: metadata.ResolveVNDBID,
		Fetch:   metadata.GetVNDBMetadata,
	},
	metadata.SourceBangumi: {
		Resolve: metadata.ResolveBangumiID,
		Fetch:   metadata.GetBangumiMetadata,
	},
}

// defaultMetadataSources 是没有用户配置时采用的默认来源顺序。
//
// 顺序同时决定候选数据的默认排列顺序。
var defaultMetadataSources = []metadata.Source{
	metadata.SourceVNDB,
	metadata.SourceBangumi,
}

// DefaultMetadataSources 返回默认元数据源。
//
// 返回副本，避免外部代码意外修改包内默认配置。
func DefaultMetadataSources() []metadata.Source {
	result := make([]metadata.Source, len(defaultMetadataSources))
	copy(result, defaultMetadataSources)
	return result
}

// CollectMetadata 从多个启用的数据源收集元数据。
//
// 单个来源失败只记录到 Issues，不影响其他来源。
// Results 的顺序与 sources 的顺序保持一致。
func CollectMetadata(keyword string, sources []metadata.Source) (MetadataCollection, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return MetadataCollection{}, fmt.Errorf("元数据搜索关键词不能为空")
	}

	collection := MetadataCollection{
		Results: make([]metadata.Result, 0, len(sources)),
		Issues:  make([]MetadataSourceIssue, 0),
	}

	// 避免同一个来源被重复查询。
	seen := make(map[metadata.Source]bool)

	for _, source := range sources {
		if seen[source] {
			continue
		}
		seen[source] = true

		handler, exists := metadataSourceHandlers[source]
		if !exists {
			collection.Issues = append(
				collection.Issues,
				MetadataSourceIssue{
					Source:  source,
					Message: "不支持的元数据源",
				},
			)
			continue
		}

		// 第一步：根据搜索词确定这个来源中的条目 ID。
		externalID, err := handler.Resolve(keyword)
		if err != nil {
			collection.Issues = append(
				collection.Issues,
				MetadataSourceIssue{
					Source: source,
					Message: fmt.Sprintf(
						"匹配条目失败: %v",
						err,
					),
				},
			)
			continue
		}

		externalID = strings.TrimSpace(externalID)
		if externalID == "" {
			collection.Issues = append(
				collection.Issues,
				MetadataSourceIssue{
					Source:  source,
					Message: "匹配结果没有条目 ID",
				},
			)
			continue
		}

		// 第二步：根据已经确认的条目 ID 获取完整元数据。
		result, err := handler.Fetch(externalID)
		if err != nil {
			collection.Issues = append(
				collection.Issues,
				MetadataSourceIssue{
					Source: source,
					Message: fmt.Sprintf(
						"获取元数据失败: %v",
						err,
					),
				},
			)
			continue
		}

		collection.Results = append(
			collection.Results,
			result,
		)
	}

	return collection, nil
}
