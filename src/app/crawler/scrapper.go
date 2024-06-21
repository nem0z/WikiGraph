package crawler

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/nem0z/WikiGraph/app/entity"
)

const WikiBaseUrl string = "https://fr.wikipedia.org/wiki/"

var validLinkPattern = regexp.MustCompile(`^/wiki/([^:]*)$`)

type InvalidUrl struct {
	URL string
}

func (e *InvalidUrl) Error() string {
	return fmt.Sprintf("invalid URL: %s", e.URL)
}

type URLError struct {
	URL string
}

func (e *URLError) Error() string {
	return fmt.Sprintf("invalid URL: %s", e.URL)
}

type Scraper struct {
	*colly.Collector
	url string
}

func NewScraper(url string) *Scraper {
	url = fmt.Sprintf("%v%v", WikiBaseUrl, url)
	return &Scraper{colly.NewCollector(), url}
}

func isValidLink(link string) (string, error) {
	matches := validLinkPattern.FindStringSubmatch(link)
	if len(matches) > 1 && matches[1] != "" {
		return matches[1], nil
	}
	return "", &InvalidUrl{link}
}

func (s *Scraper) GetArticles() (articles []*entity.Article, finalError error) {
	s.OnHTML("#mw-content-text a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.Attr("title")

		if url, err := isValidLink(link); err == nil {
			articles = append(articles, entity.NewArticle(url, title))
		}
	})

	s.OnError(func(r *colly.Response, err error) {
		if err != nil {
			formatedError := fmt.Sprintf("Scrapper on URL : %v failed with response: %v\nError : %v", r.Request.URL, r, err)
			finalError = errors.New(formatedError)
		}
	})

	err := s.Visit(s.url)
	if finalError != nil {
		finalError = err
	}

	return articles, finalError
}
