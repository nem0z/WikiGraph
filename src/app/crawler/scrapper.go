package crawler

import (
	"fmt"
	"regexp"
	"time"

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

type Scraper struct {
	*colly.Collector
	url string
}

func NewScraper(url string, proxy string) (*Scraper, error) {
	//url = urlpkg.QueryEscape(url)
	url = fmt.Sprintf("%v%v", WikiBaseUrl, url)

	collector := colly.NewCollector()
	collector.SetRequestTimeout(30 * time.Second)
	if proxy != "" {
		if err := collector.SetProxy(proxy); err != nil {
			return nil, err
		}
	}

	return &Scraper{collector, url}, nil
}

func isValidLink(link string) (string, error) {
	matches := validLinkPattern.FindStringSubmatch(link)
	if len(matches) > 1 && matches[1] != "" {
		return matches[1], nil
	}
	return "", &InvalidUrl{link}
}

func (s *Scraper) GetArticles() (articles []*entity.Article, err error) {
	s.OnHTML("#mw-content-text a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.Attr("title")

		if url, err := isValidLink(link); err == nil {
			articles = append(articles, entity.NewArticle(url, title))
		}
	})

	err = s.Visit(s.url)
	if err != nil {
		err = fmt.Errorf("GetArticles (%v) failed with error: %v", s.url, err)
	}

	return articles, err
}
