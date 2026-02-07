package page

import (
	"net/http"

	"github.com/antchfx/htmlquery"
	"github.com/stjudewashere/seonaut/internal/issues/errors"
	"github.com/stjudewashere/seonaut/internal/models"

	"golang.org/x/net/html"
)

// Returns a report_manager.PageIssueReporter with a callback function that returns true if
// the media type is text/html and the page has more than one H1 tag.
func NewMultipleH1TagsReporter() *models.PageIssueReporter {
	c := func(pageReport *models.PageReport, htmlNode *html.Node, header *http.Header) bool {
		if !pageReport.Crawled {
			return false
		}

		if pageReport.MediaType != "text/html" {
			return false
		}

		tags, err := htmlquery.QueryAll(htmlNode, "//h1")
		if err != nil {
			return false
		}

		return len(tags) > 1
	}

	return &models.PageIssueReporter{
		ErrorType: errors.ErrorMultipleH1Tags,
		Callback:  c,
	}
}

// Returns a report_manager.PageIssueReporter with a callback function that returns true if
// the media type is text/html and the page has more than one H2 tag.
func NewMultipleH2TagsReporter() *models.PageIssueReporter {
	c := func(pageReport *models.PageReport, htmlNode *html.Node, header *http.Header) bool {
		if !pageReport.Crawled {
			return false
		}

		if pageReport.MediaType != "text/html" {
			return false
		}

		tags, err := htmlquery.QueryAll(htmlNode, "//h2")
		if err != nil {
			return false
		}

		return len(tags) > 1
	}

	return &models.PageIssueReporter{
		ErrorType: errors.ErrorMultipleH2Tags,
		Callback:  c,
	}
}
