package repository

import (
	"database/sql"

	"github.com/stjudewashere/seonaut/internal/models"
)

type IntegrationsRepository struct {
	DB *sql.DB
}

// SaveAPICredential stores or updates API credentials for a user/service.
func (r *IntegrationsRepository) SaveAPICredential(c *models.APICredential) error {
	query := `INSERT INTO api_credentials (user_id, service, credentials)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE credentials = VALUES(credentials)`

	result, err := r.DB.Exec(query, c.UserId, c.Service, c.Credentials)
	if err != nil {
		return err
	}

	c.Id, _ = result.LastInsertId()
	return nil
}

// FindAPICredential returns API credentials for a specific user and service.
func (r *IntegrationsRepository) FindAPICredential(userId int64, service string) (*models.APICredential, error) {
	query := `SELECT id, user_id, service, credentials, created, updated
		FROM api_credentials WHERE user_id = ? AND service = ?`

	c := &models.APICredential{}
	err := r.DB.QueryRow(query, userId, service).Scan(
		&c.Id, &c.UserId, &c.Service, &c.Credentials, &c.Created, &c.Updated,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// DeleteAPICredential removes API credentials.
func (r *IntegrationsRepository) DeleteAPICredential(userId int64, service string) error {
	_, err := r.DB.Exec("DELETE FROM api_credentials WHERE user_id = ? AND service = ?", userId, service)
	return err
}

// SavePageSpeedResult persists a PageSpeed Insights result.
func (r *IntegrationsRepository) SavePageSpeedResult(ps *models.PageSpeedResult) error {
	query := `INSERT INTO pagespeed_results
		(pagereport_id, crawl_id, performance_score, fcp, lcp, tbt, cls, speed_index, tti, recommendations)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(query,
		ps.PageReportId, ps.CrawlId, ps.PerformanceScore,
		ps.FCP, ps.LCP, ps.TBT, ps.CLS, ps.SpeedIndex, ps.TTI, ps.Recommendations,
	)
	if err != nil {
		return err
	}

	ps.Id, err = result.LastInsertId()
	return err
}

// FindPageSpeedResult returns PageSpeed results for a page report.
func (r *IntegrationsRepository) FindPageSpeedResult(pageReportId int64) (*models.PageSpeedResult, error) {
	query := `SELECT id, pagereport_id, crawl_id, performance_score, fcp, lcp, tbt, cls, speed_index, tti, recommendations, created
		FROM pagespeed_results WHERE pagereport_id = ? ORDER BY created DESC LIMIT 1`

	ps := &models.PageSpeedResult{}
	err := r.DB.QueryRow(query, pageReportId).Scan(
		&ps.Id, &ps.PageReportId, &ps.CrawlId,
		&ps.PerformanceScore, &ps.FCP, &ps.LCP, &ps.TBT, &ps.CLS,
		&ps.SpeedIndex, &ps.TTI, &ps.Recommendations, &ps.Created,
	)
	if err != nil {
		return nil, err
	}

	return ps, nil
}
