package adfilter

import (
	"path/filepath"
	"regexp"
	"strings"
)

// Filter represents an ad/spam filter
type Filter struct {
	enabled          bool
	keywords         []string
	wildcardKeywords []string
	patterns         []*regexp.Regexp
	urlPatterns      []*regexp.Regexp
	maxURLCount      int
}

// Config holds the configuration for the ad filter
type Config struct {
	Enabled          bool     `json:"enabled"`
	Keywords         []string `json:"keywords"`
	WildcardKeywords []string `json:"wildcard_keywords"`
	Patterns         []string `json:"patterns"`
	URLPatterns      []string `json:"url_patterns"`
	MaxURLCount      int      `json:"max_url_count"`
}

// NewFilter creates a new ad filter with the given configuration
func NewFilter(config *Config) (*Filter, error) {
	if config == nil {
		return &Filter{
			enabled:          false,
			keywords:         []string{},
			wildcardKeywords: []string{},
			patterns:         []*regexp.Regexp{},
			urlPatterns:      []*regexp.Regexp{},
		}, nil
	}

	keywords := config.Keywords
	if keywords == nil {
		keywords = []string{}
	}

	wildcardKeywords := config.WildcardKeywords
	if wildcardKeywords == nil {
		wildcardKeywords = []string{}
	}

	filter := &Filter{
		enabled:          config.Enabled,
		keywords:         keywords,
		wildcardKeywords: wildcardKeywords,
		maxURLCount:      config.MaxURLCount,
		patterns:         []*regexp.Regexp{},
		urlPatterns:      []*regexp.Regexp{},
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

// matchWildcard checks if a text matches a wildcard pattern
// Supports * (match any sequence) and ? (match single character)
// Case-insensitive matching
func matchWildcard(pattern, text string) bool {
	// Convert both to lowercase for case-insensitive matching
	pattern = strings.ToLower(pattern)
	text = strings.ToLower(text)

	// Try matching at every position in the text for substring matching
	// This allows patterns like "*test*" to match anywhere in the text
	for i := 0; i <= len(text); i++ {
		for j := i; j <= len(text); j++ {
			substring := text[i:j]
			if matched, _ := filepath.Match(pattern, substring); matched {
				return true
			}
		}
	}
	
	// Also check if the pattern matches the entire text
	if matched, err := filepath.Match(pattern, text); err == nil && matched {
		return true
	}

	return false
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

	// Check for exact keyword matches
	for _, keyword := range f.keywords {
		if strings.Contains(lowerContent, strings.ToLower(keyword)) {
			return true
		}
	}

	// Check for wildcard keyword matches
	for _, wildcardKeyword := range f.wildcardKeywords {
		if matchWildcard(wildcardKeyword, content) {
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
			// Common spam keywords (exact match, can be customized)
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
		WildcardKeywords: []string{
			// Wildcard patterns using * and ?
			"*微信*",      // Matches any text containing "微信"
			"*vx*",       // Matches any text containing "vx" (common WeChat abbreviation)
			"加*好友",     // Matches "加" followed by anything and then "好友"
			"*免费*",      // Matches any text containing "免费"
			"*优惠*码",    // Matches text with "优惠" followed by anything and ending with "码"
			"代理*",      // Matches "代理" followed by anything
			"*赚钱*",     // Matches any text containing "赚钱"
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
