package database

import (
	"database/sql"
	"log"
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
