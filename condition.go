package pgs

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type Identifiable interface {
	exp.Inable
	exp.Comparable
	exp.Likeable
	exp.Isable
}

type Conditional interface {
	Condition(inUpdate bool) (goqu.Expression, error)
	getJoiners() []*joiner
	//getOp() string
	//getValue() interface{}
}

type Condition struct {
	Field   fieldI
	Op      string
	Value   interface{}
	joiners []*joiner
}

func (c Condition) Condition(inUpdate bool) (exp.Expression, error) {
	var ident exp.IdentifierExpression
	ident = c.Field.getIdent()
	if inUpdate {
		model := c.Field.getModel()
		if model.joiner != nil {
			ident = goqu.I(fmt.Sprintf("%s.%s", model.joiner.ParentTable, model.joiner.From))
		}
	}
	var condition exp.Expression
	switch c.Op {
	case opIn:
		condition = ident.In(c.Value)
	case opNotIn:
		condition = ident.NotIn(c.Value)
	case opEq:
		condition = ident.Eq(c.Value)
	case opNotEq:
		condition = ident.Neq(c.Value)
	case opLike:
		condition = ident.Like(c.Value)
	case opNotLike:
		condition = ident.NotLike(c.Value)
	case opRegex:
		condition = ident.RegexpLike(c.Value)
	case opRegexI:
		condition = ident.RegexpILike(c.Value)
	case opNotRegex:
		condition = ident.RegexpNotLike(c.Value)
	case opNotRegexI:
		condition = ident.RegexpNotILike(c.Value)
	case opLt:
		condition = ident.Lt(c.Value)
	case opLte:
		condition = ident.Lte(c.Value)
	case opGt:
		condition = ident.Gt(c.Value)
	case opGte:
		condition = ident.Gte(c.Value)
	case opIsNotNull:
		condition = ident.IsNotNull()
	case opIsNull:
		condition = ident.IsNull()
	default:
		return condition, fmt.Errorf("operator %s can not be found", c.Op)
	}
	return condition, nil
}

func (c Condition) getJoiners() []*joiner {
	return c.joiners
}

type OrCondition struct {
	Conditions []Condition
}

func Or(conditions ...Condition) OrCondition {
	return OrCondition{conditions}
}

func (oe OrCondition) Condition(inUpdate bool) (exp.Expression, error) {
	var exps []exp.Expression
	for _, cond := range oe.Conditions {
		expr, err := cond.Condition(inUpdate)
		if err != nil {
			return nil, err
		}
		exps = append(exps, expr)
	}
	return goqu.Or(exps...), nil
}

func (oe OrCondition) getJoiners() []*joiner {
	var joiners []*joiner
	for _, cond := range oe.Conditions {
		joiners = append(joiners, cond.getJoiners()...)
	}
	return joiners
}
