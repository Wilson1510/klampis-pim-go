package utils

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// GenerateSlug creates a URL-friendly slug from the given text
func GenerateSlug(text string) string {
	if text == "" {
		return ""
	}

	// Convert to lowercase
	slug := strings.ToLower(text)

	// Remove accents and normalize unicode
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	slug, _, _ = transform.String(t, slug)

	// Replace spaces and special characters with hyphens
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Replace multiple consecutive hyphens with single hyphen
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	return slug
}

// isMn reports whether the rune is a nonspacing mark.
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// GenerateUniqueSlug creates a unique slug by appending numbers if needed
func GenerateUniqueSlug(baseSlug string, existingSlugs []string) string {
	if baseSlug == "" {
		return ""
	}

	// Check if base slug is unique
	if !contains(existingSlugs, baseSlug) {
		return baseSlug
	}

	// If not unique, append numbers until we find a unique one
	counter := 1
	for {
		candidateSlug := baseSlug + "-" + strconv.Itoa(counter)

		if !contains(existingSlugs, candidateSlug) {
			return candidateSlug
		}
		counter++

		// Safety check to prevent infinite loop
		if counter > 1000 {
			break
		}
	}

	return baseSlug
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
