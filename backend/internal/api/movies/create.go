package movies

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	service "github.com/leak-streaming/leak-streaming/backend/internal/service/movies"
)

type CreateHandler struct {
	service *service.Service
}

func NewCreateHandler(service *service.Service) *CreateHandler {
	return &CreateHandler{service: service}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.service == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "service unavailable", nil)
		return
	}

	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload createMovieRequest
	if err := decoder.Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "ไม่สามารถอ่านข้อมูลที่ส่งมาได้", nil)
		return
	}

	isVisible := true
	if payload.IsVisible != nil {
		isVisible = *payload.IsVisible
	}

	inputs := make([]service.CaptionInput, 0, len(payload.Captions))
	for _, caption := range payload.Captions {
		inputs = append(inputs, service.CaptionInput{
			LanguageCode: caption.LanguageCode,
			Label:        caption.Label,
			CaptionURL:   caption.CaptionURL,
		})
	}

	movie, err := h.service.CreateMovie(r.Context(), service.CreateMovieInput{
		Title:             payload.Title,
		Synopsis:          payload.Synopsis,
		PosterURL:         payload.PosterURL,
		AvailabilityStart: payload.AvailabilityStart,
		AvailabilityEnd:   payload.AvailabilityEnd,
		IsVisible:         isVisible,
		StreamURL:         payload.StreamURL,
		DRMKeyID:          payload.DRMKeyID,
		AllowedHosts:      payload.AllowedHosts,
		Captions:          inputs,
	})
	if err != nil {
		var validationErr service.ValidationError
		if errors.As(err, &validationErr) {
			writeJSONError(w, http.StatusUnprocessableEntity, "ข้อมูลไม่ถูกต้อง", validationErr.Fields)
			return
		}
		if errors.Is(err, service.ErrDuplicateMovieTitle) {
			writeJSONError(w, http.StatusConflict, "มีภาพยนตร์ที่ใช้ชื่อนี้อยู่แล้ว", nil)
			return
		}

		writeJSONError(w, http.StatusInternalServerError, "ไม่สามารถบันทึกข้อมูลได้", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/movies/%s", movie.Slug))
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(movieResponseFromDomain(movie)); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

type createMovieRequest struct {
	Title             string               `json:"title"`
	Synopsis          string               `json:"synopsis"`
	PosterURL         string               `json:"posterUrl"`
	AvailabilityStart string               `json:"availabilityStart"`
	AvailabilityEnd   string               `json:"availabilityEnd"`
	IsVisible         *bool                `json:"isVisible"`
	StreamURL         string               `json:"streamUrl"`
	DRMKeyID          string               `json:"drmKeyId"`
	AllowedHosts      []string             `json:"allowedHosts"`
	Captions          []createCaptionInput `json:"captions"`
}

type createCaptionInput struct {
	LanguageCode string `json:"languageCode"`
	Label        string `json:"label"`
	CaptionURL   string `json:"captionUrl"`
}

func writeJSONError(w http.ResponseWriter, status int, message string, details map[string]string) {
	response := map[string]any{
		"error": message,
	}
	if len(details) > 0 {
		response["details"] = details
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
