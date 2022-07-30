package clientocol

import (
	"fmt"
	"strings"
	"sync"
)

type dbxGetSelect interface {
	Get(any, string, ...any) error
	Select(any, string, ...any) error
}

type SQLSelect interface {
	From(string) SQLSelect
	Join(string) SQLSelect
	On(string) SQLSelect
	Where(map[string]any) SQLSelect
	Limit(int) SQLSelect
	Offset(int) SQLSelect
	Group(string) SQLSelect
	Having(string) SQLSelect
	Order(string) SQLSelect
	Get(dbxGetSelect, any) error
	Select(dbxGetSelect, any) error
	Clear()
	String() string
}
type sqlselect struct {
	sel    string
	from   string
	join   string
	on     string
	where  string
	group  string
	having string
	order  string
	limit  int
	offset int
	args   []any
}

var _ SQLSelect = &sqlselect{}

var sqlstatpl = sync.Pool{New: func() any { return &sqlselect{} }}

func Select(es EgoSQL) SQLSelect {
	ret := sqlstatpl.Get().(*sqlselect)
	ret.Clear()

	ret.sel = strings.Join(es.SQLNames(), " , ")
	ret.args = append(ret.args, es.SQLValues()...)

	return ret
}

func (ss *sqlselect) From(tablename string) SQLSelect {
	ss.from = tablename
	return ss
}

func (ss *sqlselect) Join(tablename string) SQLSelect {
	ss.join = tablename
	return ss
}

func (ss *sqlselect) On(tablename string) SQLSelect {
	ss.on = tablename
	return ss
}

func (ss *sqlselect) Where(args map[string]any) SQLSelect {
	if len(args) > 0 {
		field := []string{}
		for k, v := range args {
			field = append(field, k+" = ? ")
			ss.args = append(ss.args, v)
		}
		ss.where = " WHERE " + strings.Join(field, " AND ")
	}
	return ss
}

func (ss *sqlselect) Limit(l int) SQLSelect {
	ss.limit = l
	return ss
}

func (ss *sqlselect) Offset(o int) SQLSelect {
	ss.offset = o
	return ss
}

func (ss *sqlselect) Group(group string) SQLSelect {
	ss.group = group
	return ss
}

func (ss *sqlselect) Having(having string) SQLSelect {
	ss.having = having
	return ss
}

func (ss *sqlselect) Order(order string) SQLSelect {
	ss.order = order
	return ss
}

func (ss *sqlselect) Get(dbx dbxGetSelect, data any) error {
	defer sqlstatpl.Put(ss)
	return dbx.Get(data, ss.String(), ss.args...)
}

func (ss *sqlselect) Select(dbx dbxGetSelect, data any) error {
	defer sqlstatpl.Put(ss)
	return dbx.Select(data, ss.String(), ss.args...)
}

func (ss *sqlselect) String() string {
	bs := &strings.Builder{}
	fmt.Fprintf(bs, "SELECT %s FROM %s ", ss.sel, ss.from)

	if ss.join != "" {
		bs.WriteByte(' ')
		bs.WriteString(ss.join)
		bs.WriteByte(' ')
		if ss.on != "" {
			fmt.Fprintf(bs, " on %s ", ss.on)
		}
	}

	if ss.where != "" {
		bs.WriteByte(' ')
		bs.WriteString(ss.where)
		bs.WriteByte(' ')
	}

	if ss.group != "" {
		fmt.Fprintf(bs, " GROUP BY %s ", ss.group)
		if ss.having != "" {
			fmt.Fprintf(bs, " HAVING %s ", ss.having)
		}
	}

	if ss.order != "" {
		fmt.Fprintf(bs, " ORDER BY %s ", ss.group)
	}

	if ss.limit != 0 {
		fmt.Fprintf(bs, " LIMIT %d ", ss.limit)
	}

	if ss.offset != 0 {
		fmt.Fprintf(bs, " OFFSET %d ", ss.offset)
	}
	return bs.String()
}

func (ss *sqlselect) Clear() {
	ss.sel = ""
	ss.from = ""
	ss.join = ""
	ss.on = ""
	ss.where = ""
	ss.limit = 0
	ss.offset = 0
	ss.args = ss.args[:0]
}
