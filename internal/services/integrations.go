package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/stjudewashere/seonaut/internal/models"
)

// IntegrationsServiceRepository defines the data access interface for API integrations.
type IntegrationsServiceRepository interface {
	SaveAPICredential(c *models.APICredential) error
	FindAPICredential(userId int64, service string) (*models.APICredential, error)
	DeleteAPICredential(userId int64, service string) error
	SavePageSpeedResult(ps *models.PageSpeedResult) error
	FindPageSpeedResult(pageReportId int64) (*models.PageSpeedResult, error)
}

// IntegrationsConfig holds configuration for external API integrations.
type IntegrationsConfig struct {
	PageSpeedAPIKey string
	AhrefsAPIKey    string
}

// IntegrationsService manages external API integrations.
type IntegrationsService struct {
	repository IntegrationsServiceRepository
	config     *IntegrationsConfig
	client     *http.Client
}

// NewIntegrationsService creates a new IntegrationsService.
func NewIntegrationsService(r IntegrationsServiceRepository, config *IntegrationsConfig) *IntegrationsService {
	return &IntegrationsService{
		repository: r,
		config:     config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RunPageSpeedInsights fetches PageSpeed Insights data for a URL.
func (s *IntegrationsService) RunPageSpeedInsights(pageReport *models.PageReport, crawlId int64) (*models.PageSpeedResult, error) {
	if s.config.PageSpeedAPIKey == "" {
		return nil, fmt.Errorf("PageSpeed Insights API key not configured")
	}

	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=%s&key=%s&strategy=mobile&category=PERFORMANCE",
		url.QueryEscape(pageReport.URL),
		s.config.PageSpeedAPIKey,
	)

	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("PageSpeed API error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var psiResp pageSpeedResponse
	if err := json.Unmarshal(body, &psiResp); err != nil {
		return nil, fmt.Errorf("PageSpeed response parse error: %w", err)
	}

	result := &models.PageSpeedResult{
		PageReportId:     pageReport.Id,
		CrawlId:          crawlId,
		PerformanceScore: psiResp.LighthouseResult.Categories.Performance.Score * 100,
		FCP:              psiResp.LighthouseResult.Audits.FCP.NumericValue,
		LCP:              psiResp.LighthouseResult.Audits.LCP.NumericValue,
		TBT:              psiResp.LighthouseResult.Audits.TBT.NumericValue,
		CLS:              psiResp.LighthouseResult.Audits.CLS.NumericValue,
		SpeedIndex:       psiResp.LighthouseResult.Audits.SpeedIndex.NumericValue,
		TTI:              psiResp.LighthouseResult.Audits.TTI.NumericValue,
	}

	// Serialize recommendations
	recommendations := s.extractRecommendations(psiResp)
	recJSON, _ := json.Marshal(recommendations)
	result.Recommendations = string(recJSON)

	if err := s.repository.SavePageSpeedResult(result); err != nil {
		log.Printf("save PageSpeed result error: %v", err)
	}

	return result, nil
}

// extractRecommendations extracts actionable recommendations from PSI response.
func (s *IntegrationsService) extractRecommendations(resp pageSpeedResponse) []string {
	var recommendations []string

	for _, audit := range resp.LighthouseResult.AuditsList {
		if audit.Score != nil && *audit.Score < 0.9 && audit.Title != "" {
			recommendations = append(recommendations, audit.Title)
		}
	}

	return recommendations
}

// GetPageSpeedResult returns cached PageSpeed results for a page.
func (s *IntegrationsService) GetPageSpeedResult(pageReportId int64) (*models.PageSpeedResult, error) {
	return s.repository.FindPageSpeedResult(pageReportId)
}

// SaveCredentials stores API credentials for a user.
func (s *IntegrationsService) SaveCredentials(userId int64, service string, credentials string) error {
	cred := &models.APICredential{
		UserId:      userId,
		Service:     service,
		Credentials: credentials,
	}
	return s.repository.SaveAPICredential(cred)
}

// GetCredentials returns API credentials for a user and service.
func (s *IntegrationsService) GetCredentials(userId int64, service string) (*models.APICredential, error) {
	return s.repository.FindAPICredential(userId, service)
}

// DeleteCredentials removes API credentials.
func (s *IntegrationsService) DeleteCredentials(userId int64, service string) error {
	return s.repository.DeleteAPICredential(userId, service)
}

// PageSpeed Insights API response structures
type pageSpeedResponse struct {
	LighthouseResult struct {
		Categories struct {
			Performance struct {
				Score float64 `json:"score"`
			} `json:"performance"`
		} `json:"categories"`
		Audits struct {
			FCP        pageSpeedAudit `json:"first-contentful-paint"`
			LCP        pageSpeedAudit `json:"largest-contentful-paint"`
			TBT        pageSpeedAudit `json:"total-blocking-time"`
			CLS        pageSpeedAudit `json:"cumulative-layout-shift"`
			SpeedIndex pageSpeedAudit `json:"speed-index"`
			TTI        pageSpeedAudit `json:"interactive"`
		} `json:"audits"`
		AuditsList []pageSpeedAuditItem `json:"-"` // populated separately
	} `json:"lighthouseResult"`
}

type pageSpeedAudit struct {
	NumericValue float64 `json:"numericValue"`
}

type pageSpeedAuditItem struct {
	Title string   `json:"title"`
	Score *float64 `json:"score"`
}
