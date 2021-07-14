package table

import (
	"github.com/jonathangjertsen/jco-go/ops"
	"github.com/olekukonko/tablewriter"
	"os"
)

type Table struct {
	table *tablewriter.Table
	bytes uint
}

func NewTable(bits uint) *Table {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Formula", "|", "Decimal", "Hexadecimal", "Binary"})
	t.SetBorder(false)
	t.SetHeaderLine(false)
	t.SetTablePadding("\t")
	t.SetCenterSeparator("")
	t.SetColumnSeparator("")
	t.SetHeaderAlignment(tablewriter.ALIGN_RIGHT)
	t.SetAlignment(tablewriter.ALIGN_RIGHT)
	return &Table{
		table: t,
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

	t.table.Append([]string{
		name,
		"|",
		dec,
		hex,
		bin,
	})
}

func (t *Table) Render() {
	t.table.Render()
}
