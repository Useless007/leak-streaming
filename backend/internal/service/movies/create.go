package movies

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	domain "github.com/leak-streaming/leak-streaming/backend/internal/domain/movies"
	"github.com/leak-streaming/leak-streaming/backend/internal/persistence/repository"
)

const maxSlugAttempts = 5

type CreateMovieInput struct {
	Title             string
	Synopsis          string
	PosterURL         string
	AvailabilityStart string
	AvailabilityEnd   string
	IsVisible         bool
	StreamURL         string
	DRMKeyID          string
	AllowedHosts      []string
	Captions          []CaptionInput
}

type CaptionInput struct {
	LanguageCode string
	Label        string
	CaptionURL   string
}

type ValidationError struct {
	Fields map[string]string
}

func (e ValidationError) Error() string {
	return "invalid input"
}

var (
	ErrDuplicateMovieTitle = errors.New("movie title already exists")
)

func (s *Service) CreateMovie(ctx context.Context, input CreateMovieInput) (domain.Movie, error) {
	if s == nil || s.repo == nil {
		return domain.Movie{}, errors.New("movie service not configured")
	}

	issues := make(map[string]string)

	title := strings.TrimSpace(input.Title)
	if title == "" {
		issues["title"] = "กรุณาระบุชื่อเรื่อง"
	} else if utf8.RuneCountInString(title) > 255 {
		issues["title"] = "ชื่อเรื่องต้องไม่ยาวเกิน 255 ตัวอักษร"
	}

	synopsis := strings.TrimSpace(input.Synopsis)

	posterURL := strings.TrimSpace(input.PosterURL)
	if posterURL != "" && !isValidHTTPURL(posterURL) {
		issues["posterUrl"] = "โปสเตอร์ต้องเป็น URL แบบ http(s)"
	}

	streamURL := strings.TrimSpace(input.StreamURL)
	if streamURL == "" {
		issues["streamUrl"] = "กรุณาระบุลิงก์ .m3u8"
	} else if !isValidStreamURL(streamURL) {
		issues["streamUrl"] = "ต้องเป็น URL แบบ http(s) และลงท้ายด้วย .m3u8"
	}

	availabilityStart := parseOptionalTime(input.AvailabilityStart, "availabilityStart", issues)
	availabilityEnd := parseOptionalTime(input.AvailabilityEnd, "availabilityEnd", issues)
	if availabilityStart != nil && availabilityEnd != nil && availabilityEnd.Before(*availabilityStart) {
		issues["availabilityEnd"] = "วันที่สิ้นสุดต้องอยู่หลังหรือเท่ากับวันที่เริ่มฉาย"
	}

	drmKeyID := strings.TrimSpace(input.DRMKeyID)

	normalizedCaptions, captionIssues := normalizeCaptions(input.Captions)
	for field, message := range captionIssues {
		issues[field] = message
	}

	if len(issues) > 0 {
		return domain.Movie{}, ValidationError{Fields: issues}
	}

	slugBase := slugify(title)
	if slugBase == "" {
		return domain.Movie{}, ValidationError{Fields: map[string]string{"title": "ไม่สามารถสร้าง slug จากชื่อเรื่องนี้ได้"}}
	}
	if utf8.RuneCountInString(slugBase) > 120 {
		slugBase = string([]rune(slugBase)[:120])
	}

	allowedHosts := normalizeAllowedHosts(streamURL, input.AllowedHosts)
	params := repository.CreateMovieParams{
		Title:             title,
		Synopsis:          synopsis,
		PosterURL:         posterURL,
		AvailabilityStart: availabilityStart,
		AvailabilityEnd:   availabilityEnd,
		IsVisible:         input.IsVisible,
		StreamURL:         streamURL,
		DRMKeyID:          drmKeyID,
		AllowedHosts:      allowedHosts,
		Captions:          normalizedCaptions,
	}

	slug := slugBase
	for attempt := 0; attempt < maxSlugAttempts; attempt++ {
		params.Slug = slug
		movie, err := s.repo.CreateMovie(ctx, params)
		if err == nil {
			return movie, nil
		}

		if errors.Is(err, repository.ErrDuplicateSlug) {
			slug = fmt.Sprintf("%s-%d", slugBase, attempt+2)
			if utf8.RuneCountInString(slug) > 128 {
				runes := []rune(slug)
				if len(runes) > 128 {
					slug = string(runes[:128])
				}
			}
			continue
		}

		if errors.Is(err, repository.ErrDuplicateTitle) {
			return domain.Movie{}, ErrDuplicateMovieTitle
		}

		return domain.Movie{}, err
	}

	return domain.Movie{}, ErrDuplicateMovieTitle
}

func parseOptionalTime(raw string, field string, issues map[string]string) *time.Time {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		issues[field] = "รูปแบบวันที่ต้องเป็น RFC3339"
		return nil
	}
	t := parsed.UTC()
	return &t
}

func normalizeCaptions(inputs []CaptionInput) ([]domain.Caption, map[string]string) {
	issues := make(map[string]string)
	normalized := make([]domain.Caption, 0, len(inputs))
	seenLanguages := make(map[string]struct{})

	for idx, input := range inputs {
		lang := strings.TrimSpace(input.LanguageCode)
		label := strings.TrimSpace(input.Label)
		captionURL := strings.TrimSpace(input.CaptionURL)

		if lang == "" && label == "" && captionURL == "" {
			continue
		}

		fieldPrefix := fmt.Sprintf("captions.%d", idx)
		valid := true

		if lang == "" {
			issues[fieldPrefix+".languageCode"] = "กรุณาระบุรหัสภาษา"
			valid = false
		} else {
			lang = strings.ToLower(lang)
			if utf8.RuneCountInString(lang) < 2 || utf8.RuneCountInString(lang) > 10 {
				issues[fieldPrefix+".languageCode"] = "รหัสภาษาต้องมีความยาว 2-10 ตัว"
				valid = false
			}
			if _, exists := seenLanguages[lang]; exists {
				issues[fieldPrefix+".languageCode"] = "ภาษานี้ถูกเพิ่มแล้ว"
				valid = false
			} else {
				seenLanguages[lang] = struct{}{}
			}
		}

		if label == "" {
			issues[fieldPrefix+".label"] = "กรุณาระบุชื่อคำบรรยาย"
			valid = false
		}

		if captionURL == "" {
			issues[fieldPrefix+".captionUrl"] = "กรุณาระบุ URL ของคำบรรยาย"
			valid = false
		} else if !isValidCaptionURL(captionURL) {
			issues[fieldPrefix+".captionUrl"] = "ต้องเป็น URL แบบ http(s) หรือ path ที่ขึ้นต้นด้วย /"
			valid = false
		}

		if valid {
			normalized = append(normalized, domain.Caption{
				LanguageCode: lang,
				Label:        label,
				CaptionURL:   captionURL,
			})
		}
	}

	return normalized, issues
}

func normalizeAllowedHosts(streamURL string, provided []string) []string {
	seen := make(map[string]struct{})
	hosts := make([]string, 0, len(provided)+1)

	if parsed, err := url.Parse(streamURL); err == nil {
		if host := strings.ToLower(parsed.Hostname()); host != "" {
			seen[host] = struct{}{}
			hosts = append(hosts, host)
		}
	}

	for _, raw := range provided {
		host := sanitizeHost(raw)
		if host == "" {
			continue
		}
		if _, exists := seen[host]; exists {
			continue
		}
		seen[host] = struct{}{}
		hosts = append(hosts, host)
	}

	return hosts
}

func sanitizeHost(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	if strings.Contains(value, "://") {
		parsed, err := url.Parse(value)
		if err != nil {
			return ""
		}
		value = parsed.Hostname()
	} else {
		if slash := strings.IndexRune(value, '/'); slash >= 0 {
			value = value[:slash]
		}
		if colon := strings.IndexRune(value, ':'); colon >= 0 {
			value = value[:colon]
		}
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return strings.ToLower(value)
}

func isValidHTTPURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	return parsed.Host != ""
}

func isValidStreamURL(raw string) bool {
	if !isValidHTTPURL(raw) {
		return false
	}
	parsed, _ := url.Parse(raw)
	return strings.HasSuffix(strings.ToLower(parsed.Path), ".m3u8")
}

func isValidCaptionURL(raw string) bool {
	if strings.HasPrefix(raw, "/") {
		return true
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	return parsed.Host != ""
}

func slugify(value string) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.Grow(len(runes))
	prevHyphen := false
	for _, r := range runes {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			builder.WriteRune(unicode.ToLower(r))
			prevHyphen = false
		case r == ' ' || r == '-' || r == '_':
			if !prevHyphen && builder.Len() > 0 {
				builder.WriteRune('-')
				prevHyphen = true
			}
		default:
			// skip other characters
		}
	}
	slug := builder.String()
	slug = strings.Trim(slug, "-")
	return slug
}
