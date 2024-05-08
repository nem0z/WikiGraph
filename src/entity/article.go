package entity

type Article struct {
	Id    int64  `json:"id"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

func NewArticle(url, title string) *Article {
	return &Article{Url: url, Title: title}
}
