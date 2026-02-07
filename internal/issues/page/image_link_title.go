package page

import (
	"net/http"

	"github.com/antchfx/htmlquery"
	"github.com/stjudewashere/seonaut/internal/issues/errors"
	"github.com/stjudewashere/seonaut/internal/models"

	"golang.org/x/net/html"
)

// Returns a report_manager.PageIssueReporter with a callback function that returns true if
// the media type is text/html and the page contains anchor elements with image children
// that are missing the title attribute on the anchor tag.
func NewImageLinkMissingTitleReporter() *models.PageIssueReporter {
	c := func(pageReport *models.PageReport, htmlNode *html.Node, header *http.Header) bool {
		if !pageReport.Crawled {
			return false
		}

		if pageReport.MediaType != "text/html" {
			return false
		}

		imageLinks, err := htmlquery.QueryAll(htmlNode, "//a[.//img]")
		if err != nil {
			return false
		}

		for _, link := range imageLinks {
			title := htmlquery.SelectAttr(link, "title")
			if title == "" {
				return true
			}
		}

		return false
	}

	return &models.PageIssueReporter{
		ErrorType: errors.ErrorImageLinkMissingTitle,
		Callback:  c,
	}
}
