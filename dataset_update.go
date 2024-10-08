package pgs

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type UpdateDataset struct {
	model   *Model
	dataset *goqu.UpdateDataset
	err     error
	tx      pgx.Tx
}

func (d *UpdateDataset) Where(conditions ...Conditional) *UpdateDataset {
	var exps []exp.Expression
	for _, condition := range conditions {
		cond, err := condition.Condition(true)
		d.err = err
		exps = append(exps, cond)
	}
	d.dataset = d.dataset.Where(exps...)
	return d
}

func (d *UpdateDataset) Exec() error {
	query, _, _ := d.dataset.ToSQL()
	var err error
	if d.tx != nil {
		_, err = d.tx.Exec(d.model.db.Ctx, query)
	} else {
		_, err = d.model.db.Pool.Exec(d.model.db.Ctx, query)
	}
	return err
}

func (d *UpdateDataset) WithTx(tx pgx.Tx) *UpdateDataset {
	d.tx = tx
	return d
}

func (d *UpdateDataset) Returning(fields ...fieldI) *UpdateDataset {
	var rValues []interface{}
	for _, field := range fields {
		rValues = append(rValues, field.getIdent())
	}
	d.dataset = d.dataset.Returning(rValues...)
	return d
}

func (d *UpdateDataset) Scan(dst interface{}) error {
	query, _, _ := d.dataset.ToSQL()
	var err error
	if d.tx != nil {
		err = pgxscan.Select(d.model.db.Ctx, d.tx, dst, query)
	} else {
		err = pgxscan.Select(d.model.db.Ctx, d.model.db.Pool, dst, query)
	}
	return err
}

func (d *UpdateDataset) ScanOne(dst interface{}) error {
	query, _, _ := d.dataset.ToSQL()
	var err error
	if d.tx != nil {
		err = pgxscan.Get(d.model.db.Ctx, d.tx, dst, query)
	} else {
		err = pgxscan.Get(d.model.db.Ctx, d.model.db.Pool, dst, query)
	}
	return err
}

func (d *UpdateDataset) Query() string {
	query, _, _ := d.dataset.ToSQL()
	return query
}
