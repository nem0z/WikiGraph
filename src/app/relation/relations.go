package relation

import (
	"github.com/go-sql-driver/mysql"
	"github.com/nem0z/WikiGraph/app/article"
	"github.com/nem0z/WikiGraph/database"
)

type Relation struct {
	ParentLink string             `json:"parentLink"`
	Childs     []*article.Article `json:"childs"`
}

func NewRelation(parentLink string, childs ...*article.Article) *Relation {
	return &Relation{ParentLink: parentLink, Childs: childs}
}

func CreateDB(q database.Queryable, parentId int64, childArticle *article.Article) error {
	const query string = "INSERT INTO relations (parent, child) VALUES (?, ?)"

	childId, err := childArticle.GetOrCreateDB(q)
	if err != nil {
		return err
	}

	err = q.QueryRow(query, parentId, childId).Err()

	// handle unique constraint violation
	if mysqlErr, ok := err.(*mysql.MySQLError); ok &&
		mysqlErr.Number == database.UniqueConstraintViolationErrCode {
		return nil
	}

	return err
}

func (relation *Relation) Create(db *database.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	parentId, err := article.GetIdFromLink(tx, relation.ParentLink)
	if err != nil {
		database.RollbackTransaction(tx, relation.ParentLink)
		return err
	}

	for _, child := range relation.Childs {
		err = CreateDB(tx, parentId, child)
		if err != nil {
			database.RollbackTransaction(tx, relation.ParentLink)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return article.Process(db, parentId)
}
