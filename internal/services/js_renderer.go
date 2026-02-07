package services

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/stjudewashere/seonaut/internal/models"
)

const (
	// jsRenderTimeout is the maximum time to wait for JS rendering.
	jsRenderTimeout = 30 * time.Second

	// jsWaitAfterLoad is how long to wait after page load for JS to execute.
	jsWaitAfterLoad = 2 * time.Second
)

// JSRenderer provides headless browser rendering capabilities.
type JSRenderer struct {
	allocCtx context.Context
	cancel   context.CancelFunc
}

// NewJSRenderer creates a new JSRenderer with a headless Chrome instance.
func NewJSRenderer() *JSRenderer {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-background-networking", true),
	)

	// Use system Chromium if CHROMIUM_PATH is set (e.g. in Docker)
	if chromiumPath := os.Getenv("CHROMIUM_PATH"); chromiumPath != "" {
		opts = append(opts, chromedp.ExecPath(chromiumPath))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	return &JSRenderer{
		allocCtx: allocCtx,
		cancel:   cancel,
	}
}

// Close shuts down the headless browser.
func (jr *JSRenderer) Close() {
	jr.cancel()
}

// RenderPage fetches a URL using a headless browser and returns the rendered HTML.
func (jr *JSRenderer) RenderPage(pageURL string, userAgent string) (string, error) {
	ctx, cancel := chromedp.NewContext(jr.allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, jsRenderTimeout)
	defer cancel()

	var renderedHTML string

	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(pageURL),
		chromedp.Sleep(jsWaitAfterLoad),
		chromedp.OuterHTML("html", &renderedHTML),
	)
	if err != nil {
		return "", err
	}

	return renderedHTML, nil
}

// JSRenderingResult holds the comparison between raw and rendered HTML.
type JSRenderingResult struct {
	RenderedHTML   string
	RawWordCount  int
	JSWordCount   int
	HasDifference bool
}

// CompareRenderedContent compares the original HTTP response with JS-rendered content.
func (jr *JSRenderer) CompareRenderedContent(pageURL string, rawPageReport *models.PageReport, userAgent string) (*JSRenderingResult, error) {
	renderedHTML, err := jr.RenderPage(pageURL, userAgent)
	if err != nil {
		return nil, err
	}

	// Parse the rendered HTML to count words
	renderedURL, _ := url.Parse(pageURL)
	renderedBody := []byte(renderedHTML)
	renderedReport, _, err := NewHTMLParser(
		renderedURL,
		200,
		&http.Header{"Content-Type": []string{"text/html; charset=utf-8"}},
		renderedBody,
		int64(len(renderedBody)),
	)
	if err != nil {
		return nil, err
	}

	result := &JSRenderingResult{
		RenderedHTML:  renderedHTML,
		RawWordCount: rawPageReport.Words,
		JSWordCount:  renderedReport.Words,
	}

	// Determine if there's a significant difference
	if result.RawWordCount == 0 && result.JSWordCount > 0 {
		result.HasDifference = true
	} else if result.RawWordCount > 0 {
		ratio := float64(result.JSWordCount) / float64(result.RawWordCount)
		// If JS rendering produces significantly different content (>20% change)
		result.HasDifference = ratio > 1.2 || ratio < 0.8
	}

	return result, nil
}
