package sqlite_repo

import (
	"os"

	"github.com/jmoiron/sqlx"
)

func CreateTebles(db *sqlx.DB, scriptPath string) error {
	query, err := os.ReadFile(scriptPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(query))

	return err
}
