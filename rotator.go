package main

import (
	"context"
	"time"

	"github.com/Aldiwildan77/dfs/entity"
	"github.com/rs/zerolog/log"

	"gorm.io/gorm"
)

type rotator interface {
	StartRotator(ctx context.Context, tick time.Duration)
	RotateFile(ctx context.Context, file entity.FileMetadata) (string, error)
}

type Rotator struct {
	cfg Config
	db  *gorm.DB

	discord discord
}

func NewRotator(cfg Config, db *gorm.DB, discord discord) rotator {
	return &Rotator{
		cfg:     cfg,
		db:      db,
		discord: discord,
	}
}

func (r *Rotator) StartRotator(ctx context.Context, tick time.Duration) {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for range ticker.C {
		log.Info().Msg("Rotating expiring files...")
		r.rotateExpiringFiles(ctx)
	}
}

func (r *Rotator) rotateExpiringFiles(ctx context.Context) {
	var files []entity.FileMetadata
	now := time.Now()
	expiryThreshold := now.Add(time.Duration(r.cfg.RotatorThreshold) * time.Hour)

	if err := r.db.Where("last_rotated_at <= ?", expiryThreshold).Find(&files).Error; err != nil {
		log.Error().Err(err).Msg("Failed to fetch expiring files")
		return
	}

	for _, file := range files {
		_, err := r.RotateFile(ctx, file)
		if err != nil {
			log.Error().Err(err).Msg("Failed to rotate file")
		}
	}
}

func (r *Rotator) RotateFile(ctx context.Context, file entity.FileMetadata) (string, error) {
	now := time.Now()

	req := &FetchFileRequest{
		GuildID:   file.GuildID,
		ChannelID: file.ChannelID,
		MessageID: file.MessageID,
	}

	newURL, err := r.discord.FetchFile(ctx, req)
	if err != nil {
		return "", err
	}

	file.ExpiredAt, err = GetExURL(newURL)
	if err != nil {
		return "", err
	}

	file.IssuedAt, err = GetIsURL(newURL)
	if err != nil {
		return "", err
	}

	file.Hash, err = GetHashURL(newURL)
	if err != nil {
		return "", err
	}

	file.URL = newURL
	file.LastRotatedAt = &now

	if err := r.db.Save(&file).Error; err != nil {
		// Allow the rotator to retry on the next tick
		log.Error().Err(err).Msg("Failed to save rotated file")
		return "", err
	}

	log.Info().Msgf("Rotated file %s", file.Filename)
	return newURL, nil
}
