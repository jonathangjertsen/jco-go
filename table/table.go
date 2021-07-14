package table

import (
	"fmt"
	"github.com/jonathangjertsen/jco-go/ops"
)

const (
	N_COLUMNS = 5
	PADDING   = 3
)

type Table struct {
	table [][N_COLUMNS]string
	bytes uint
}

func NewTable(bits uint) *Table {
	return &Table{
		table: [][N_COLUMNS]string{
			{"FORMULA", "|", "DECIMAL", "HEXADECIMAL", "BINARY"},
		},
		bytes: (bits + 7) / 8,
	}
}

func (t *Table) Add(name string, value []byte) {
	if t.bytes > uint(len(value)) {
		padding := t.bytes - uint(len(value))
		value = ops.PrependZeros(value, uint(padding))
	}
	valueTruncated := ops.Truncate(value, t.bytes)

	dec := ops.BytesToDec(valueTruncated, t.bytes)
	hex := ops.BytesToHex(valueTruncated, t.bytes)
	bin := ops.BytesToBin(valueTruncated, t.bytes)

	if !ops.Equivalent(value, valueTruncated) {
		dec = "*" + dec
		hex = "*" + hex
		bin = "*" + bin
	}

	t.table = append(t.table, [N_COLUMNS]string{
		name,
		"|",
		dec,
		hex,
		bin,
	})
}

func (t *Table) Render() {
	nRows := len(t.table)
	widths := [N_COLUMNS]int{}
	for c := 0; c < N_COLUMNS; c++ {
		for r := 0; r < nRows; r++ {
			widths[c] = ops.Intmax(widths[c], len(t.table[r][c]))
		}
	}
	formatStrings := [N_COLUMNS]string{}
	for c := 0; c < N_COLUMNS; c++ {
		formatStrings[c] = fmt.Sprintf("%%%ds", widths[c]+PADDING)
	}
	for r := 0; r < nRows; r++ {
		for c := 0; c < N_COLUMNS; c++ {
			fmt.Printf(formatStrings[c], t.table[r][c])
		}
		fmt.Println("")
	}
}
