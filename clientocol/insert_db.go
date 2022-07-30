package clientocol

import (
	"fmt"
	"strings"
	"sync"
)

type dbxExec interface {
	Exec(string, ...any) error
}

type SQLInsert interface {
	Fields(eg EgoSQL) SQLInsert
	Values(vs ...EgoSQL) SQLInsert
	Insert(dbxExec) error
	Clear()
	String() string
}

var sqlinsertpl = sync.Pool{
	New: func() any { return &sqlinsert{} },
}

type sqlinsert struct {
	tablename string
	fields    string
	bindvars  string
	values    []any
}

var _ SQLInsert = &sqlinsert{}

func InsertInto(tablename string) SQLInsert {
	ret := sqlinsertpl.Get().(*sqlinsert)
	ret.Clear()
	ret.tablename = tablename

	return ret
}

func (si *sqlinsert) Fields(eg EgoSQL) SQLInsert {
	si.fields = strings.Join(eg.SQLNames(), " , ")
	bv := strings.Repeat(" ?,", len(eg.SQLNames()))
	si.bindvars = " ( " + bv[:len(bv)-1] + " ) "
	return si
}

func (si *sqlinsert) Values(vs ...EgoSQL) SQLInsert {
	for _, v := range vs {
		si.values = append(si.values, v.SQLValues()...)
	}
	si.bindvars = strings.Repeat(si.bindvars+",", len(vs))
	si.bindvars = si.bindvars[:len(si.bindvars)]

	return si
}

func (si *sqlinsert) Insert(dbx dbxExec) error {
	return dbx.Exec(si.String(), si.values...)
}

func (si *sqlinsert) Clear() {
	si.tablename = ""
	si.fields = ""
	si.bindvars = ""
	si.values = si.values[:0]
}

// insert into user(name, age) values (?,?)
func (si *sqlinsert) String() string {
	bs := &strings.Builder{}

	fmt.Fprintf(
		bs,
		"INSERT INTO %s (%s) VALUES %s ",
		si.tablename, si.fields, si.bindvars,
	)

	return bs.String()
}
