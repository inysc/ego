package clientocol

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type EgoSQL interface {
	SQLNames() []string
	SQLValues() []any
}

var dbx *sqlx.DB

func InitDB(dsn string) *sqlx.DB {

	dbx = sqlx.MustConnect("mysql", dsn)

	return dbx
}

func CloseDB() {
	dbx.Close()
}

func Transaction(dbx *sqlx.DB, fc func(*sqlx.Tx) error) error {
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
