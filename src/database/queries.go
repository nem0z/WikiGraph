package database

import (
	"database/sql"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/nem0z/WikiGraph/app/entity"
)

const UniqueConstraintViolationErrCode = 1062

type Queryable interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func RollbackTransaction(tx *sql.Tx, identifier string) {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		log.Printf("error on tx rollback (%v) : %v\n", identifier, rollbackErr)
	}
}

func (db *DB) GetIdFromLink(link string) (id int64, err error) {
	id, err = db.cache.GetInt64(link)
	if err == nil && id >= 0 {
		return id, nil
	}

	const query string = "SELECT id FROM articles WHERE link = ?"
	err = db.QueryRow(query, link).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, db.cache.Set(link, id)
}

func (db *DB) CreateArticle(article *entity.Article) (id int64, err error) {
	const insertQuery string = "INSERT INTO articles (link, title) VALUES (?, ?)"

	res, err := db.Exec(insertQuery, article.Link, article.Title)
	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()
	if err != nil || id <= 0 {
		return id, err
	}

	if db.cache.Set(article.Link, id) != nil {
		log.Printf("Error inserting to cache : %v => %v\n", article.Link, id)
	}

	go db.onInsertArticle(article)

	return id, nil
}

func (db *DB) GetOrCreateArticle(article *entity.Article) (id int64, err error) {

	id, err = db.GetIdFromLink(article.Link)
	if err == sql.ErrNoRows {
		return db.CreateArticle(article)
	}

	return id, err
}

func (db *DB) ResolveArticleIds(articles ...*entity.Article) (ids []int64, finalErr error) {
	ids = make([]int64, len(articles))

	for i, article := range articles {
		id, err := db.GetOrCreateArticle(article)
		if err != nil {
			finalErr = err
			continue
		}

		ids[i] = id
	}

	return ids, finalErr
}

func (db *DB) CreateLinks(parentId int64, childIds ...int64) error {
	if len(childIds) == 0 {
		return nil
	}

	const baseQuery string = "INSERT IGNORE INTO relations (parent, child) VALUES "
	values := make([]string, len(childIds))
	args := make([]interface{}, 2*len(childIds))

	for i, childId := range childIds {
		values[i] = "(?, ?)"
		args[2*i] = parentId
		args[2*i+1] = childId
	}

	query := baseQuery + strings.Join(values, ",")
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok &&
		mysqlErr.Number == UniqueConstraintViolationErrCode {
		return nil
	}

	return err
}

func (db *DB) ProcessArticle(id int64) error {
	const query string = "UPDATE articles SET processed = 1 WHERE id = ?"
	_, err := db.Exec(query, id)
	return err
}

func (db *DB) GetUnprocessedArticleLinks(n uint) (links []string, err error) {
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
