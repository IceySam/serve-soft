package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/IceySam/serve-soft/utility"
)

type Completer interface {
	/*
	 *	set update data
	 * 	m equals struct matching db relation
	 */
	Set(m map[string]interface{}) Completer
	/*
	 * 	execute non returning query like update and delete
	 */
	Apply() error
	/*
	 * 	execute non returning query like update and delete with context
	 */
	ApplyCtx(ctx context.Context) error
	/*
	 * 	execute non returning query like update and delete with context in a transaction
	 */
	TxApplyCtx(ctx context.Context, tx *sql.Tx) error
	/*
	 * 	execute select or queries returning many rows
	 */
	Many(i interface{}) error
	/*
	 * 	execute select or queries returning many rows with context
	 */
	ManyCtx(ctx context.Context, i interface{}) error
	/*
	 * 	execute select or queries returning many rows with context in a transaction
	 */
	TxManyCtx(ctx context.Context, tx *sql.Tx, i interface{}) error
	/*
	 * 	execute select scan into provided interface
	 */
	One(i interface{}) error
	/*
	 * 	execute select scan into provided interface with context
	 */
	OneCtx(ctx context.Context, i interface{}) error
	/*
	 * 	execute select scan into provided interface with context in a transaction
	 */
	TxOneCtx(ctx context.Context, tx *sql.Tx, i interface{}) error
	/*
	 *	where condition
	 *	could be map[string]interface{} or []map[string]interface{}
	 *  e.g. map[string]interface{ "id": 1 }
	 */
	Where(any) Completer
	/*
	 * 	perform query for values in list
	 */
	In(field string, values []interface{}) Completer
}

type partialQuery struct {
	part      string
	query     Query
	data      map[string]interface{}
	strutType reflect.Type
}

type Query struct {
	Conn *sql.DB
}

// create table
func (q Query) Create(name string, definition ...string) error {
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(\n", name)
	for i := 0; i < len(definition); i++ {
		stmt = fmt.Sprintf("%s %s", stmt, definition[i])
		if i+1 < len((definition)) {
			stmt = fmt.Sprintf("%s,\n", stmt)
		} else {
			stmt = fmt.Sprintf("%s\n);", stmt)
		}
	}

	_, err := q.Conn.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

// create table with context
func (q Query) CreateCtx(ctx context.Context, name string, definition ...string) error {
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(\n", name)
	for i := 0; i < len(definition); i++ {
		stmt = fmt.Sprintf("%s %s", stmt, definition[i])
		if i+1 < len((definition)) {
			stmt = fmt.Sprintf("%s,\n", stmt)
		} else {
			stmt = fmt.Sprintf("%s\n);", stmt)
		}
	}

	_, err := q.Conn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

// tx create table with context
func (q Query) TxCreateCtx(ctx context.Context, tx *sql.Tx, name string, definition ...string) error {
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(\n", name)
	for i := 0; i < len(definition); i++ {
		stmt = fmt.Sprintf("%s %s", stmt, definition[i])
		if i+1 < len((definition)) {
			stmt = fmt.Sprintf("%s,\n", stmt)
		} else {
			stmt = fmt.Sprintf("%s\n);", stmt)
		}
	}

	_, err := tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

// insert into table
func (q Query) Insert(i interface{}) (int64, error) {
	stmt, values, err := prepareInsert(i)
	if err != nil {
		return 0, err
	}
	res, err := q.Conn.Exec(stmt, values...)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

// insert into table with context
func (q Query) InsertCtx(ctx context.Context, i interface{}) (int64, error) {
	stmt, values, err := prepareInsert(i)
	if err != nil {
		return 0, err
	}
	res, err := q.Conn.ExecContext(ctx, stmt, values...)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

// TX insert into table with context
func (q Query) TxInsertCtx(ctx context.Context, tx *sql.Tx, i interface{}) (int64, error) {
	stmt, values, err := prepareInsert(i)
	if err != nil {
		return 0, err
	}
	res, err := tx.ExecContext(ctx, stmt, values...)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func prepareInsert(i interface{}) (string, []interface{}, error) {
	m, _, name, err := utility.ToMap(i)
	if err != nil {
		return "", []interface{}{}, err
	}

	keys := make([]string, 0, len(m))
	placeholders := make([]string, 0, len(m))
	values := make([]interface{}, 0, len(m))

	for k, v := range m {
		if v != nil {
			keys = append(keys, fmt.Sprintf("`%s`",k))	
			placeholders = append(placeholders, "?")
			values = append(values, v)
		}
	}

	stmt := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", name, strings.Join(keys, ", "), strings.Join(placeholders, ", "))
	return stmt, values, nil
}

// set update data, implements Completer.
func (p *partialQuery) Set(m map[string]interface{}) Completer {
	stmt := ""
	index := 1
	for k, v := range m {
		k = fmt.Sprintf("`%s`", k)
		if reflect.TypeOf(v).ConvertibleTo(reflect.TypeOf("")) {
			v = fmt.Sprintf("'%v'", v)
		}
		if index == len(m) {
			stmt = fmt.Sprintf("%s%s=%v", stmt, k, v)
		} else {
			stmt = fmt.Sprintf("%s%s=%v,", stmt, k, v)
		}
		index++
	}
	stmt = fmt.Sprintf("%s%s", p.part, stmt)
	p.part = stmt
	return p
}

// Apply implements Completer.
func (p *partialQuery) Apply() error {
	_, err := p.query.Conn.Exec(p.part + ";")
	if err != nil {
		return err
	}
	return nil
}

// Apply with context implements Completer.
func (p *partialQuery) ApplyCtx(ctx context.Context) error {
	_, err := p.query.Conn.ExecContext(ctx, p.part+";")
	if err != nil {
		return err
	}
	return nil
}

// TX Apply with context implements Completer.
func (p *partialQuery) TxApplyCtx(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, p.part+";")
	if err != nil {
		return err
	}
	return nil
}

// Many implements Completer.
func (p *partialQuery) Many(i interface{}) error {
	stmt := fmt.Sprintf("%s;", p.part)

	items, err := p.fetchData(stmt)
	if err != nil {
		return err
	}

	err = utility.ToStructArray(items, i)
	if err != nil {
		return err
	}
	return nil
}

// Many with context implements Completer.
func (p *partialQuery) ManyCtx(ctx context.Context, i interface{}) error {
	stmt := fmt.Sprintf("%s;", p.part)

	items, err := p.fetchData(stmt, ctx)
	if err != nil {
		return err
	}

	err = utility.ToStructArray(items, i)
	if err != nil {
		return err
	}
	return nil
}

// TX Many with context implements Completer.
func (p *partialQuery) TxManyCtx(ctx context.Context, tx *sql.Tx, i interface{}) error {
	stmt := fmt.Sprintf("%s;", p.part)

	items, err := p.fetchData(stmt, ctx, tx)
	if err != nil {
		return err
	}

	err = utility.ToStructArray(items, i)
	if err != nil {
		return err
	}
	return nil
}

// One implements Completer.
func (p *partialQuery) One(i interface{}) error {
	stmt := fmt.Sprintf("%s LIMIT 1;", p.part)

	items, err := p.fetchData(stmt)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return fmt.Errorf("query did not return any result")
	}
	err = utility.ToStruct(items[0], i)
	if err != nil {
		return err
	}
	return nil
}

// One with context implements Completer.
func (p *partialQuery) OneCtx(ctx context.Context, i interface{}) error {
	stmt := fmt.Sprintf("%s LIMIT 1;", p.part)

	items, err := p.fetchData(stmt, ctx)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return fmt.Errorf("query did not return any result")
	}
	err = utility.ToStruct(items[0], i)
	if err != nil {
		return err
	}
	return nil
}

// TX One with context implements Completer.
func (p *partialQuery) TxOneCtx(ctx context.Context, tx *sql.Tx, i interface{}) error {
	stmt := fmt.Sprintf("%s LIMIT 1;", p.part)

	items, err := p.fetchData(stmt, ctx, tx)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return fmt.Errorf("query did not return any result")
	}
	err = utility.ToStruct(items[0], i)
	if err != nil {
		return err
	}
	return nil
}

// Where implements Completer.
func (p *partialQuery) Where(m any) Completer {
	stmt := " WHERE "

	var list []map[string]interface{}

	hasOr := utility.TypeEquals(m, list)
	if hasOr {
		list = m.([]map[string]interface{})
	} else {
		list = append(list, m.(map[string]interface{}))
	}

	index := 1
	for i := range list {
		if index > 1 {
			stmt = fmt.Sprintf("%s OR (", stmt)
		}

		x := 1
		for k, v := range list[i] {
			stmt = fmt.Sprintf("%s`%s`='%v'", stmt, k, v)
			if x < len(list[i]) {
				stmt = fmt.Sprintf("%s AND ", stmt)
			}
			x++
		}

		if index > 1 {
			stmt = fmt.Sprintf("%s)", stmt)
		}

		index++
	}

	stmt = fmt.Sprintf("%s%s", p.part, stmt)
	p.part = stmt
	return p
}

// find where in list implements Completer.
func (p *partialQuery) In(field string, values []interface{}) Completer {
	stmt := fmt.Sprintf(" WHERE %s IN(", field)
	for i := 0; i < len(values); i++ {
		if i+1 == len(values) {
			stmt = fmt.Sprintf("%s'%v')", stmt, values[i])
		} else {
			stmt = fmt.Sprintf("%s'%v', ", stmt, values[i])
		}
	}
	stmt = fmt.Sprintf("%s%s", p.part, stmt)
	p.part = stmt
	return p
}

// update table
func (q Query) Update(i interface{}) Completer {
	m, ty, name, err := utility.ToMap(i)
	if err != nil {
		log.Fatalln(err)
	}
	stmt := fmt.Sprintf("UPDATE %s SET ", name)

	return &partialQuery{part: stmt, query: q, data: m, strutType: ty}
}

// delete data
func (q Query) Delete(i interface{}) Completer {
	m, ty, name, err := utility.ToMap(i)
	if err != nil {
		log.Fatalln(err)
	}
	stmt := fmt.Sprintf("DELETE FROM %s", name)
	return &partialQuery{part: stmt, query: q, data: m, strutType: ty}
}

// find all from relation
func (q Query) FindAll(i interface{}) ([]map[string]interface{}, error) {
	m, ty, name, err := utility.ToMap(i)
	if err != nil {
		return nil, err
	}
	p := &partialQuery{part: "", query: q, data: m, strutType: ty}
	return p.fetchData(fmt.Sprintf("SELECT * FROM %s;", name))
}

// find all with context from relation
func (q Query) FindAllCtx(ctx context.Context, i interface{}) ([]map[string]interface{}, error) {
	m, ty, name, err := utility.ToMap(i)
	if err != nil {
		return nil, err
	}
	p := &partialQuery{part: "", query: q, data: m, strutType: ty}
	return p.fetchData(fmt.Sprintf("SELECT * FROM %s;", name), ctx)
}

// TX find all with context from relation
func (q Query) TxFindAllCtx(ctx context.Context, tx *sql.Tx, i interface{}) ([]map[string]interface{}, error) {
	m, ty, name, err := utility.ToMap(i)
	if err != nil {
		return nil, err
	}
	p := &partialQuery{part: "", query: q, data: m, strutType: ty}
	return p.fetchData(fmt.Sprintf("SELECT * FROM %s;", name), ctx, tx)
}

// find all from relation
func (p *partialQuery) fetchData(stmt string, args ...interface{}) ([]map[string]interface{}, error) {
	var rows *sql.Rows
	var err error

	if len(args) == 0 {
		rows, err = p.query.Conn.Query(stmt)
	} else if len(args) == 1 {
		rows, err = p.query.Conn.QueryContext(args[0].(context.Context), stmt)
	} else {
		rows, err = args[1].(*sql.Tx).QueryContext(args[0].(context.Context), stmt)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fieldDes, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0)

	for rows.Next() {
		columns := make([]sql.RawBytes, len(fieldDes))
		item := make([]interface{}, len(fieldDes))
		for x := range columns {
			item[x] = &columns[x]
		}

		if err := rows.Scan(item...); err != nil {
			return nil, err
		}

		res := make(map[string]interface{}, len(p.data))
		for i, col := range columns {
			str := string(col)
			if col == nil {
				res[fieldDes[i]] = nil
			} else {
				res[fieldDes[i]] = str
			}
		}
		items = append(items, res)
	}
	return items, nil
}

// find data from relation
func (q Query) Find(i interface{}) Completer {
	m, ty, name, err := utility.ToMap(i)
	if err != nil {
		log.Fatalln(err)
	}
	stmt := fmt.Sprintf("SELECT * FROM %s", name)
	return &partialQuery{part: stmt, query: q, data: m, strutType: ty}
}
