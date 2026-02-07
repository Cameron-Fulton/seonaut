package page

import (
	"net/http"

	"github.com/stjudewashere/seonaut/internal/issues/errors"
	"github.com/stjudewashere/seonaut/internal/models"

	"golang.org/x/net/html"
)

// Returns a report_manager.PageIssueReporter with a callback function that returns true if
// the media type is text/html and the page uses a meta refresh redirect.
func NewMetaRefreshRedirectReporter() *models.PageIssueReporter {
	c := func(pageReport *models.PageReport, htmlNode *html.Node, header *http.Header) bool {
		if !pageReport.Crawled {
			return false
		}

		if pageReport.MediaType != "text/html" {
			return false
		}

		return pageReport.Refresh != ""
	}

	return &models.PageIssueReporter{
		ErrorType: errors.ErrorMetaRefreshRedirect,
		Callback:  c,
	}
}
