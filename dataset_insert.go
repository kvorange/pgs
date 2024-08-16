package pgs

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type InsertDataset struct {
	model   *Model
	dataset *goqu.InsertDataset
	err     error
	tx      pgx.Tx
}

func (d *InsertDataset) Exec() error {
	query, _, _ := d.dataset.ToSQL()
	var err error
	if d.tx != nil {
		_, err = d.tx.Exec(d.model.db.Ctx, query)
	} else {
		_, err = d.model.db.Pool.Exec(d.model.db.Ctx, query)
	}
	return err
}

func (d *InsertDataset) WithTx(tx pgx.Tx) *InsertDataset {
	d.tx = tx
	return d
}

func (d *InsertDataset) Returning(fields ...fieldI) *InsertDataset {
	var rValues []interface{}
	for _, field := range fields {
		rValues = append(rValues, field.getIdent())
	}
	d.dataset = d.dataset.Returning(rValues...)
	return d
}

func (d *InsertDataset) Scan(dst interface{}) error {
	query, _, _ := d.dataset.ToSQL()
	var err error
	if d.tx != nil {
		err = pgxscan.Select(d.model.db.Ctx, d.tx, dst, query)
	} else {
		err = pgxscan.Select(d.model.db.Ctx, d.model.db.Pool, dst, query)
	}
	return err
}

func (d *InsertDataset) ScanOne(dst interface{}) error {
	query, _, _ := d.dataset.ToSQL()
	var err error
	if d.tx != nil {
		err = pgxscan.Get(d.model.db.Ctx, d.tx, dst, query)
	} else {
		err = pgxscan.Get(d.model.db.Ctx, d.model.db.Pool, dst, query)
	}
	return err
}

func (d *InsertDataset) Query() string {
	query, _, _ := d.dataset.ToSQL()
	return query
}
