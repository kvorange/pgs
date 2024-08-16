package pgs

import (
	"encoding/json"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"reflect"
)

type Field[T any] struct {
	Value T
	field string
	as    string
	model *Model
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

func (f *Field[T]) MarshalJSON() ([]byte, error) {
	// Сериализуем значение поля
	data, err := json.Marshal(f.Value)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (f *Field[T]) As(as string) *Field[T] {
	f.as = as
	return f
}

func (f *Field[T]) Scan(src interface{}) error {
	v := reflect.ValueOf(&f.Value)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("value is not a pointer")
	}

	// Проверяем, что у типа, на который указывает v, есть метод Scan
	elem := v.Elem()
	if !elem.CanAddr() {
		return fmt.Errorf("cannot address element")
	}

	scanMethod := elem.Addr().MethodByName("Scan")
	if !scanMethod.IsValid() {
		return fmt.Errorf("type does not have a Scan method")
	}

	// Вызываем метод Scan с аргументом src
	scanMethod.Call([]reflect.Value{reflect.ValueOf(src)})
	return nil
}

func (f *Field[T]) In(value interface{}) Condition {
	ds, ok := value.(*SelectDataset)
	if ok {
		// goqu use Eq for op IN for nested select
		return Condition{
			Field:   f.getIdent(),
			Op:      opEq,
			Value:   ds.dataset,
			joiners: []*joiner{f.getJoiner()},
		}
	}
	return Condition{
		Field:   f.getIdent(),
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
			Field:   f.getIdent(),
			Op:      opNotEq,
			Value:   ds.dataset,
			joiners: []*joiner{f.getJoiner()},
		}
	}
	return Condition{
		Field:   f.getIdent(),
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
		Field:   f.getIdent(),
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
		Field:   f.getIdent(),
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
		Field:   f.getIdent(),
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
		Field:   f.getIdent(),
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
		Field:   f.getIdent(),
		Op:      opRegex,
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
		Field:   f.getIdent(),
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
		Field:   f.getIdent(),
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
		Field:   f.getIdent(),
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
		Field:   f.getIdent(),
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
		Field:   f.getIdent(),
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
		Field:   f.getIdent(),
		Op:      opGte,
		Value:   value,
		joiners: []*joiner{f.getJoiner()},
	}
}

func (f *Field[T]) IsNotNull() Condition {
	return Condition{
		Field: f.getIdent(),
		Op:    opIsNotNull,
		Value: nil,
	}
}

func (f *Field[T]) IsNull() Condition {
	return Condition{
		Field: f.getIdent(),
		Op:    opIsNull,
		Value: nil,
	}
}
