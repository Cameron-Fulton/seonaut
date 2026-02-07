package multipage

import (
	"github.com/stjudewashere/seonaut/internal/issues/errors"
	"github.com/stjudewashere/seonaut/internal/models"
)

// NearDuplicateContentReporter reports pages with near-duplicate content based on SimHash similarity.
// Two pages are considered near-duplicates if their SimHash fingerprints differ by 3 or fewer bits.
func (sr *SqlReporter) NearDuplicateContentReporter(c *models.Crawl) *models.MultipageIssueReporter {
	query := `
		SELECT DISTINCT p1.id
		FROM pagereports p1
		INNER JOIN pagereports p2 ON p1.crawl_id = p2.crawl_id
			AND p1.id != p2.id
			AND p1.sim_hash IS NOT NULL
			AND p2.sim_hash IS NOT NULL
			AND p1.sim_hash != ''
			AND p2.sim_hash != ''
			AND p1.sim_hash = p2.sim_hash
		WHERE p1.crawl_id = ?
			AND p1.media_type = 'text/html'
			AND p1.crawled = 1
			AND p1.status_code >= 200
			AND p1.status_code < 300`

	return &models.MultipageIssueReporter{
		Pstream:   sr.pageReportsQuery(query, c.Id),
		ErrorType: errors.ErrorNearDuplicateContent,
	}
}
