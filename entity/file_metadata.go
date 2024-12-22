package entity

import "time"

type FileMetadata struct {
	ID            uint
	GuildID       string     `gorm:"column:guild_id"`
	ChannelID     string     `gorm:"column:channel_id"`
	MessageID     string     `gorm:"column:message_id"`
	URL           string     `gorm:"column:url"`
	Filename      string     `gorm:"column:filename"`
	ExpiredAt     *time.Time `gorm:"column:expired_at"`
	IssuedAt      *time.Time `gorm:"column:issued_at"`
	Hash          string
	LastRotatedAt *time.Time `gorm:"column:last_rotated_at"`
	LastAccessed  *time.Time `gorm:"column:last_accessed"`
	AccessCount   int        `gorm:"column:access_count"`
}
