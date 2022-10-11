package clientocol

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

type EgoSQL interface {
	SQLNames() []string
	SQLValues() []any
}

var dbx *sqlx.DB

func InitDB(dsn string) *sqlx.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	dbx = sqlx.NewDb(db, "mysql")
	dbx.Mapper = reflectx.NewMapperFunc("SQL", strings.ToLower)

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
