package db0305

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

type DB struct {
	KV     KV
	tables map[string]Schema
}

func (db *DB) Open() error {
	db.tables = map[string]Schema{}
	return db.KV.Open()
}

func (db *DB) Close() error { return db.KV.Close() }

func (db *DB) Select(schema *Schema, row Row) (ok bool, err error) {
	key := row.EncodeKey(schema)
	val, ok, err := db.KV.Get(key)
	if err != nil || !ok {
		return ok, err
	}
	if err = row.DecodeVal(schema, val); err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) Insert(schema *Schema, row Row) (updated bool, err error) {
	key := row.EncodeKey(schema)
	val := row.EncodeVal(schema)
	return db.KV.SetEx(key, val, ModeInsert)
}

func (db *DB) Upsert(schema *Schema, row Row) (updated bool, err error) {
	key := row.EncodeKey(schema)
	val := row.EncodeVal(schema)
	return db.KV.SetEx(key, val, ModeUpsert)
}

func (db *DB) Update(schema *Schema, row Row) (updated bool, err error) {
	key := row.EncodeKey(schema)
	val := row.EncodeVal(schema)
	return db.KV.SetEx(key, val, ModeUpdate)
}

func (db *DB) Delete(schema *Schema, row Row) (deleted bool, err error) {
	key := row.EncodeKey(schema)
	return db.KV.Del(key)
}

type SQLResult struct {
	Updated int
	Header  []string //named attributes - columns
	Values  []Row
}

func (db *DB) ExecStmt(stmt interface{}) (r SQLResult, err error) {
	switch ptr := stmt.(type) { //type switch
	case *StmtCreatTable:
		err = db.execCreateTable(ptr)
	case *StmtSelect:
		r.Header = ptr.cols
		r.Values, err = db.execSelect(ptr)
	case *StmtInsert:
		r.Updated, err = db.execInsert(ptr)
	case *StmtUpdate:
		r.Updated, err = db.execUpdate(ptr)
	case *StmtDelete:
		r.Updated, err = db.execDelete(ptr)
	default:
		panic("unreachable")
	}
	return
}

func (db *DB) execCreateTable(stmt *StmtCreatTable) (err error) {
	pkeys := []int{}

	for _, key := range stmt.pkey {
		idx := slices.IndexFunc(stmt.cols, func(c Column) bool { return c.Name == key })
		if idx < 0 {
			return fmt.Errorf("row should have at least one primary key")
		}
		pkeys = append(pkeys, idx)
	}
	schema := Schema{Table: stmt.table, Cols: stmt.cols, PKey: pkeys}
	val, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	if _, err = db.KV.SetEx([]byte("@schema_"+stmt.table), val, ModeUpsert); err != nil {
		return err
	}
	db.tables[stmt.table] = schema
	return nil
}

func (db *DB) GetSchema(table string) (Schema, error) {
	schema, ok := db.tables[table]
	if !ok {
		val, ok, err := db.KV.Get([]byte("@schema_" + table))
		if err == nil && ok {
			err = json.Unmarshal(val, &schema)
		}
		if err != nil {
			return Schema{}, err
		}
		if !ok {
			return Schema{}, errors.New("table is not found")
		}
		db.tables[table] = schema
	}
	return schema, nil
}

func (db *DB) execSelect(stmt *StmtSelect) ([]Row, error) {
	schema, err := db.GetSchema(stmt.table)
	if err != nil {
		return nil, fmt.Errorf("expected correct table name. Table %s not found", stmt.table)
	}
	pKeys := []string{}
	for _, val := range schema.PKey {
		pKeys = append(pKeys, schema.Cols[val].Name)
	}
	for _, val := range stmt.keys {
		if !slices.Contains(pKeys, val.column) {
			return nil, fmt.Errorf("%s not a primary key", val.column)
		}
	}
	row := schema.NewRow()
	for _, key := range stmt.keys {
		idx := slices.IndexFunc(schema.Cols, func(c Column) bool { return c.Name == key.column })
		if idx < 0 {
			return nil, fmt.Errorf("column %s not found in table %s", key.column, stmt.table)
		}
		if key.value.Type != schema.Cols[idx].Type {
			return nil, fmt.Errorf("type mismatch for column %s", key.column)
		}
		row[idx] = key.value
	}
	for _, idx := range schema.PKey {
		if row[idx].Type == 0 {
			return nil, fmt.Errorf("missing value for primary key column %s", schema.Cols[idx].Name)
		}
	}
	ok, err := db.Select(&schema, row)
	if err != nil {
		return nil, err
	}
	if !ok {
		return []Row{}, nil
	}
	result := make(Row, 0, len(stmt.cols))
	for _, name := range stmt.cols {
		idx := slices.IndexFunc(schema.Cols, func(c Column) bool { return c.Name == name })
		if idx < 0 {
			return nil, fmt.Errorf("column %s not found in table %s", name, stmt.table)
		}
		result = append(result, row[idx])
	}
	return []Row{result}, nil
}

func (db *DB) execInsert(stmt *StmtInsert) (count int, err error) {
	if schema, err := db.GetSchema(stmt.table); err != nil {
		return 0, fmt.Errorf("expected correct table name. Table %s not found", stmt.table)
	} else {
		row := Row(stmt.value)
		_, err := db.Insert(&schema, row)
		if err != nil {
			return 0, nil
		}
		return 1, nil // count is amount of rows affected?
	}
}

func (db *DB) execUpdate(stmt *StmtUpdate) (count int, err error) {
	schema, err := db.GetSchema(stmt.table)
	if err != nil {
		return 0, fmt.Errorf("expected correct table name. Table %s not found", stmt.table)
	}
	pKeys := []string{}
	for _, val := range schema.PKey {
		pKeys = append(pKeys, schema.Cols[val].Name)
	}
	for _, val := range stmt.keys {
		if !slices.Contains(pKeys, val.column) {
			return 0, fmt.Errorf("%s not a primary key", val.column)
		}
	}
	row := schema.NewRow()
	for _, key := range stmt.keys {
		idx := slices.IndexFunc(schema.Cols, func(c Column) bool { return c.Name == key.column })
		if idx < 0 {
			return 0, fmt.Errorf("column %s not found in table %s", key.column, stmt.table)
		}
		if key.value.Type != schema.Cols[idx].Type {
			return 0, fmt.Errorf("type mismatch for column %s", key.column)
		}
		row[idx] = key.value
	}
	for _, idx := range schema.PKey {
		if row[idx].Type == 0 {
			return 0, fmt.Errorf("missing value for primary key column %s", schema.Cols[idx].Name)
		}
	}
	ok, err := db.Select(&schema, row)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	for _, assign := range stmt.value {
		idx := slices.IndexFunc(schema.Cols, func(c Column) bool { return c.Name == assign.column })
		if idx < 0 {
			return 0, fmt.Errorf("column %s not found in table %s", assign.column, stmt.table)
		}
		if slices.Contains(schema.PKey, idx) {
			return 0, fmt.Errorf("cannot update primary key column %s", assign.column)
		}
		if assign.value.Type != schema.Cols[idx].Type {
			return 0, fmt.Errorf("type mismatch for column %s", assign.column)
		}
		row[idx] = assign.value
	}
	ok, err = db.Update(&schema, row)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	return 1, nil
}

func (db *DB) execDelete(stmt *StmtDelete) (count int, err error) {
	schema, err := db.GetSchema(stmt.table)
	if err != nil {
		return 0, fmt.Errorf("expected correct table name. Table %s not found", stmt.table)
	}
	pKeys := []string{}
	for _, val := range schema.PKey {
		pKeys = append(pKeys, schema.Cols[val].Name)
	}
	for _, val := range stmt.keys {
		if !slices.Contains(pKeys, val.column) {
			return 0, fmt.Errorf("%s not a primary key", val.column)
		}
	}
	row := schema.NewRow()
	for _, key := range stmt.keys {
		idx := slices.IndexFunc(schema.Cols, func(c Column) bool { return c.Name == key.column })
		if idx < 0 {
			return 0, fmt.Errorf("column %s not found in table %s", key.column, stmt.table)
		}
		if key.value.Type != schema.Cols[idx].Type {
			return 0, fmt.Errorf("type mismatch for column %s", key.column)
		}
		row[idx] = key.value
	}
	for _, idx := range schema.PKey {
		if row[idx].Type == 0 {
			return 0, fmt.Errorf("missing value for primary key column %s", schema.Cols[idx].Name)
		}
	}
	ok, err := db.Delete(&schema, row)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	return 1, nil
}

// QzBQWVJJOUhU https://trialofcode.org/
