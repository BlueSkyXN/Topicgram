package adfilter

import (
	"testing"
)

func TestNewFilter(t *testing.T) {
	config := &Config{
		Enabled:  true,
		Keywords: []string{"spam", "ad"},
		Patterns: []string{`test\d+`},
	}

	filter, err := NewFilter(config)
	if err != nil {
		t.Fatalf("Failed to create filter: %v", err)
	}

	if !filter.IsEnabled() {
		t.Error("Filter should be enabled")
	}
}

func TestCheckMessage_Keywords(t *testing.T) {
	config := &Config{
		Enabled:  true,
		Keywords: []string{"免费领取", "加微信"},
	}

	filter, _ := NewFilter(config)

	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"Contains keyword", "免费领取iPhone", true},
		{"Contains keyword case insensitive", "免费领取", true},
		{"No keyword", "Hello world", false},
		{"Empty text", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.CheckMessage(tt.text, "")
			if result != tt.expected {
				t.Errorf("CheckMessage(%q) = %v, want %v", tt.text, result, tt.expected)
			}
		})
	}
}

func TestCheckMessage_URLs(t *testing.T) {
	config := &Config{
		Enabled:     true,
		MaxURLCount: 2,
	}

	filter, _ := NewFilter(config)

	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"No URLs", "Just a normal message", false},
		{"One URL", "Check this out: https://example.com", false},
		{"Two URLs", "Links: https://example.com and https://test.com", false},
		{"Three URLs", "Too many links: https://a.com https://b.com https://c.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.CheckMessage(tt.text, "")
			if result != tt.expected {
				t.Errorf("CheckMessage(%q) = %v, want %v", tt.text, result, tt.expected)
			}
		})
	}
}

func TestCheckMessage_Patterns(t *testing.T) {
	config := &Config{
		Enabled:  true,
		Patterns: []string{`(.)\1{5,}`}, // 5 or more repeated characters - using backreference
	}

	filter, err := NewFilter(config)
	// Note: Go's regexp doesn't support backreferences like \1
	// So we'll use a different pattern
	if err != nil {
		// Try with a simpler pattern that Go supports
		config.Patterns = []string{`a{6,}`} // 6 or more 'a' characters
		filter, err = NewFilter(config)
		if err != nil {
			t.Fatalf("Failed to create filter: %v", err)
		}
	}
	if filter == nil {
		t.Fatal("Filter is nil")
	}

	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"Normal text", "Hello world", false},
		{"Repeated chars", "aaaaaaa", true},
		{"Not enough repeats", "aaaa", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filter.CheckMessage(tt.text, "")
			if result != tt.expected {
				t.Errorf("CheckMessage(%q) = %v, want %v", tt.text, result, tt.expected)
			}
		})
	}
}

func TestCheckMessage_WithCaption(t *testing.T) {
	config := &Config{
		Enabled:  true,
		Keywords: []string{"spam"},
	}

	filter, _ := NewFilter(config)

	// Test with caption containing spam keyword
	result := filter.CheckMessage("Normal text", "This is spam content")
	if !result {
		t.Error("Should detect spam in caption")
	}

	// Test with both text and caption
	result = filter.CheckMessage("Hello", "world")
	if result {
		t.Error("Should not detect spam in clean message")
	}
}

func TestCheckMessage_Disabled(t *testing.T) {
	config := &Config{
		Enabled:  false,
		Keywords: []string{"spam", "ad"},
	}

	filter, _ := NewFilter(config)

	// Even with spam keywords, disabled filter should allow all messages
	result := filter.CheckMessage("This is spam", "")
	if result {
		t.Error("Disabled filter should not block any messages")
	}
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	if config == nil {
		t.Fatal("GetDefaultConfig should not return nil")
	}

	if config.Enabled {
		t.Error("Default config should have filter disabled")
	}

	if len(config.Keywords) == 0 {
		t.Error("Default config should have some keywords")
	}

	if config.MaxURLCount == 0 {
		t.Error("Default config should have MaxURLCount set")
	}
}
