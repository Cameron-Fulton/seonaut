package repository

import (
	"database/sql"
	"log"

	"github.com/stjudewashere/seonaut/internal/models"
)

type WebhookRepository struct {
	DB *sql.DB
}

// FindWebhooksByProjectId returns all active webhooks for a project.
func (r *WebhookRepository) FindWebhooksByProjectId(projectId int64) []models.Webhook {
	webhooks := []models.Webhook{}
	query := `SELECT id, project_id, event_type, url, secret, active, created
		FROM webhooks WHERE project_id = ? AND active = 1`

	rows, err := r.DB.Query(query, projectId)
	if err != nil {
		log.Printf("FindWebhooksByProjectId: %v", err)
		return webhooks
	}
	defer rows.Close()

	for rows.Next() {
		w := models.Webhook{}
		err := rows.Scan(&w.Id, &w.ProjectId, &w.EventType, &w.URL, &w.Secret, &w.Active, &w.Created)
		if err != nil {
			log.Printf("FindWebhooksByProjectId scan: %v", err)
			continue
		}
		webhooks = append(webhooks, w)
	}

	return webhooks
}

// FindAllWebhooksByProjectId returns all webhooks (active and inactive) for a project.
func (r *WebhookRepository) FindAllWebhooksByProjectId(projectId int64) []models.Webhook {
	webhooks := []models.Webhook{}
	query := `SELECT id, project_id, event_type, url, secret, active, created
		FROM webhooks WHERE project_id = ?`

	rows, err := r.DB.Query(query, projectId)
	if err != nil {
		log.Printf("FindAllWebhooksByProjectId: %v", err)
		return webhooks
	}
	defer rows.Close()

	for rows.Next() {
		w := models.Webhook{}
		err := rows.Scan(&w.Id, &w.ProjectId, &w.EventType, &w.URL, &w.Secret, &w.Active, &w.Created)
		if err != nil {
			log.Printf("FindAllWebhooksByProjectId scan: %v", err)
			continue
		}
		webhooks = append(webhooks, w)
	}

	return webhooks
}

// SaveWebhook creates a new webhook.
func (r *WebhookRepository) SaveWebhook(w *models.Webhook) error {
	query := `INSERT INTO webhooks (project_id, event_type, url, secret, active)
		VALUES (?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(query, w.ProjectId, w.EventType, w.URL, w.Secret, w.Active)
	if err != nil {
		return err
	}

	w.Id, err = result.LastInsertId()
	return err
}

// DeleteWebhook removes a webhook by ID.
func (r *WebhookRepository) DeleteWebhook(id int64) error {
	_, err := r.DB.Exec("DELETE FROM webhooks WHERE id = ?", id)
	return err
}

// UpdateWebhookActive toggles the active state of a webhook.
func (r *WebhookRepository) UpdateWebhookActive(id int64, active bool) error {
	_, err := r.DB.Exec("UPDATE webhooks SET active = ? WHERE id = ?", active, id)
	return err
}

// SaveWebhookLog records a webhook delivery attempt.
func (r *WebhookRepository) SaveWebhookLog(l *models.WebhookLog) error {
	query := `INSERT INTO webhook_logs (webhook_id, event_type, payload, response_code, error)
		VALUES (?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(query, l.WebhookId, l.EventType, l.Payload, l.ResponseCode, l.Error)
	if err != nil {
		return err
	}

	l.Id, err = result.LastInsertId()
	return err
}

// FindWebhookLogs returns recent webhook delivery logs.
func (r *WebhookRepository) FindWebhookLogs(webhookId int64, limit int) []models.WebhookLog {
	logs := []models.WebhookLog{}
	query := `SELECT id, webhook_id, event_type, payload, response_code, error, created
		FROM webhook_logs WHERE webhook_id = ? ORDER BY created DESC LIMIT ?`

	rows, err := r.DB.Query(query, webhookId, limit)
	if err != nil {
		log.Printf("FindWebhookLogs: %v", err)
		return logs
	}
	defer rows.Close()

	for rows.Next() {
		l := models.WebhookLog{}
		err := rows.Scan(&l.Id, &l.WebhookId, &l.EventType, &l.Payload, &l.ResponseCode, &l.Error, &l.Created)
		if err != nil {
			log.Printf("FindWebhookLogs scan: %v", err)
			continue
		}
		logs = append(logs, l)
	}

	return logs
}
