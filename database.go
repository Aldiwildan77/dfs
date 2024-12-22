package main

import (
	"log/slog"
	"os"
	"strings"

	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(databaseURL string) (*gorm.DB, error) {
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			a.Value = slog.StringValue(strings.ToLower(a.Value.String()))
		}
		return a
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replaceAttr,
	}))

	gormLogger := slogGorm.New(
		slogGorm.WithHandler(logger.Handler()),
		slogGorm.WithTraceAll(),
	)

	db, err := gorm.Open(mysql.Open(databaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	return db.Debug(), nil
}
