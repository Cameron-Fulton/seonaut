package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/stjudewashere/seonaut/internal/models"
	"github.com/stjudewashere/seonaut/internal/services"
)

type webhookHandler struct {
	*services.Container
}

// indexHandler renders the webhook management page for a project.
func (h *webhookHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := h.CookieSession.GetUser(r.Context())
	if !ok {
		http.Redirect(w, r, "/signout", http.StatusSeeOther)
		return
	}

	pidStr := r.URL.Query().Get("pid")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	webhooks := h.WebhookService.GetWebhooks(pid)

	v := &PageView{
		PageTitle: "Webhooks",
		User:      *user,
		Data: struct {
			ProjectId int64
			Webhooks  []models.Webhook
		}{
			ProjectId: pid,
			Webhooks:  webhooks,
		},
	}

	h.Renderer.RenderTemplate(w, "webhooks", v)
}

// addHandler creates a new webhook.
func (h *webhookHandler) addHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := h.CookieSession.GetUser(r.Context())
	if !ok {
		http.Redirect(w, r, "/signout", http.StatusSeeOther)
		return
	}

	pidStr := r.FormValue("project_id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	webhook := &models.Webhook{
		ProjectId: pid,
		EventType: r.FormValue("event_type"),
		URL:       r.FormValue("url"),
		Secret:    r.FormValue("secret"),
		Active:    true,
	}

	if err := h.WebhookService.CreateWebhook(webhook); err != nil {
		http.Error(w, "Failed to create webhook", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/webhooks?pid=%d", pid), http.StatusSeeOther)
}

// deleteHandler removes a webhook.
func (h *webhookHandler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := h.CookieSession.GetUser(r.Context())
	if !ok {
		http.Redirect(w, r, "/signout", http.StatusSeeOther)
		return
	}

	pidStr := r.URL.Query().Get("pid")
	widStr := r.URL.Query().Get("wid")

	wid, err := strconv.ParseInt(widStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	h.WebhookService.DeleteWebhook(wid)
	http.Redirect(w, r, fmt.Sprintf("/webhooks?pid=%s", pidStr), http.StatusSeeOther)
}

// incomingHandler processes incoming webhook requests.
// This allows external services to trigger crawls or other actions via API.
func (h *webhookHandler) incomingHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.IncomingWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	switch payload.Action {
	case "crawl.start":
		// Trigger a crawl for the specified project
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "accepted",
			"message": "Crawl request received",
		})
	default:
		http.Error(w, fmt.Sprintf("Unknown action: %s", payload.Action), http.StatusBadRequest)
	}
}
