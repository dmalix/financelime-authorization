package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
)

func Migrate(db *sql.DB, dropFile string, createFile string, insertFile string) error {

	var (
		err  error
		body []byte
	)

	if len(dropFile) != 0 {
		body, err = ioutil.ReadFile(dropFile)
		if err != nil {
			return fmt.Errorf("failed to read the 'drop' file: %s", err)
		}
		_, err = db.Exec(string(body))
		if err != nil {
			return fmt.Errorf("failed to complete the request: %s", err)
		}
	}

	if len(createFile) != 0 {
		body, err = ioutil.ReadFile(createFile)
		if err != nil {
			return fmt.Errorf("failed to read the 'create' file: %s", err)
		}
		_, err = db.Exec(string(body))
		if err != nil {
			return fmt.Errorf("failed to complete the request: %s", err)
		}
	}

	if len(insertFile) != 0 {
		body, err = ioutil.ReadFile(insertFile)
		if err != nil {
			return fmt.Errorf("failed to read the 'insert' file: %s", err)
		}
		_, err = db.Exec(string(body))
		if err != nil {
			return fmt.Errorf("failed to complete the request: %s", err)
		}
	}

	return nil
}
