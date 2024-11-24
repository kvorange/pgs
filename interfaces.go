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
	getField() string
	getModel() *Model

	Selectable
	Ordered
}

type Selectable interface {
	getSelectors() []interface{}
	getJoiners() []*joiner
}

type Ordered interface {
	getIdent() exp.IdentifierExpression
}
