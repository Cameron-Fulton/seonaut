package page

import (
	"net/http"

	"github.com/stjudewashere/seonaut/internal/issues/errors"
	"github.com/stjudewashere/seonaut/internal/models"

	"golang.org/x/net/html"
)

// Returns a report_manager.PageIssueReporter with a callback function that returns true if
// the media type is text/html, the page has more than 200 words, and the readability score
// is below 30, indicating the content is very hard to read.
func NewLowReadabilityReporter() *models.PageIssueReporter {
	c := func(pageReport *models.PageReport, htmlNode *html.Node, header *http.Header) bool {
		if !pageReport.Crawled {
			return false
		}

		if pageReport.MediaType != "text/html" {
			return false
		}

		if pageReport.Words <= 200 {
			return false
		}

		return pageReport.ReadabilityScore < 30
	}

	return &models.PageIssueReporter{
		ErrorType: errors.ErrorLowReadability,
		Callback:  c,
	}
}
