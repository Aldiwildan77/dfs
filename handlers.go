package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Aldiwildan77/dfs/entity"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type handlers interface {
	SaveDiscordLink(w http.ResponseWriter, r *http.Request)
	GetDiscordLink(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	cfg     *Config
	db      *gorm.DB
	rotator rotator
}

func NewHandlers(cfg *Config, db *gorm.DB, rotator rotator) handlers {
	return &Handlers{
		cfg:     cfg,
		db:      db,
		rotator: rotator,
	}
}

func (h *Handlers) SaveDiscordLink(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	var req struct {
		GuildID   string `json:"guild_id" validate:"required"`
		ChannelID string `json:"channel_id" validate:"required"`
		MessageID string `json:"message_id" validate:"required"`
		URL       string `json:"url" validate:"required,url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid request body"})
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		var errMsg []string
		for _, e := range errors {
			errMsg = append(errMsg, fmt.Sprintf("Key: %s Error:Field validation for '%s' failed on the '%s' tag", e.Namespace(), e.Field(), e.Tag()))
		}
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid request body", "errors": errMsg})
		return
	}

	queryParams, err := getQueryParams(req.URL)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid URL"})
		return
	}

	ex := queryParams.Get("ex")
	is := queryParams.Get("is")
	hash := queryParams.Get("hm")

	if ex == "" || is == "" || hash == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid URL Query Parameters"})
		return
	}

	exTime, err := parseHexToUnix(ex)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid Expired URL"})
		return
	}

	isTime, err := parseHexToUnix(is)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid Issued URL"})
		return
	}

	URL, err := url.Parse(req.URL)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid URL"})
		return
	}

	fileName := strings.Split(URL.Path, "/")[len(strings.Split(URL.Path, "/"))-1]

	fileMetadata := entity.FileMetadata{
		GuildID:       req.GuildID,
		ChannelID:     req.ChannelID,
		MessageID:     req.MessageID,
		URL:           req.URL,
		Filename:      fileName,
		Hash:          hash,
		ExpiredAt:     &exTime,
		IssuedAt:      &isTime,
		LastRotatedAt: &now,
		LastAccessed:  &now,
		AccessCount:   0,
	}

	err = h.db.Create(&fileMetadata).Error
	if err != nil {
		if IsDuplicateKeyError(err) {
			WriteJSON(w, http.StatusConflict, map[string]interface{}{"message": "File already exists"})
			return
		}

		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": "Failed to save file metadata"})
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "File metadata saved successfully",
		"data": map[string]interface{}{
			"url":        req.URL,
			"filename":   fileName,
			"expired_at": exTime,
			"issued_at":  isTime,
			"hash":       hash,
		},
	})
}

func (h *Handlers) GetDiscordLink(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid request body"})
		return
	}

	queryParamsURL, err := getQueryParams(req.URL)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid URL"})
		return
	}

	ex := queryParamsURL.Get("ex")
	is := queryParamsURL.Get("is")
	hm := queryParamsURL.Get("hm")

	if ex == "" || is == "" || hm == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid URL Query Parameters"})
		return
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid URL"})
		return
	}

	filteredURL := parsedURL.Scheme + "://" + parsedURL.Host + parsedURL.Path

	var fileMetadata entity.FileMetadata
	if err := h.db.Where("url LIKE ?", filteredURL+"%").First(&fileMetadata).Error; err != nil {
		WriteJSON(w, http.StatusNotFound, map[string]interface{}{"message": "File not found"})
		return
	}

	if fileMetadata.ExpiredAt.Before(now) {
		newURL, err := h.rotator.RotateFile(r.Context(), fileMetadata)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": "Failed to rotate file"})
			return
		}

		newURLParsed, err := url.Parse(newURL)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid Parsed New URL"})
			return
		}

		newEx := newURLParsed.Query().Get("ex")
		newIs := newURLParsed.Query().Get("is")
		newHm := newURLParsed.Query().Get("hm")

		newExTime, err := parseHexToUnix(newEx)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid Expired URL"})
			return
		}

		newIsTime, err := parseHexToUnix(newIs)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid Issued URL"})
			return
		}

		WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": "File rotated successfully",
			"data": map[string]interface{}{
				"url":        newURL,
				"filename":   fileMetadata.Filename,
				"expired_at": newExTime,
				"issued_at":  newIsTime,
				"hash":       newHm,
			},
		})
		return
	}

	fileMetadata.AccessCount++
	fileMetadata.LastAccessed = &now

	if err := h.db.Save(&fileMetadata).Error; err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": "Failed to update file metadata"})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "File accessed successfully",
		"data": map[string]interface{}{
			"url":        fileMetadata.URL,
			"filename":   fileMetadata.Filename,
			"expired_at": fileMetadata.ExpiredAt,
			"issued_at":  fileMetadata.IssuedAt,
			"hash":       fileMetadata.Hash,
		},
	})
}
