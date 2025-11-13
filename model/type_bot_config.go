package model

import "Topicgram/pkg/adfilter"

type BotConfig struct {
	Token        string
	GroupId      int64
	LanguageCode string

	WebHook struct {
		Host string
	}

	AdFilter *adfilter.Config `json:"ad_filter"`
}
