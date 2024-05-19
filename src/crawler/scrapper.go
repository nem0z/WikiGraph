package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/nem0z/WikiGraph/app/article"
)

const WikiBaseUrl string = "https://fr.wikipedia.org/wiki/"

type Scraper struct {
	*colly.Collector
	url string
}

func NewScraper(url string) *Scraper {
	url = fmt.Sprintf("%v%v", WikiBaseUrl, url)
	return &Scraper{colly.NewCollector(), url}
}

func isValidLink(link string) bool {
	return strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":")
}

func formateUrl(link string) string {
	return strings.TrimPrefix(link, "/wiki/")
}

func (s *Scraper) GetArticles() (articles []*article.Article, finalError error) {
	s.OnHTML("#mw-content-text a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.Attr("title")

		if isValidLink(link) {
			link, err := url.QueryUnescape(link)
			if err != nil {
				finalError = err
			}

			articles = append(articles, article.NewArticle(formateUrl(link), title))
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
