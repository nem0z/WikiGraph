package entity

type Article struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

func NewArticle(link, title string) *Article {
	return &Article{Link: link, Title: title}
}
