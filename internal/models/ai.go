package models

import "time"

// AIAnalysis represents the result of an AI analysis on a page.
type AIAnalysis struct {
	Id           int64
	PageReportId int64
	CrawlId      int64
	Provider     string // "claude" or "lmstudio"
	AnalysisType string // "intent", "content_quality", "seo_suggestions"
	Result       string
	Created      time.Time
}

// AIProvider defines the available AI providers.
const (
	AIProviderClaude   = "claude"
	AIProviderLMStudio = "lmstudio"
)

// AIAnalysisType defines the types of AI analysis.
const (
	AIAnalysisIntent         = "intent"
	AIAnalysisContentQuality = "content_quality"
	AIAnalysisSEOSuggestions = "seo_suggestions"
)
