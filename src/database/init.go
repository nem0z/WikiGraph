package database

import (
	"database/sql"
	"os"
	"regexp"
)

const queryPattern = `\b(CREATE|INSERT)\b[\s\S]*?;`

func Init(db *sql.DB, path string) error {
	initScript, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(queryPattern)
	queries := re.FindAllString(string(initScript), -1)

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
