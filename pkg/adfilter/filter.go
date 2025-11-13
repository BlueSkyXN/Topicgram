package adfilter

import (
	"regexp"
	"strings"
)

// Filter represents an ad/spam filter
type Filter struct {
	enabled       bool
	keywords      []string
	patterns      []*regexp.Regexp
	urlPatterns   []*regexp.Regexp
	maxURLCount   int
}

// Config holds the configuration for the ad filter
type Config struct {
	Enabled     bool     `json:"enabled"`
	Keywords    []string `json:"keywords"`
	Patterns    []string `json:"patterns"`
	URLPatterns []string `json:"url_patterns"`
	MaxURLCount int      `json:"max_url_count"`
}

// NewFilter creates a new ad filter with the given configuration
func NewFilter(config *Config) (*Filter, error) {
	if config == nil {
		return &Filter{
			enabled:     false,
			keywords:    []string{},
			patterns:    []*regexp.Regexp{},
			urlPatterns: []*regexp.Regexp{},
		}, nil
	}

	keywords := config.Keywords
	if keywords == nil {
		keywords = []string{}
	}

	filter := &Filter{
		enabled:     config.Enabled,
		keywords:    keywords,
		maxURLCount: config.MaxURLCount,
		patterns:    []*regexp.Regexp{},
		urlPatterns: []*regexp.Regexp{},
	}

	// Compile regex patterns
	for _, pattern := range config.Patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		filter.patterns = append(filter.patterns, re)
	}

	// Compile URL patterns
	for _, pattern := range config.URLPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		filter.urlPatterns = append(filter.urlPatterns, re)
	}

	// Default URL patterns if none provided
	if len(filter.urlPatterns) == 0 {
		defaultURLPattern := regexp.MustCompile(`(?i)(https?://|www\.|t\.me/)[^\s]+`)
		filter.urlPatterns = append(filter.urlPatterns, defaultURLPattern)
	}

	// Set default max URL count if not specified
	if filter.maxURLCount == 0 {
		filter.maxURLCount = 3
	}

	return filter, nil
}

// IsEnabled returns whether the filter is enabled
func (f *Filter) IsEnabled() bool {
	return f.enabled
}

// CheckMessage checks if a message should be blocked
// Returns true if the message is spam/ad, false otherwise
func (f *Filter) CheckMessage(text string, caption string) bool {
	if !f.enabled {
		return false
	}

	// Combine text and caption for checking
	content := text
	if caption != "" {
		if content != "" {
			content += " " + caption
		} else {
			content = caption
		}
	}

	if content == "" {
		return false
	}

	// Convert to lowercase for case-insensitive matching
	lowerContent := strings.ToLower(content)

	// Check for blocked keywords
	for _, keyword := range f.keywords {
		if strings.Contains(lowerContent, strings.ToLower(keyword)) {
			return true
		}
	}

	// Check against regex patterns
	for _, pattern := range f.patterns {
		if pattern.MatchString(content) {
			return true
		}
	}

	// Count URLs and check if exceeds limit
	urlCount := 0
	for _, urlPattern := range f.urlPatterns {
		matches := urlPattern.FindAllString(content, -1)
		urlCount += len(matches)
	}

	if urlCount > f.maxURLCount {
		return true
	}

	return false
}

// GetDefaultConfig returns a default configuration with common spam patterns
func GetDefaultConfig() *Config {
	return &Config{
		Enabled: false, // Disabled by default
		Keywords: []string{
			// Common spam keywords (can be customized)
			"免费领取",
			"点击领取",
			"限时优惠",
			"立即购买",
			"加微信",
			"添加微信",
			"扫码添加",
			"代理加盟",
			"兼职",
			"刷单",
			"投资理财",
			"贷款",
		},
		Patterns: []string{
			// Pattern for many repeated 'a' characters (as example)
			`a{6,}`,
			// Pattern for many repeated numbers
			`\d{10,}`,
		},
		URLPatterns: []string{
			`(?i)(https?://|www\.|t\.me/)[^\s]+`,
		},
		MaxURLCount: 3,
	}
}
