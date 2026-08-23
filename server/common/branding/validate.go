package branding

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	slugPattern  = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	colorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
	reserved     = map[string]struct{}{"api": {}, "admin": {}, "app": {}, "www": {}}
)

func Validate(slug, name, description, logoURL, primaryColor string) error {
	if len(slug) < 3 || len(slug) > 32 || !slugPattern.MatchString(slug) {
		return errors.New("invalid slug")
	}
	if _, exists := reserved[slug]; exists {
		return errors.New("reserved slug")
	}
	nameLength := utf8.RuneCountInString(strings.TrimSpace(name))
	if nameLength < 1 || nameLength > 80 {
		return errors.New("name must contain 1-80 characters")
	}
	if len([]byte(description)) > 1000 {
		return errors.New("description is too long")
	}
	if logoURL != "" {
		parsed, err := url.Parse(logoURL)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
			return errors.New("logo URL must use HTTPS")
		}
	}
	if primaryColor != "" && !colorPattern.MatchString(primaryColor) {
		return errors.New("primary color must be #RRGGBB")
	}
	return nil
}
