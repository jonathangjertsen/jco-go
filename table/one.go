package table

import (
	"fmt"
	"github.com/jonathangjertsen/jco-go/ops"
)

func (t *Table) One(a []byte, metavar string) {
	t.Add(fmt.Sprintf("%s ", metavar), a)
	t.Add(fmt.Sprintf("~%s ", metavar), ops.Not(a))
	t.Add(fmt.Sprintf("twos_complement(%s)", metavar), ops.TwosComplement(a))
	t.Add(fmt.Sprintf("popcount(%s)", metavar), ops.Popcount(a))
	t.Add(fmt.Sprintf("clz(%s)", metavar), ops.Clz(a))
	t.Add(fmt.Sprintf("nbits(%s)", metavar), ops.Nbits(a))
}
