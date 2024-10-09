package pgs

import (
	"encoding/json"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jackc/pgx/v5/pgtype"
)

type Field[T any] struct {
	Value T
	field string
	as    string
	model *Model
}

func (f *Field[T]) getField() string {
	return f.field
}

func (f *Field[T]) getModel() *Model {
	return f.model
}

func (f *Field[T]) getSelector() exp.AliasedExpression {
	var selector exp.IdentifierExpression
	if f.model.prefix == "" {
		selector = goqu.I(fmt.Sprintf("%s.%s", f.model.tableName, f.field))
		if f.as != "" {
			return selector.As(f.as)
		}
		return selector.As(f.field)
	}
	selector = goqu.I(
		fmt.Sprintf("%s.%s", f.model.joiner.Name, f.field),
	)
	if f.as != "" {
		return selector.As(f.as)
	}
	return selector.As(goqu.S(fmt.Sprintf("%s.%s", f.model.prefix, f.field)))
}

func (f *Field[T]) getIdent() exp.IdentifierExpression {
	var ident exp.IdentifierExpression
	if f.model.prefix == "" {
		ident = goqu.I(fmt.Sprintf("%s.%s", f.model.tableName, f.field))
		return ident
	}
	ident = goqu.I(fmt.Sprintf("%s.%s", f.model.joiner.Name, f.field))
	return ident
}

func (f *Field[T]) getSelectors() []interface{} {
	return []interface{}{f.getSelector()}
}

func (f *Field[T]) getJoiner() *joiner {
	return f.model.joiner
}

func (f *Field[T]) init(model *Model, field string) {
	f.model = model
	f.field = field
}

func (f Field[T]) MarshalJSON() ([]byte, error) {
	v := interface{}(&f.Value)

	// rewrite pgtype Float Marshal
	if floatValue, ok := v.(*pgtype.Float8); ok {
		if !floatValue.Valid {
			return json.Marshal(nil)
		}
		return json.Marshal(floatValue.Float64)
	}

	if floatValue, ok := v.(*pgtype.Float4); ok {
		if !floatValue.Valid {
			return json.Marshal(nil)
		}
		return json.Marshal(floatValue.Float32)
	}

	return json.Marshal(f.Value)
}

func (f *Field[T]) As(as string) *Field[T] {
	f.as = as
	return f
}

func (f *Field[T]) Scan(src interface{}) error {
	v := interface{}(&f.Value)
	if pgValue, ok := v.(*pgtype.Int2); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Int4); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Int8); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Float4); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Float8); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Text); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Time); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Timestamp); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Timestamptz); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Point); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Polygon); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Bool); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Bits); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Box); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Circle); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Date); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Interval); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Numeric); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Line); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Hstore); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Lseg); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.Uint32); ok {
		return pgValue.Scan(src)
	}
	if pgValue, ok := v.(*pgtype.UUID); ok {
		return pgValue.Scan(src)
	}
	return fmt.Errorf("unsupported type: %T", v)
}

func (f *Field[T]) In(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		// goqu use Eq for op IN for nested select
		return Condition{
			Field:   f,
			Op:      opEq,
			Value:   ds.dataset,
			joiners: []*joiner{f.getJoiner()},
		}
	}
	return Condition{
		Field:   f,
		Op:      opIn,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) NotIn(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		// goqu use notEq for op IN for nested select
		return Condition{
			Field:   f,
			Op:      opNotEq,
			Value:   ds.dataset,
			joiners: []*joiner{f.getJoiner()},
		}
	}
	return Condition{
		Field:   f,
		Op:      opNotIn,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) Eq(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opEq,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) NotEq(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opNotEq,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) Like(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opLike,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) NotLike(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opNotLike,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) Regex(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opRegex,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) RegexI(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opRegexI,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) NotRegex(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opNotRegex,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) NotRegexI(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opNotRegexI,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) Lt(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opLt,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) Lte(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opLte,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) Gt(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opGt,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) Gte(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		value = ds.dataset
	}
	return Condition{
		Field:   f,
		Op:      opGte,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) IsNotNull() Condition {
	return Condition{
		Field: f,
		Op:    opIsNotNull,
		Value: nil,
	}
}

func (f *Field[T]) IsNull() Condition {
	return Condition{
		Field: f,
		Op:    opIsNull,
		Value: nil,
	}
}
