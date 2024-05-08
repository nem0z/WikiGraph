package entity

type Article struct {
	Id    int64
	Url   string
	Title string
}

func NewArticle(url, title string) *Article {
	return &Article{Url: url, Title: title}
}
