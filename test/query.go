package test

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/IceySam/serve-soft/utility"
	"github.com/jackc/pgx/v5"
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
	 * 	execute select or queries returning many rows
	 */
	Many() ([]map[string]interface{}, error)
	/*
	 * 	execute select scan into provided interface
	 */
	One(i interface{}) error
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
	strutType string
}

type Query struct {
	Conn *pgx.Conn
}

// create table
func (q Query) Create(name string, definition ...string) error {
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(\n", name)
	for i := 0; i < len(definition); i++ {
		stmt = fmt.Sprintf("%s %s", stmt, definition[i])
		if i+1 < len((definition)) {
			stmt = fmt.Sprintf("%s,\n", stmt)
		} else {
			stmt = fmt.Sprintf("%s\n);", stmt)
		}
	}

	_, err := q.Conn.Exec(context.Background(), stmt)
	if err != nil {
		return err
	}
	return nil
}

// insert into table
func (q Query) Insert(i interface{}) error {
	m, ty, err := utility.ToMap(i)
	if err != nil {
		return err
	}
	stmt := fmt.Sprintf("INSERT INTO %s (", ty)
	x := 1
	for k := range m {
		if x == len(m) {
			stmt = fmt.Sprintf("%s%s)", stmt, k)
		} else {
			stmt = fmt.Sprintf("%s%s,", stmt, k)
		}
		x++
	}
	stmt = fmt.Sprintf("%s VALUES (", stmt)
	y := 1
	for _, v := range m {
		if reflect.TypeOf(v).ConvertibleTo(reflect.TypeOf("")) {
			v = fmt.Sprintf("'%v'", v)
		}
		if y == len(m) {
			stmt = fmt.Sprintf("%s%v);", stmt, v)
		} else {
			stmt = fmt.Sprintf("%s%v,", stmt, v)
		}
		y++
	}

	_, err = q.Conn.Exec(context.Background(), stmt)
	if err != nil {
		return err
	}
	return nil
}

// set update data, implements Completer.
func (p *partialQuery) Set(m map[string]interface{}) Completer {
	stmt := ""
	index := 1
	for k, v := range m {
		if index == len(m) {
			stmt = fmt.Sprintf("%s%s='%v'", stmt, k, v)
		} else {
			stmt = fmt.Sprintf("%s%s='%v',", stmt, k, v)
		}
		index++
	}
	stmt = fmt.Sprintf("%s%s", p.part, stmt)
	p.part = stmt
	return p
}

// Apply implements Completer.
func (p *partialQuery) Apply() error {
	_, err := p.query.Conn.Exec(context.Background(), p.part+";")
	if err != nil {
		return err
	}
	return nil
}

// Many implements Completer.
func (p *partialQuery) Many() ([]map[string]interface{}, error) {
	stmt := fmt.Sprintf("%s;", p.part)
	return p.fetchData(stmt)
}

// One implements Completer.
func (p *partialQuery) One(i interface{}) error {
	stmt := fmt.Sprintf("%s LIMIT 1;", p.part)

	items, err := p.fetchData(stmt)
	if err != nil {
		return err
	}

	utility.ToStuct(items[0], i)

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
			stmt = fmt.Sprintf("%s%s='%v'", stmt, k, v)
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
	m, ty, err := utility.ToMap(i)
	if err != nil {
		log.Fatalln(err)
	}
	stmt := fmt.Sprintf("UPDATE %s SET ", ty)

	return &partialQuery{part: stmt, query: q, data: m, strutType: ty}
}

// delete data
func (q Query) Delete(i interface{}) Completer {
	m, ty, err := utility.ToMap(i)
	if err != nil {
		log.Fatalln(err)
	}
	stmt := fmt.Sprintf("DELETE FROM %s", ty)
	return &partialQuery{part: stmt, query: q, data: m, strutType: ty}
}

// find all from relation
func (q Query) FindAll(i interface{}) ([]map[string]interface{}, error) {
	m, ty, err := utility.ToMap(i)
	if err != nil {
		return nil, err
	}
	p := &partialQuery{part: "", query: q, data: m, strutType: ty}
	return p.fetchData(fmt.Sprintf("SELECT * FROM %s;", ty))
}

// find all from relation
func (p *partialQuery) fetchData(stmt string) ([]map[string]interface{}, error) {
	rows, err := p.query.Conn.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fieldDes := rows.FieldDescriptions()
	items := make([]map[string]interface{}, 0)

	for rows.Next() {
		columns := make([]interface{}, len(fieldDes))
		item := make([]interface{}, len(fieldDes))
		for x := range columns {
			item[x] = &columns[x]
		}

		if err := rows.Scan(item...); err != nil {
			return nil, err
		}
		res := make(map[string]interface{}, len(p.data))
		for index, val := range fieldDes {
			r := item[index].(*interface{})
			res[val.Name] = *r
		}

		items = append(items, res)
	}
	return items, nil
}

// find data from relation
func (q Query) Find(i interface{}) Completer {
	m, ty, err := utility.ToMap(i)
	if err != nil {
		log.Fatalln(err)
	}
	stmt := fmt.Sprintf("SELECT * FROM %s", ty)
	return &partialQuery{part: stmt, query: q, data: m, strutType: ty}
}
