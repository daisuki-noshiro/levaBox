package model

// GameMetadataSource 记录游戏关联的外部元数据条目。
type GameMetadataSource struct {
	GameID     string
	Source     string
	ExternalID string
}
