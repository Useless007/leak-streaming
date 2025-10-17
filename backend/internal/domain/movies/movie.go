package movies

import "time"

type Movie struct {
	ID                 string
	Slug               string
	Title              string
	Synopsis           string
	PosterURL          string
	AvailabilityStart  time.Time
	AvailabilityEnd    time.Time
	IsVisible          bool
	StreamURL          string
	DRMKeyID           string
	Captions           []Caption
	AllowedStreamHosts []string
}

type Caption struct {
	LanguageCode string
	Label        string
	CaptionURL   string
}

func (m Movie) IsAvailable(now time.Time) bool {
	if !m.IsVisible {
		return false
	}
	if !m.AvailabilityStart.IsZero() && now.Before(m.AvailabilityStart) {
		return false
	}
	if !m.AvailabilityEnd.IsZero() && now.After(m.AvailabilityEnd) {
		return false
	}
	return true
}
