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
	Condition() (goqu.Expression, error)
	getJoiners() []*joiner
}

type Condition struct {
	Field   Identifiable
	Op      string
	Value   interface{}
	joiners []*joiner
}

func (c Condition) Condition() (exp.Expression, error) {
	var condition exp.Expression
	switch c.Op {
	case opIn:
		condition = c.Field.In(c.Value)
	case opNotIn:
		condition = c.Field.NotIn(c.Value)
	case opEq:
		condition = c.Field.Eq(c.Value)
	case opNotEq:
		condition = c.Field.Neq(c.Value)
	case opLike:
		condition = c.Field.Like(c.Value)
	case opNotLike:
		condition = c.Field.NotLike(c.Value)
	case opRegex:
		condition = c.Field.RegexpLike(c.Value)
	case opRegexI:
		condition = c.Field.RegexpILike(c.Value)
	case opNotRegex:
		condition = c.Field.RegexpNotLike(c.Value)
	case opNotRegexI:
		condition = c.Field.RegexpNotILike(c.Value)
	case opLt:
		condition = c.Field.Lt(c.Value)
	case opLte:
		condition = c.Field.Lte(c.Value)
	case opGt:
		condition = c.Field.Gt(c.Value)
	case opGte:
		condition = c.Field.Gte(c.Value)
	case opIsNotNull:
		condition = c.Field.IsNotNull()
	case opIsNull:
		condition = c.Field.IsNull()
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

func (oe OrCondition) Condition() (exp.Expression, error) {
	var exps []exp.Expression
	for _, cond := range oe.Conditions {
		expr, err := cond.Condition()
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
