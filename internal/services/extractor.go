package services

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/stjudewashere/seonaut/internal/models"
	"golang.org/x/net/html"
)

// ExtractorServiceRepository defines the data access interface for the extractor service.
type ExtractorServiceRepository interface {
	FindExtractorsByProjectId(projectId int64) []models.CustomExtractor
	SaveExtractor(e *models.CustomExtractor) error
	DeleteExtractor(id int64) error
	SaveExtractionResult(result *models.ExtractionResult) error
	FindExtractionResultsByPageReport(pageReportId int64) []models.ExtractionResult
}

// ExtractorService manages custom data extraction rules.
type ExtractorService struct {
	repository ExtractorServiceRepository
}

// NewExtractorService creates a new ExtractorService.
func NewExtractorService(r ExtractorServiceRepository) *ExtractorService {
	return &ExtractorService{
		repository: r,
	}
}

// RunExtractors applies all project extractors to an HTML document and saves results.
func (s *ExtractorService) RunExtractors(projectId int64, pageReport *models.PageReport, htmlNode *html.Node, crawlId int64, bodyText string) {
	extractors := s.repository.FindExtractorsByProjectId(projectId)

	for _, extractor := range extractors {
		value, err := s.extract(extractor, htmlNode, bodyText)
		if err != nil {
			log.Printf("extractor %s error: %v", extractor.Name, err)
			continue
		}

		if value == "" {
			continue
		}

		result := &models.ExtractionResult{
			PageReportId: pageReport.Id,
			CrawlId:      crawlId,
			ExtractorId:  extractor.Id,
			Value:        value,
		}

		if err := s.repository.SaveExtractionResult(result); err != nil {
			log.Printf("save extraction result error: %v", err)
		}
	}
}

// extract applies a single extractor to the content.
func (s *ExtractorService) extract(e models.CustomExtractor, htmlNode *html.Node, bodyText string) (string, error) {
	switch e.Type {
	case models.ExtractorTypeRegex:
		return s.extractRegex(e.Pattern, bodyText)
	case models.ExtractorTypeCSS:
		return s.extractCSS(e.Pattern, htmlNode)
	case models.ExtractorTypeXPath:
		return s.extractXPath(e.Pattern, htmlNode)
	default:
		return "", fmt.Errorf("unknown extractor type: %s", e.Type)
	}
}

// extractRegex applies a regex pattern to the body text.
func (s *ExtractorService) extractRegex(pattern string, text string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern: %w", err)
	}

	matches := re.FindAllString(text, 10)
	if len(matches) == 0 {
		return "", nil
	}

	return strings.Join(matches, " | "), nil
}

// extractCSS converts a simple CSS selector to XPath and extracts content.
// Supports basic selectors: tag, .class, #id, tag.class, tag#id
func (s *ExtractorService) extractCSS(selector string, htmlNode *html.Node) (string, error) {
	xpath := cssToXPath(selector)
	return s.extractXPath(xpath, htmlNode)
}

// extractXPath applies an XPath expression to the HTML node.
func (s *ExtractorService) extractXPath(expr string, htmlNode *html.Node) (string, error) {
	nodes, err := htmlquery.QueryAll(htmlNode, expr)
	if err != nil {
		return "", fmt.Errorf("xpath error: %w", err)
	}

	var results []string
	for _, n := range nodes {
		text := strings.TrimSpace(htmlquery.InnerText(n))
		if text != "" {
			results = append(results, text)
		}
	}

	if len(results) == 0 {
		return "", nil
	}

	return strings.Join(results, " | "), nil
}

// cssToXPath converts basic CSS selectors to XPath expressions.
func cssToXPath(selector string) string {
	selector = strings.TrimSpace(selector)

	// Handle #id
	if strings.HasPrefix(selector, "#") {
		id := selector[1:]
		return fmt.Sprintf("//*[@id='%s']", id)
	}

	// Handle .class
	if strings.HasPrefix(selector, ".") {
		class := selector[1:]
		return fmt.Sprintf("//*[contains(@class, '%s')]", class)
	}

	// Handle tag#id
	if idx := strings.Index(selector, "#"); idx > 0 {
		tag := selector[:idx]
		id := selector[idx+1:]
		return fmt.Sprintf("//%s[@id='%s']", tag, id)
	}

	// Handle tag.class
	if idx := strings.Index(selector, "."); idx > 0 {
		tag := selector[:idx]
		class := selector[idx+1:]
		return fmt.Sprintf("//%s[contains(@class, '%s')]", tag, class)
	}

	// Handle tag[attr]
	if idx := strings.Index(selector, "["); idx > 0 {
		tag := selector[:idx]
		attr := strings.Trim(selector[idx:], "[]")
		if eqIdx := strings.Index(attr, "="); eqIdx > 0 {
			key := attr[:eqIdx]
			val := strings.Trim(attr[eqIdx+1:], "\"'")
			return fmt.Sprintf("//%s[@%s='%s']", tag, key, val)
		}
		return fmt.Sprintf("//%s[@%s]", tag, attr)
	}

	// Handle plain tag name
	return fmt.Sprintf("//%s", selector)
}

// GetExtractors returns all extractors for a project.
func (s *ExtractorService) GetExtractors(projectId int64) []models.CustomExtractor {
	return s.repository.FindExtractorsByProjectId(projectId)
}

// CreateExtractor creates a new custom extractor.
func (s *ExtractorService) CreateExtractor(e *models.CustomExtractor) error {
	// Validate the pattern based on type
	switch e.Type {
	case models.ExtractorTypeRegex:
		if _, err := regexp.Compile(e.Pattern); err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
	case models.ExtractorTypeCSS, models.ExtractorTypeXPath:
		// Basic validation - will fail at extraction time if invalid
	default:
		return fmt.Errorf("invalid extractor type: %s", e.Type)
	}

	return s.repository.SaveExtractor(e)
}

// DeleteExtractor removes an extractor.
func (s *ExtractorService) DeleteExtractor(id int64) error {
	return s.repository.DeleteExtractor(id)
}

// GetExtractionResults returns extraction results for a page report.
func (s *ExtractorService) GetExtractionResults(pageReportId int64) []models.ExtractionResult {
	return s.repository.FindExtractionResultsByPageReport(pageReportId)
}
