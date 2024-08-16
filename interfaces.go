package pgs

import "github.com/doug-martin/goqu/v9/exp"

type modelI interface {
	Init(db *DbClient, model interface{}) error
}

type fieldI interface {
	MarshalJSON() ([]byte, error)
	Scan(value interface{}) error

	init(model *Model, field string)
	getSelector() exp.AliasedExpression
	Selectable
	Ordered
}

type Selectable interface {
	getSelectors() []interface{}
	getJoiner() *joiner
}

type Ordered interface {
	getIdent() exp.IdentifierExpression
}
