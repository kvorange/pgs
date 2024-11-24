package pgs

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type CountExpression struct {
	field fieldI
	as    string
}

func Count(field fieldI) CountExpression {
	return CountExpression{field, ""}
}

func (c CountExpression) getSelectors() []interface{} {
	ident := c.field.getIdent()
	if c.as != "" {
		return []interface{}{goqu.COUNT(ident).As(c.as)}
	}
	return []interface{}{goqu.COUNT(ident)}
}

func (c CountExpression) getJoiners() []*joiner {
	return c.field.getJoiners()
}

func (c CountExpression) As(as string) CountExpression {
	c.as = as
	return c
}

type LiteralExpression struct {
	as         string
	expression exp.LiteralExpression
	joiners    []*joiner
}

func L(sql string, values ...interface{}) LiteralExpression {
	var j []*joiner
	var args []interface{}
	for _, value := range values {
		field, ok := value.(fieldI)
		if ok {
			j = append(j, field.getJoiners()...)
			args = append(args, field.getIdent())
		} else {
			args = append(args, value)
		}

	}
	return LiteralExpression{
		expression: goqu.L(sql, args...),
		joiners:    j,
	}
}

func (l LiteralExpression) As(as string) LiteralExpression {
	l.as = as
	return l
}

func (l LiteralExpression) getSelectors() []interface{} {
	if l.as != "" {
		return []interface{}{l.expression.As(l.as)}
	}
	return []interface{}{l.expression}
}

func (l LiteralExpression) getJoiners() []*joiner {
	return l.joiners
}

func (l LiteralExpression) Condition(inUpdate bool) (goqu.Expression, error) {
	return l.expression.Expression(), nil
}
