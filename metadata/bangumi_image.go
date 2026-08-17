package metadata

import "strings"

type bangumiImages struct {
	Large  string `json:"large"`
	Common string `json:"common"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
	Grid   string `json:"grid"`
}

// GetBangumiImages 获取指定 Bangumi 条目的封面候选。
func GetBangumiImages(subjectID string) ([]ImageCandidate, error) {
	subject, err := getBangumiSubject(subjectID)
	if err != nil {
		return nil, err
	}

	return buildBangumiImages(subject.Images), nil
}

func buildBangumiImages(images *bangumiImages) []ImageCandidate {
	if images == nil {
		return nil
	}

	url := strings.TrimSpace(images.Large)
	if url == "" {
		url = strings.TrimSpace(images.Common)
	}

	if url == "" {
		return nil
	}

	var thumbnail *string

	if value := strings.TrimSpace(images.Common); value != "" {
		thumbnail = &value
	}

	var width *int
	var height *int

	w, h, err := getRemoteImageSize(url)
	if err == nil && w > 0 && h > 0 {
		width = &w
		height = &h
	}

	return []ImageCandidate{
		{
			Source:    SourceBangumi,
			URL:       url,
			Thumbnail: thumbnail,
			Width:     width,
			Height:    height,
		},
	}
}
