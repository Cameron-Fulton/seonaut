package repository

import (
	"database/sql"
	"log"

	"github.com/stjudewashere/seonaut/internal/models"
)

type ExtractorRepository struct {
	DB *sql.DB
}

// FindExtractorsByProjectId returns all extractors for a project.
func (r *ExtractorRepository) FindExtractorsByProjectId(projectId int64) []models.CustomExtractor {
	extractors := []models.CustomExtractor{}
	query := `SELECT id, project_id, name, type, pattern, created FROM custom_extractors WHERE project_id = ?`

	rows, err := r.DB.Query(query, projectId)
	if err != nil {
		log.Printf("FindExtractorsByProjectId: %v", err)
		return extractors
	}
	defer rows.Close()

	for rows.Next() {
		e := models.CustomExtractor{}
		err := rows.Scan(&e.Id, &e.ProjectId, &e.Name, &e.Type, &e.Pattern, &e.Created)
		if err != nil {
			log.Printf("FindExtractorsByProjectId scan: %v", err)
			continue
		}
		extractors = append(extractors, e)
	}

	return extractors
}

// SaveExtractor creates a new custom extractor.
func (r *ExtractorRepository) SaveExtractor(e *models.CustomExtractor) error {
	query := `INSERT INTO custom_extractors (project_id, name, type, pattern) VALUES (?, ?, ?, ?)`
	result, err := r.DB.Exec(query, e.ProjectId, e.Name, e.Type, e.Pattern)
	if err != nil {
		return err
	}

	e.Id, err = result.LastInsertId()
	return err
}

// DeleteExtractor removes an extractor by ID.
func (r *ExtractorRepository) DeleteExtractor(id int64) error {
	_, err := r.DB.Exec("DELETE FROM custom_extractors WHERE id = ?", id)
	return err
}

// SaveExtractionResult persists an extraction result.
func (r *ExtractorRepository) SaveExtractionResult(result *models.ExtractionResult) error {
	query := `INSERT INTO extraction_results (pagereport_id, crawl_id, extractor_id, value)
		VALUES (?, ?, ?, ?)`
	res, err := r.DB.Exec(query, result.PageReportId, result.CrawlId, result.ExtractorId, result.Value)
	if err != nil {
		return err
	}

	result.Id, err = res.LastInsertId()
	return err
}

// FindExtractionResultsByPageReport returns extraction results for a page report.
func (r *ExtractorRepository) FindExtractionResultsByPageReport(pageReportId int64) []models.ExtractionResult {
	results := []models.ExtractionResult{}
	query := `SELECT id, pagereport_id, crawl_id, extractor_id, value
		FROM extraction_results WHERE pagereport_id = ?`

	rows, err := r.DB.Query(query, pageReportId)
	if err != nil {
		log.Printf("FindExtractionResultsByPageReport: %v", err)
		return results
	}
	defer rows.Close()

	for rows.Next() {
		er := models.ExtractionResult{}
		err := rows.Scan(&er.Id, &er.PageReportId, &er.CrawlId, &er.ExtractorId, &er.Value)
		if err != nil {
			log.Printf("FindExtractionResultsByPageReport scan: %v", err)
			continue
		}
		results = append(results, er)
	}

	return results
}
