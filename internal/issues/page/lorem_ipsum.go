package page

import (
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/stjudewashere/seonaut/internal/issues/errors"
	"github.com/stjudewashere/seonaut/internal/models"

	"golang.org/x/net/html"
)

// Returns a report_manager.PageIssueReporter with a callback function that returns true if
// the media type is text/html and the page's body contains "lorem ipsum" placeholder text.
func NewLoremIpsumReporter() *models.PageIssueReporter {
	c := func(pageReport *models.PageReport, htmlNode *html.Node, header *http.Header) bool {
		if !pageReport.Crawled {
			return false
		}

		if pageReport.MediaType != "text/html" {
			return false
		}

		body, err := htmlquery.Query(htmlNode, "//body")
		if err != nil || body == nil {
			return false
		}

		bodyText := htmlquery.InnerText(body)

		return strings.Contains(strings.ToLower(bodyText), "lorem ipsum")
	}

	return &models.PageIssueReporter{
		ErrorType: errors.ErrorLoremIpsum,
		Callback:  c,
	}
}
