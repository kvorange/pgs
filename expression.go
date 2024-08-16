package pgs

import "github.com/doug-martin/goqu/v9"

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

func (c CountExpression) getJoiner() *joiner {
	joiner := c.field.getJoiner()
	return joiner
}

func (c CountExpression) As(as string) CountExpression {
	c.as = as
	return c
}
