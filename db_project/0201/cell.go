package db0201

import (
	"encoding/binary"
)

type CellType uint8

const (
	TypeI64 CellType = 1
	TypeStr CellType = 2
)

type Cell struct {
	Type CellType
	I64  int64
	Str  []byte
}

func (cell *Cell) Encode(toAppend []byte) []byte {
	switch cell.Type {
	case TypeI64:
		toAppend = binary.LittleEndian.AppendUint64(toAppend, uint64(cell.I64))

	case TypeStr:
		toAppend = binary.LittleEndian.AppendUint32(toAppend, uint32(len(cell.Str)))
		toAppend = append(toAppend, cell.Str...)
	}

	return toAppend
}

func (cell *Cell) Decode(data []byte) (rest []byte, err error) {
	switch cell.Type {
	case TypeI64:
		cell.I64 = int64(binary.LittleEndian.Uint64(data[:]))
		rest = data[8:] //8 for in64
	case TypeStr:
		strLen := binary.LittleEndian.Uint32(data[:4])
		cell.Str = data[4 : 4+strLen]
		rest = data[4+strLen:]
	}
	return rest, nil
}

// QzBQWVJJOUhU https://trialofcode.org/
