package article

import (
	"database/sql"

	"github.com/nem0z/WikiGraph/database"
)

type Article struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

func NewArticle(link, title string) *Article {
	return &Article{Link: link, Title: title}
}

func GetIdFromLink(q database.Queryable, link string) (id int64, err error) {
	const query string = "SELECT id FROM articles WHERE link = ?"
	return id, q.QueryRow(query, link).Scan(&id)
}

func Process(q database.Queryable, id int64) error {
	const query string = "UPDATE articles SET processed = 1 WHERE id = ?"
	_, err := q.Exec(query, id)
	return err
}

func (article *Article) GetOrCreateDB(q database.Queryable) (id int64, err error) {
	const selectQuery string = "SELECT id FROM articles WHERE link = ?"
	const insertQuery string = "INSERT INTO articles (link, title) VALUES (?, ?) ON DUPLICATE KEY UPDATE title = VALUES(title);"

	err = q.QueryRow(selectQuery, article.Link).Scan(&id)

	if err == sql.ErrNoRows {
		res, err := q.Exec(insertQuery, article.Link, article.Title)
		if err != nil {
			if res == nil {
				return id, err
			}
			return res.RowsAffected()
		}

		id, err := res.LastInsertId()
		if id == 0 || err != nil {
			id, err = res.RowsAffected()
		}

		return id, err
	}

	return id, err
}

func GetUnprocessedArticleLinks(db *database.DB, n uint) (links []string, err error) {
	const query string = "SELECT link FROM articles WHERE processed = 0 LIMIT ?"

	rows, err := db.Query(query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var link string

		err = rows.Scan(&link)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, rows.Err()
}
