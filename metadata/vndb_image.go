package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const vndbReleaseURL = "https://api.vndb.org/kana/release"

// VNDBImages 是一个 VNDB 条目可提供的图片候选。
type VNDBImages struct {
	Covers      []ImageCandidate
	Backgrounds []ImageCandidate
}

// VNDB /release 返回的原始图片。
type vndbReleaseImage struct {
	ID        string  `json:"id"`
	URL       string  `json:"url"`
	Thumbnail string  `json:"thumbnail"`
	Dims      []int   `json:"dims"`
	Type      string  `json:"type"`
	VN        string  `json:"vn"`
	Sexual    float64 `json:"sexual"`
}

type vndbReleaseResult struct {
	Images []vndbReleaseImage `json:"images"`
}

type vndbReleaseResponse struct {
	Results []vndbReleaseResult `json:"results"`
	More    bool                `json:"more"`
}

// GetVNDBImages 根据 VNDB ID 获取封面和背景图候选。
func GetVNDBImages(vndbID string) (VNDBImages, error) {
	vndbID = strings.TrimSpace(vndbID)
	if vndbID == "" {
		return VNDBImages{}, fmt.Errorf("VNDB ID 不能为空")
	}

	var allImages []vndbReleaseImage

	page := 1

	for {
		requestData := struct {
			Filters []any  `json:"filters"`
			Fields  string `json:"fields"`
			Results int    `json:"results"`
			Page    int    `json:"page"`
		}{
			Filters: []any{
				"and",
				[]any{
					"vn",
					"=",
					[]any{"id", "=", vndbID},
				},
				[]any{
					"image",
					"=",
					"dig",
				},
			},
			Fields:  "images.id,images.url,images.thumbnail,images.dims,images.type,images.vn,images.sexual",
			Results: 100,
			Page:    page,
		}

		body, err := json.Marshal(requestData)
		if err != nil {
			return VNDBImages{}, err
		}

		req, err := http.NewRequest(
			http.MethodPost,
			vndbReleaseURL,
			bytes.NewReader(body),
		)
		if err != nil {
			return VNDBImages{}, err
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := vndbHTTPClient.Do(req)
		if err != nil {
			return VNDBImages{}, err
		}

		var response vndbReleaseResponse

		decodeErr := json.NewDecoder(resp.Body).Decode(&response)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return VNDBImages{}, fmt.Errorf(
				"VNDB release 请求失败，状态码: %d",
				resp.StatusCode,
			)
		}

		if decodeErr != nil {
			return VNDBImages{}, decodeErr
		}

		for _, release := range response.Results {
			allImages = append(allImages, release.Images...)
		}

		if !response.More {
			break
		}

		page++
	}

	return buildVNDBImages(vndbID, allImages), nil
}

// buildVNDBImages 整理 dig 图片并根据长宽分类。
func buildVNDBImages(vndbID string, images []vndbReleaseImage) VNDBImages {
	result := VNDBImages{
		Covers:      make([]ImageCandidate, 0),
		Backgrounds: make([]ImageCandidate, 0),
	}

	seen := make(map[string]bool)

	for _, image := range images {
		// 只使用数字版图片。
		if image.Type != "dig" {
			continue
		}

		// bundle release 中，图片可能只属于其中某一个 VN。
		if image.VN != "" && image.VN != vndbID {
			continue
		}

		// 不将被标记为带有性内容的图片作为候选。
		if image.Sexual > 0 {
			continue
		}

		if image.URL == "" || len(image.Dims) != 2 {
			continue
		}

		if seen[image.ID] {
			continue
		}
		seen[image.ID] = true

		width := image.Dims[0]
		height := image.Dims[1]

		var thumbnail *string
		if image.Thumbnail != "" {
			thumbnail = &image.Thumbnail
		}

		candidate := ImageCandidate{
			Source:    SourceVNDB,
			URL:       image.URL,
			Thumbnail: thumbnail,
			Width:     &width,
			Height:    &height,
		}

		switch {
		case width > height:
			result.Backgrounds = append(
				result.Backgrounds,
				candidate,
			)

		case height > width:
			result.Covers = append(
				result.Covers,
				candidate,
			)
		}
	}

	return result
}
