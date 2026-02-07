package models

import (
	"net/url"
)

type PageReport struct {
	Id                 int64
	URL                string
	ParsedURL          *url.URL
	RedirectURL        string
	Refresh            string
	StatusCode         int
	ContentType        string
	MediaType          string
	Lang               string
	Title              string
	Description        string
	Robots             string
	Noindex            bool
	Nofollow           bool
	Canonical          string
	H1                 string
	H2                 string
	Links              []Link
	ExternalLinks      []Link
	Words              int
	Hreflangs          []Hreflang
	Size               int64
	Images             []Image
	Scripts            []string
	Styles             []string
	Iframes            []string
	Audios             []string
	Videos             []Video
	BlockedByRobotstxt bool
	Crawled            bool
	InSitemap          bool
	InternalLinks      []InternalLink
	Depth              int
	BodyHash           string
	Timeout            bool
	TTFB               int

	// Agency custom fields
	MetaKeywords     string
	H1Count          int
	H2Count          int
	ReadabilityScore float64
	SimHash          string
	SchemaTypes      string // comma-separated list of detected schema.org types
	HasOpenGraph     bool
	HasTwitterCard   bool
	BodyText         string // extracted body text, not persisted to DB
}
