package model

import "time"

type GameProgress string
type BackgroundType string

const (
	ProgressNotStarted GameProgress = "not_started"
	ProgressPlaying    GameProgress = "playing"
	ProgressCompleted  GameProgress = "completed"

	BackgroundImage BackgroundType = "image"
	BackgroundVideo BackgroundType = "video"
)

type Background struct {
	Type BackgroundType
	Path string
}

// 启动游戏的配置
type LaunchConfig struct {
	ExecutablePath   string
	WorkingDirectory string
}

type Game struct {
	ID     string
	Launch LaunchConfig
	
	//基本资料
	Title       string
	Company     string
	Year        int
	Description string

	//媒体
	CoverPath  string
	Background Background
	BGMPath    string
	BGMEnabled bool

	// 用户状态
	Favorite         bool
	Progress         GameProgress
	TotalPlaySeconds int64
	LastPlayedAt     *time.Time

	// 持有Tag
	TagIDs []string
}

type GameQueue struct {
	DefaultGameIDs []string
	CurrentGameIDs []string
}

type Tag struct {
	ID   string
	Name string
}
