package metadata

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"
	"time"
)

var remoteImageHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// getRemoteImageSize 获取远程图片的实际尺寸。
func getRemoteImageSize(url string) (width int, height int, err error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return 0, 0, fmt.Errorf("图片 URL 不能为空")
	}

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return 0, 0, err
	}

	req.Header.Set("User-Agent", "daisuki-noshiro/levaBox")
	req.Header.Set("Accept", "image/*")

	resp, err := remoteImageHTTPClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, 0, fmt.Errorf(
			"图片请求失败，状态码: %d",
			resp.StatusCode,
		)
	}

	config, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	return config.Width, config.Height, nil
}
