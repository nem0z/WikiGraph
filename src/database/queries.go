package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/nem0z/WikiGraph/entity"
)

const UniqueConstraintViolationErrCode = 1062

type Queryable interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func rollbackTransaction(tx *sql.Tx, identifier string) {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		log.Printf("error on tx rollback (%v) : %v\n", identifier, rollbackErr)
	}
}

func GetUnprocessedArticleLinks(db *DB, n uint) (links []string, err error) {
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

func ProcessArticle(q Queryable, id int64) error {
	const query string = "UPDATE articles SET processed = 1 WHERE id = ?"
	_, err := q.Exec(query, id)
	return err
}

func GetArticleId(q Queryable, link string) (id int64, err error) {
	const query string = "SELECT id FROM articles WHERE link = ?"
	return id, q.QueryRow(query, link).Scan(&id)
}

func GetOrCreateArticle(q Queryable, article *entity.Article) (id int64, err error) {
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

func CreateRelation(q Queryable, parentId int64, childArticle *entity.Article) error {
	const query string = "INSERT INTO relations (parent, child) VALUES (?, ?)"

	childId, err := GetOrCreateArticle(q, childArticle)
	if err != nil {
		return err
	}

	err = q.QueryRow(query, parentId, childId).Err()

	// handle unique constraint violation
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == UniqueConstraintViolationErrCode {
		return nil
	}

	return err
}

func CreateRelations(db *DB, relation *entity.Relation) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	parentId, err := GetArticleId(tx, relation.ParentLink)
	if err != nil {
		rollbackTransaction(tx, relation.ParentLink)
		return err
	}

	for _, child := range relation.Childs {
		err = CreateRelation(tx, parentId, child)
		if err != nil {
			log.Printf("Relation with child (%v) not created : %v\n", child.Link, err)
			rollbackTransaction(tx, relation.ParentLink)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return ProcessArticle(db, parentId)
}
