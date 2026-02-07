package repository

import (
	"database/sql"
	"log"

	"github.com/stjudewashere/seonaut/internal/models"
)

type AIRepository struct {
	DB *sql.DB
}

// SaveAIAnalysis persists an AI analysis result to the database.
func (r *AIRepository) SaveAIAnalysis(a *models.AIAnalysis) error {
	query := `INSERT INTO ai_analyses (pagereport_id, crawl_id, provider, analysis_type, result)
		VALUES (?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(query, a.PageReportId, a.CrawlId, a.Provider, a.AnalysisType, a.Result)
	if err != nil {
		return err
	}

	a.Id, err = result.LastInsertId()
	return err
}

// FindAIAnalysesByPageReport returns all AI analyses for a specific page report.
func (r *AIRepository) FindAIAnalysesByPageReport(pageReportId int64) []models.AIAnalysis {
	analyses := []models.AIAnalysis{}
	query := `SELECT id, pagereport_id, crawl_id, provider, analysis_type, result, created
		FROM ai_analyses WHERE pagereport_id = ? ORDER BY created DESC`

	rows, err := r.DB.Query(query, pageReportId)
	if err != nil {
		log.Printf("FindAIAnalysesByPageReport: %v", err)
		return analyses
	}
	defer rows.Close()

	for rows.Next() {
		a := models.AIAnalysis{}
		err := rows.Scan(&a.Id, &a.PageReportId, &a.CrawlId, &a.Provider, &a.AnalysisType, &a.Result, &a.Created)
		if err != nil {
			log.Printf("FindAIAnalysesByPageReport scan: %v", err)
			continue
		}
		analyses = append(analyses, a)
	}

	return analyses
}
