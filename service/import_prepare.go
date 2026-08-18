package service

import (
	"fmt"
	"strings"

	"GalgameBox/metadata"
)

// PrepareImportMetadata 查询启用的元数据源并整理成可确认的导入草稿。
func (s *ImportService) PrepareImportMetadata(
	base ImportDraft,
	sources []metadata.Source,
) (ImportMetadataResult, error) {
	base.SearchKeyword = strings.TrimSpace(base.SearchKeyword)
	if base.SearchKeyword == "" {
		return ImportMetadataResult{}, fmt.Errorf("元数据搜索关键词不能为空")
	}

	if sources == nil {
		sources = DefaultMetadataSources()
	}

	collection, err := CollectMetadata(base.SearchKeyword, sources)
	if err != nil {
		return ImportMetadataResult{}, err
	}

	return ImportMetadataResult{
		Draft:  BuildImportDraft(base, collection.Results),
		Issues: collection.Issues,
	}, nil
}
