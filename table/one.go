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
	t.Add(fmt.Sprintf("reverse_bitstring(%s)", metavar), ops.BitstringReverse(a))
	t.Add(fmt.Sprintf("reverse_bitorder(%s)", metavar), ops.BitReverse(a))
	t.Add(fmt.Sprintf("reverse_byteorder(%s)", metavar), ops.ByteReverse(a))
	t.Add(fmt.Sprintf("reverse_nibbleorder(%s)", metavar), ops.NibbleSwap(a))
}
