package table

import (
	"fmt"
	"github.com/jonathangjertsen/jco-go/ops"
)

func (t *Table) Two(a []byte, b []byte, metavar1 string, metavar2 string) {
	t.Add(fmt.Sprintf("      %s", metavar1), a)
	t.Add(fmt.Sprintf("      %s", metavar2), b)

	t.Add(fmt.Sprintf("%s  + %s", metavar1, metavar2), ops.Add(a, b))
	t.Add(fmt.Sprintf("%s  | %s", metavar1, metavar2), ops.Or(a, b))
	t.Add(fmt.Sprintf("%s  & %s", metavar1, metavar2), ops.And(a, b))
	t.Add(fmt.Sprintf("%s  ^ %s", metavar1, metavar2), ops.Xor(a, b))
	t.Add(fmt.Sprintf("%s ^~ %s", metavar1, metavar2), ops.Xor(a, ops.Not(b)))

	t.Add(fmt.Sprintf("%s  - %s", metavar1, metavar2), ops.Subtract(a, b))
	t.Add(fmt.Sprintf("%s &~ %s", metavar1, metavar2), ops.And(a, ops.Not(b)))
	//t.Add(fmt.Sprintf("%s >> %s", metavar1, metavar2), a>>b)
	//t.Add(fmt.Sprintf("%s << %s", metavar1, metavar2), a<<b)
	t.Add(fmt.Sprintf("%s  - %s", metavar2, metavar1), ops.Subtract(b, a))
	t.Add(fmt.Sprintf("%s &~ %s", metavar2, metavar1), ops.And(b, ops.Not(a)))
	//t.Add(fmt.Sprintf("%s >> %s", metavar2, metavar1), b>>a)
	//t.Add(fmt.Sprintf("%s << %s", metavar2, metavar1), b<<a)
}
