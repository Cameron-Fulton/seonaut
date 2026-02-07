package models

import "time"

// Webhook represents a configured webhook endpoint.
type Webhook struct {
	Id        int64
	ProjectId int64
	EventType string
	URL       string
	Secret    string
	Active    bool
	Created   time.Time
}

// WebhookLog represents a log entry for a webhook delivery attempt.
type WebhookLog struct {
	Id           int64
	WebhookId    int64
	EventType    string
	Payload      string
	ResponseCode int
	Error        string
	Created      time.Time
}

// WebhookEvent represents the different event types that can trigger webhooks.
const (
	WebhookEventCrawlStarted   = "crawl.started"
	WebhookEventCrawlCompleted = "crawl.completed"
	WebhookEventIssueFound     = "issue.found"
	WebhookEventCrawlProgress  = "crawl.progress"
)

// WebhookPayload is the data structure sent in webhook requests.
type WebhookPayload struct {
	Event     string      `json:"event"`
	Timestamp time.Time   `json:"timestamp"`
	ProjectId int64       `json:"project_id"`
	Data      interface{} `json:"data"`
}

// CrawlStartedData is the webhook payload for crawl started events.
type CrawlStartedData struct {
	CrawlId int64  `json:"crawl_id"`
	URL     string `json:"url"`
}

// CrawlCompletedData is the webhook payload for crawl completed events.
type CrawlCompletedData struct {
	CrawlId        int64 `json:"crawl_id"`
	TotalURLs      int   `json:"total_urls"`
	TotalIssues    int   `json:"total_issues"`
	CriticalIssues int   `json:"critical_issues"`
	AlertIssues    int   `json:"alert_issues"`
	WarningIssues  int   `json:"warning_issues"`
}

// CrawlProgressData is the webhook payload for crawl progress events.
type CrawlProgressData struct {
	CrawlId    int64  `json:"crawl_id"`
	CrawledURL string `json:"crawled_url"`
	StatusCode int    `json:"status_code"`
	TotalURLs  int    `json:"total_urls"`
}

// IncomingWebhookPayload is the data structure for incoming webhook requests.
type IncomingWebhookPayload struct {
	Action    string                 `json:"action"`
	ProjectId int64                  `json:"project_id"`
	Data      map[string]interface{} `json:"data"`
}
