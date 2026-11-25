package db0403

import (
	"errors"
	"slices"
)

type Schema struct {
	Table string
	Cols  []Column
	PKey  []int // indexes of primary key columns
}

type Column struct {
	Name string
	Type CellType
}

type Row []Cell

func (schema *Schema) NewRow() Row {
	return make(Row, len(schema.Cols))
}

func (row Row) EncodeKey(schema *Schema) (key []byte) {
	check(len(row) == len(schema.Cols))
	key = encodeStrKey(key, []byte(schema.Table))
	for _, idx := range schema.PKey {
		check(row[idx].Type == schema.Cols[idx].Type)
		key = row[idx].EncodeKey(key)
	}
	return key
}

func (row Row) EncodeVal(schema *Schema) (val []byte) {
	check(len(row) == len(schema.Cols))
	for idx, value := range row {
		if !slices.Contains(schema.PKey, idx) {
			check(value.Type == schema.Cols[idx].Type)
			val = row[idx].EncodeVal(val)
		}
	}
	return val
}

func (row Row) DecodeKey(schema *Schema, key []byte) (err error) {
	check(len(row) == len(schema.Cols))

	table, rest, err := decodeStrKey(key)
	if err != nil {
		return err
	}
	if string(table) != schema.Table {
		return errors.New("table name mismatch")
	}
	key = rest

	for _, idx := range schema.PKey {
		row[idx] = Cell{Type: schema.Cols[idx].Type}
		if key, err = row[idx].DecodeKey(key); err != nil {
			return err
		}
	}

	if len(key) != 0 {
		return errors.New("trailing garbage")
	}
	return nil
}

func (row Row) DecodeVal(schema *Schema, val []byte) (err error) {
	check(len(row) == len(schema.Cols))

	for idx, col := range schema.Cols {
		if slices.Contains(schema.PKey, idx) {
			continue
		}
		row[idx] = Cell{Type: col.Type}
		if val, err = row[idx].DecodeVal(val); err != nil {
			return err
		}
	}

	if len(val) != 0 {
		return errors.New("trailing garbage")
	}
	return nil
}

// QzBQWVJJOUhU https://trialofcode.org/
