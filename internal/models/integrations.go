package models

import "time"

// APICredential stores credentials for external API services.
type APICredential struct {
	Id          int64
	UserId      int64
	Service     string
	Credentials string
	Created     time.Time
	Updated     time.Time
}

// PageSpeedResult stores Google PageSpeed Insights results.
type PageSpeedResult struct {
	Id               int64
	PageReportId     int64
	CrawlId          int64
	PerformanceScore float64
	FCP              float64
	LCP              float64
	TBT              float64
	CLS              float64
	SpeedIndex       float64
	TTI              float64
	Recommendations  string
	Created          time.Time
}

// API service name constants.
const (
	ServicePageSpeed      = "pagespeed"
	ServiceSearchConsole  = "searchconsole"
	ServiceGA4            = "ga4"
	ServiceAhrefs         = "ahrefs"
)
