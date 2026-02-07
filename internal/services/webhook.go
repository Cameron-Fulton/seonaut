package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/stjudewashere/seonaut/internal/models"
)

// WebhookServiceRepository defines the data access interface for the webhook service.
type WebhookServiceRepository interface {
	FindWebhooksByProjectId(projectId int64) []models.Webhook
	FindAllWebhooksByProjectId(projectId int64) []models.Webhook
	SaveWebhook(w *models.Webhook) error
	DeleteWebhook(id int64) error
	UpdateWebhookActive(id int64, active bool) error
	SaveWebhookLog(l *models.WebhookLog) error
	FindWebhookLogs(webhookId int64, limit int) []models.WebhookLog
}

// WebhookService manages webhook delivery and configuration.
type WebhookService struct {
	repository WebhookServiceRepository
	client     *http.Client
}

// NewWebhookService creates a new WebhookService.
func NewWebhookService(r WebhookServiceRepository) *WebhookService {
	return &WebhookService{
		repository: r,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send dispatches a webhook event to all matching active webhooks for a project.
func (s *WebhookService) Send(projectId int64, eventType string, data interface{}) {
	webhooks := s.repository.FindWebhooksByProjectId(projectId)

	payload := models.WebhookPayload{
		Event:     eventType,
		Timestamp: time.Now().UTC(),
		ProjectId: projectId,
		Data:      data,
	}

	for _, w := range webhooks {
		if w.EventType == eventType || w.EventType == "*" {
			go s.deliver(w, payload)
		}
	}
}

// deliver sends the webhook payload to the configured URL.
func (s *WebhookService) deliver(w models.Webhook, payload models.WebhookPayload) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("webhook marshal error: %v", err)
		return
	}

	req, err := http.NewRequest("POST", w.URL, bytes.NewReader(body))
	if err != nil {
		log.Printf("webhook request error: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SEOnaut-Webhook/1.0")

	// Sign the payload with HMAC if a secret is configured
	if w.Secret != "" {
		mac := hmac.New(sha256.New, []byte(w.Secret))
		mac.Write(body)
		signature := hex.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-Webhook-Signature", "sha256="+signature)
	}

	wl := models.WebhookLog{
		WebhookId: w.Id,
		EventType: payload.Event,
		Payload:   string(body),
	}

	resp, err := s.client.Do(req)
	if err != nil {
		wl.Error = err.Error()
		wl.ResponseCode = 0
	} else {
		wl.ResponseCode = resp.StatusCode
		resp.Body.Close()
	}

	if saveErr := s.repository.SaveWebhookLog(&wl); saveErr != nil {
		log.Printf("webhook log save error: %v", saveErr)
	}
}

// GetWebhooks returns all webhooks for a project.
func (s *WebhookService) GetWebhooks(projectId int64) []models.Webhook {
	return s.repository.FindAllWebhooksByProjectId(projectId)
}

// CreateWebhook creates a new webhook.
func (s *WebhookService) CreateWebhook(w *models.Webhook) error {
	return s.repository.SaveWebhook(w)
}

// DeleteWebhook removes a webhook.
func (s *WebhookService) DeleteWebhook(id int64) error {
	return s.repository.DeleteWebhook(id)
}

// ToggleWebhook enables or disables a webhook.
func (s *WebhookService) ToggleWebhook(id int64, active bool) error {
	return s.repository.UpdateWebhookActive(id, active)
}

// GetWebhookLogs returns recent delivery logs for a webhook.
func (s *WebhookService) GetWebhookLogs(webhookId int64) []models.WebhookLog {
	return s.repository.FindWebhookLogs(webhookId, 50)
}
