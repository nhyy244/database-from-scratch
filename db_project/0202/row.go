package db0202

import (
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

	rows := make(Row, len(schema.Cols))
	for idx := range rows {
		rows[idx].Type = schema.Cols[idx].Type
	}
	return rows
}

func (row Row) EncodeKey(schema *Schema) (key []byte) {
	key = append(key, []byte(schema.Table)...)
	key = append(key, 0x00)

	for _, v := range schema.PKey {
		key = row[v].Encode(key)
	}
	return key
}

func (row Row) EncodeVal(schema *Schema) (val []byte) {
	for idx := range schema.Cols {
		if !slices.Contains(schema.PKey, idx) {
			val = row[idx].Encode(val)
		}
	}
	return val
}

func (row Row) DecodeKey(schema *Schema, key []byte) (err error) {
	key = key[len(schema.Table)+1:]

	for _, idx := range schema.PKey {
		key, err = row[idx].Decode(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (row Row) DecodeVal(schema *Schema, val []byte) (err error) {
	for idx := range schema.Cols {
		if !slices.Contains(schema.PKey, idx) {
			row[idx].Decode(val)
		}
	}
	return nil
}

// QzBQWVJJOUhU https://trialofcode.org/
