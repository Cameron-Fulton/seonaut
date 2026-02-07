package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/stjudewashere/seonaut/internal/models"
)

// AIServiceRepository defines the data access interface for AI analysis results.
type AIServiceRepository interface {
	SaveAIAnalysis(a *models.AIAnalysis) error
	FindAIAnalysesByPageReport(pageReportId int64) []models.AIAnalysis
}

// AIConfig holds the configuration for AI providers.
type AIConfig struct {
	ClaudeAPIKey   string
	ClaudeModel    string
	LMStudioURL    string
	LMStudioModel  string
	Enabled        bool
}

// AIService manages AI-powered content analysis.
type AIService struct {
	repository AIServiceRepository
	config     *AIConfig
	client     *http.Client
}

// NewAIService creates a new AIService.
func NewAIService(r AIServiceRepository, config *AIConfig) *AIService {
	return &AIService{
		repository: r,
		config:     config,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// AnalyzePageIntent uses AI to classify page intent (Informational, Commercial, Navigational, Transactional).
func (s *AIService) AnalyzePageIntent(pageReport *models.PageReport, crawlId int64) (*models.AIAnalysis, error) {
	if !s.config.Enabled {
		return nil, fmt.Errorf("AI service is not enabled")
	}

	prompt := fmt.Sprintf(
		"Analyze this webpage and classify its search intent as one of: Informational, Commercial, Navigational, or Transactional. "+
			"Respond with ONLY a JSON object like {\"intent\": \"...\", \"confidence\": 0.0-1.0, \"reasoning\": \"...\"}.\n\n"+
			"URL: %s\nTitle: %s\nDescription: %s\nH1: %s\nWord Count: %d",
		pageReport.URL, pageReport.Title, pageReport.Description, pageReport.H1, pageReport.Words,
	)

	result, err := s.callAI(prompt)
	if err != nil {
		return nil, err
	}

	analysis := &models.AIAnalysis{
		PageReportId: pageReport.Id,
		CrawlId:      crawlId,
		Provider:     s.getActiveProvider(),
		AnalysisType: models.AIAnalysisIntent,
		Result:       result,
	}

	if err := s.repository.SaveAIAnalysis(analysis); err != nil {
		log.Printf("AI analysis save error: %v", err)
	}

	return analysis, nil
}

// AnalyzeContentQuality uses AI to assess content quality and provide suggestions.
func (s *AIService) AnalyzeContentQuality(pageReport *models.PageReport, crawlId int64) (*models.AIAnalysis, error) {
	if !s.config.Enabled {
		return nil, fmt.Errorf("AI service is not enabled")
	}

	prompt := fmt.Sprintf(
		"Analyze this webpage's SEO content quality. Respond with ONLY a JSON object like "+
			"{\"score\": 0-100, \"strengths\": [...], \"improvements\": [...], \"keyword_suggestions\": [...]}.\n\n"+
			"URL: %s\nTitle: %s\nDescription: %s\nH1: %s\nH2: %s\nWord Count: %d\nMeta Keywords: %s",
		pageReport.URL, pageReport.Title, pageReport.Description,
		pageReport.H1, pageReport.H2, pageReport.Words, pageReport.MetaKeywords,
	)

	result, err := s.callAI(prompt)
	if err != nil {
		return nil, err
	}

	analysis := &models.AIAnalysis{
		PageReportId: pageReport.Id,
		CrawlId:      crawlId,
		Provider:     s.getActiveProvider(),
		AnalysisType: models.AIAnalysisContentQuality,
		Result:       result,
	}

	if err := s.repository.SaveAIAnalysis(analysis); err != nil {
		log.Printf("AI analysis save error: %v", err)
	}

	return analysis, nil
}

// getActiveProvider returns the currently configured AI provider name.
func (s *AIService) getActiveProvider() string {
	if s.config.ClaudeAPIKey != "" {
		return models.AIProviderClaude
	}
	return models.AIProviderLMStudio
}

// callAI routes the request to the appropriate AI provider.
func (s *AIService) callAI(prompt string) (string, error) {
	if s.config.ClaudeAPIKey != "" {
		return s.callClaude(prompt)
	}
	if s.config.LMStudioURL != "" {
		return s.callLMStudio(prompt)
	}
	return "", fmt.Errorf("no AI provider configured")
}

// claudeRequest is the request body for the Claude API.
type claudeRequest struct {
	Model     string           `json:"model"`
	MaxTokens int              `json:"max_tokens"`
	Messages  []claudeMessage  `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// claudeResponse is the response body from the Claude API.
type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// callClaude sends a request to the Anthropic Claude API.
func (s *AIService) callClaude(prompt string) (string, error) {
	model := s.config.ClaudeModel
	if model == "" {
		model = "claude-sonnet-4-5-20250929"
	}

	reqBody := claudeRequest{
		Model:     model,
		MaxTokens: 1024,
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.config.ClaudeAPIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("claude API error: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(respBody, &claudeResp); err != nil {
		return "", fmt.Errorf("claude response parse error: %w", err)
	}

	if claudeResp.Error != nil {
		return "", fmt.Errorf("claude API error: %s", claudeResp.Error.Message)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("claude returned empty response")
	}

	return claudeResp.Content[0].Text, nil
}

// openAIRequest is the request body for OpenAI-compatible APIs (LM Studio).
type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIResponse is the response body from OpenAI-compatible APIs.
type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// callLMStudio sends a request to an LM Studio instance (OpenAI-compatible API).
func (s *AIService) callLMStudio(prompt string) (string, error) {
	baseURL := strings.TrimRight(s.config.LMStudioURL, "/")
	model := s.config.LMStudioModel
	if model == "" {
		model = "default"
	}

	reqBody := openAIRequest{
		Model: model,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("LM Studio API error: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(respBody, &openAIResp); err != nil {
		return "", fmt.Errorf("LM Studio response parse error: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("LM Studio API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("LM Studio returned empty response")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// GetAnalyses returns all AI analyses for a page report.
func (s *AIService) GetAnalyses(pageReportId int64) []models.AIAnalysis {
	return s.repository.FindAIAnalysesByPageReport(pageReportId)
}
