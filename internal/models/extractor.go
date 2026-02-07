package models

import "time"

// CustomExtractor represents a user-defined extraction rule.
type CustomExtractor struct {
	Id        int64
	ProjectId int64
	Name      string
	Type      string // "regex", "css", or "xpath"
	Pattern   string
	Created   time.Time
}

// ExtractionResult represents the result of applying an extractor to a page.
type ExtractionResult struct {
	Id           int64
	PageReportId int64
	CrawlId      int64
	ExtractorId  int64
	Value        string
}

// ExtractorType constants.
const (
	ExtractorTypeRegex = "regex"
	ExtractorTypeCSS   = "css"
	ExtractorTypeXPath = "xpath"
)
