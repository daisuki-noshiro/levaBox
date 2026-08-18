package service

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"GalgameBox/metadata"
)

type preparedImportMedia struct {
	Directory      string
	CoverPath      string
	BackgroundPath string
}

var importMediaHTTPClient = &http.Client{
	Timeout: 15 * time.Second,
}

// prepareImportMedia 下载用户确认的封面和背景图片。
func (s *ImportService) prepareImportMedia(
	gameID string,
	cover *metadata.ImageCandidate,
	background *metadata.ImageCandidate,
) (preparedImportMedia, error) {
	if cover == nil && background == nil {
		return preparedImportMedia{}, nil
	}

	mediaRoot, err := s.importMediaRoot()
	if err != nil {
		return preparedImportMedia{}, err
	}

	media := preparedImportMedia{
		Directory: filepath.Join(mediaRoot, gameID),
	}
	if err := os.MkdirAll(media.Directory, 0755); err != nil {
		return preparedImportMedia{}, fmt.Errorf("创建媒体目录失败: %w", err)
	}

	if cover != nil {
		media.CoverPath, err = downloadImportImage(
			cover.URL,
			media.Directory,
			"cover",
		)
		if err != nil {
			cleanupImportMedia(media.Directory)
			return preparedImportMedia{}, fmt.Errorf("下载封面失败: %w", err)
		}
	}

	if background != nil {
		media.BackgroundPath, err = downloadImportImage(
			background.URL,
			media.Directory,
			"background",
		)
		if err != nil {
			cleanupImportMedia(media.Directory)
			return preparedImportMedia{}, fmt.Errorf("下载背景失败: %w", err)
		}
	}

	return media, nil
}

func (s *ImportService) importMediaRoot() (string, error) {
	if s.mediaRoot != "" {
		return s.mediaRoot, nil
	}

	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("获取应用数据目录失败: %w", err)
	}

	return filepath.Join(configDirectory, "levaBox", "media"), nil
}

func downloadImportImage(rawURL string, directory string, name string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("图片 URL 不能为空")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("解析图片 URL 失败: %w", err)
	}

	response, err := importMediaHTTPClient.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("请求图片失败: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("图片请求状态码为 %d", response.StatusCode)
	}

	extension := importImageExtension(
		parsedURL.Path,
		response.Header.Get("Content-Type"),
	)
	targetPath := filepath.Join(directory, name+extension)

	file, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("创建图片文件失败: %w", err)
	}

	_, copyErr := io.Copy(file, response.Body)
	closeErr := file.Close()

	if copyErr != nil {
		os.Remove(targetPath)
		return "", fmt.Errorf("写入图片文件失败: %w", copyErr)
	}
	if closeErr != nil {
		os.Remove(targetPath)
		return "", fmt.Errorf("关闭图片文件失败: %w", closeErr)
	}

	return targetPath, nil
}

func importImageExtension(urlPath string, contentType string) string {
	extension := strings.ToLower(path.Ext(urlPath))
	switch extension {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".avif":
		return extension
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err == nil {
		switch strings.ToLower(mediaType) {
		case "image/jpeg":
			return ".jpg"
		case "image/png":
			return ".png"
		case "image/gif":
			return ".gif"
		case "image/webp":
			return ".webp"
		case "image/bmp":
			return ".bmp"
		case "image/avif":
			return ".avif"
		}
	}

	return ".jpg"
}

func cleanupImportMedia(directory string) {
	if directory == "" {
		return
	}

	_ = os.RemoveAll(directory)
}
