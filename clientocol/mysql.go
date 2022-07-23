package clientocol

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var dbx *sqlx.DB

func InitDB(dsn string) *sqlx.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	dbx = sqlx.NewDb(db, "SQL")

	return dbx
}

func CloseDB() {
	dbx.Close()
}

func Transaction(fc func(*sqlx.Tx) error) error {
	tx, err := dbx.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fc(tx)

	return nil
}
