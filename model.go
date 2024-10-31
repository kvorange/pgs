package pgs

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"reflect"
	"strings"
)

const separator = "__"

type joiner struct {
	ParentTable string
	Name        string
	From        string
	To          string
	Table       exp.AliasedExpression
	On          exp.JoinCondition
}

type Model struct {
	db        *DbClient
	tableName string
	fields    []fieldI
	fkModels  []*Model

	parent *Model
	asName string
	prefix string
	joiner *joiner
}

func (m *Model) Init(db *DbClient, model interface{}) error {
	m.db = db
	rValue := reflect.ValueOf(model).Elem()
	rType := rValue.Type()

	// Check model field
	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)

		if field.Type == reflect.TypeOf(Model{}) {
			m.tableName = field.Tag.Get("table")
			break
		}
		return fmt.Errorf("error in init table: not found table name")
	}

	if m.parent != nil {
		if m.parent.prefix == "" {
			m.prefix = fmt.Sprintf("%s", m.asName)
		} else {
			m.prefix = fmt.Sprintf("%s.%s", m.parent.prefix, m.asName)
		}
	}

	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)

		if field.Type == reflect.TypeOf(Model{}) {
			continue
		}

		dbTag := field.Tag.Get("db")
		if dbTag == "-" {
			continue
		}

		// Check Field[any]
		dbField, ok := rValue.Field(i).Addr().Interface().(fieldI)
		if ok {
			if dbTag == "" {
				dbField.init(m, toSnakeCase(field.Name))
			} else {
				dbField.init(m, dbTag)
			}
			m.fields = append(m.fields, dbField)
			continue
		}

		// Проверка тега fk
		fkTag := field.Tag.Get("fk")
		if fkTag != "" {
			if dbTag == "" {
				return fmt.Errorf("error in init table: not found db tag with fk %s", fkTag)
			}
			fkValues := strings.Split(fkTag, ",")
			if len(fkValues) != 2 {
				return fmt.Errorf("error in init table: uncorrect value in fk tag. Expected from_field,to_field. Goted: %s", fkTag)
			}

			fkModelInterface, ok := rValue.Field(i).Addr().Interface().(modelI)
			if !ok {
				return fmt.Errorf("error in init table: fk field %s is not a model", dbTag)
			}
			nestedModelField := reflect.ValueOf(fkModelInterface).Elem().FieldByName("Model")
			nestedModel, ok := nestedModelField.Addr().Interface().(*Model)
			if !ok {
				return fmt.Errorf("error in init table: fk field %s is not a model", dbTag)
			}
			nestedModel.asName = dbTag
			nestedModel.parent = m

			err := fkModelInterface.Init(db, fkModelInterface)
			if err != nil {
				return err
			}

			var tableAsName string
			var modelAs string
			if m.asName == "" {
				tableAsName = fmt.Sprintf("%s%s%s", m.tableName, separator, nestedModel.asName)
				modelAs = m.tableName
			} else {
				tableAsName = fmt.Sprintf("%s%s%s", m.asName, separator, nestedModel.asName)
			}
			if m.parent != nil {
				asName := m.parent.tableName
				if m.parent.asName != "" {
					asName = m.parent.asName
				}
				modelAs = fmt.Sprintf("%s%s%s", asName, separator, m.asName)
			}
			var joiner joiner
			joiner.From = fkValues[0]
			joiner.To = fkValues[1]
			joiner.Name = tableAsName
			joiner.ParentTable = m.tableName
			joiner.Table = goqu.T(nestedModel.tableName).As(tableAsName)
			joiner.On = goqu.On(
				goqu.Ex{
					fmt.Sprintf("%s.%s", modelAs, joiner.From): goqu.I(fmt.Sprintf("%s.%s", tableAsName, joiner.To)),
				},
			)
			nestedModel.joiner = &joiner
			m.fkModels = append(m.fkModels, nestedModel)
			continue
		}

		return fmt.Errorf("error in init table: unknown field %v", field.Name)
	}

	return nil
}

func (m *Model) Select(fields ...Selectable) *SelectDataset {
	dataset := goqu.From(m.tableName)
	joinedTables := make(map[string]bool)

	var selectFields []interface{}

	if len(fields) == 0 {
		selectFields = m.allSelectors()
		joiners := m.allJoiners()
		for _, joiner := range joiners {
			if joiner != nil {
				_, ok := joinedTables[joiner.Name]
				if !ok {
					dataset = dataset.LeftJoin(joiner.Table, joiner.On)
					joinedTables[joiner.Name] = true
				}
			}
		}
	}
	for _, field := range fields {
		selectFields = append(selectFields, field.getSelectors()...)
		joiner := field.getJoiner()
		if joiner != nil {
			_, ok := joinedTables[joiner.Name]
			if !ok {
				dataset = dataset.LeftJoin(joiner.Table, joiner.On)
				joinedTables[joiner.Name] = true
			}
		}
	}

	dataset = dataset.Select(selectFields...)
	return &SelectDataset{
		model:        m,
		dataset:      dataset,
		joinedTables: joinedTables,
	}
}

func (m *Model) Delete() *DeleteDataset {
	dataset := goqu.Delete(m.tableName)
	return &DeleteDataset{
		model:   m,
		dataset: dataset,
		tx:      nil,
	}
}

func (m *Model) Update(record Record) *UpdateDataset {
	values := record.toMap()
	dataset := goqu.Update(m.tableName).Set(values)
	return &UpdateDataset{
		model:   m,
		dataset: dataset,
		tx:      nil,
	}
}

func (m *Model) Insert(records ...Record) *InsertDataset {
	var rows []map[string]interface{}
	for _, record := range records {
		rows = append(rows, record.toMap())
	}
	dataset := goqu.Insert(m.tableName).Rows(rows)
	return &InsertDataset{
		model:   m,
		dataset: dataset,
		tx:      nil,
	}
}

func (m *Model) getSelectors() []interface{} {
	var selectors []interface{}
	for _, field := range m.fields {
		selectors = append(selectors, field.getSelector())
	}
	return selectors
}

func (m *Model) allSelectors() []interface{} {
	selectors := m.getSelectors()
	for _, fk := range m.fkModels {
		selectors = append(selectors, fk.allSelectors()...)
	}
	return selectors
}

func (m *Model) allJoiners() []*joiner {
	var joiners []*joiner
	joiners = append(joiners, m.joiner)
	for _, fk := range m.fkModels {
		joiners = append(joiners, fk.allJoiners()...)
	}
	return joiners
}

func (m *Model) getJoiner() *joiner {
	return m.joiner
}
