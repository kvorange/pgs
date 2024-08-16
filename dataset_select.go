package pgs

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type SelectDataset struct {
	model        *Model
	dataset      *goqu.SelectDataset
	joinedTables map[string]bool
	err          error
	tx           pgx.Tx
}

func (sd *SelectDataset) Where(conditions ...Conditional) *SelectDataset {
	var exps []exp.Expression
	for _, condition := range conditions {
		cond, err := condition.Condition()
		joiners := condition.getJoiners()
		for _, joiner := range joiners {
			if joiner == nil {
				continue
			}
			_, ok := sd.joinedTables[joiner.Name]
			if !ok {
				sd.dataset = sd.dataset.LeftJoin(joiner.Table, joiner.On)
				sd.joinedTables[joiner.Name] = true
			}
		}
		sd.err = err
		exps = append(exps, cond)
	}
	sd.dataset = sd.dataset.Where(exps...)
	return sd
}

func (sd *SelectDataset) Limit(limit uint) *SelectDataset {
	sd.dataset = sd.dataset.Limit(limit)
	return sd
}

func (sd *SelectDataset) Offset(offset uint) *SelectDataset {
	sd.dataset = sd.dataset.Offset(offset)
	return sd
}

func (sd *SelectDataset) OrderAsc(fields ...Ordered) *SelectDataset {
	for _, field := range fields {
		ident := field.getIdent()
		sd.dataset = sd.dataset.OrderAppend(ident.Asc())
	}
	return sd
}

func (sd *SelectDataset) OrderDesc(fields ...Ordered) *SelectDataset {
	for _, field := range fields {
		ident := field.getIdent()
		sd.dataset = sd.dataset.OrderAppend(ident.Desc())
	}
	return sd
}

func (sd *SelectDataset) Scan(dst interface{}) error {
	if sd.err != nil {
		return sd.err
	}
	var q pgxscan.Querier
	if sd.tx != nil {
		q = sd.tx
	} else {
		q = sd.model.db.Pool
	}
	query, _, _ := sd.dataset.ToSQL()
	err := pgxscan.Select(sd.model.db.Ctx, q, dst, query)
	return err
}

func (sd *SelectDataset) ScanOne(dst interface{}) error {
	if sd.err != nil {
		return sd.err
	}
	var q pgxscan.Querier
	if sd.tx != nil {
		q = sd.tx
	} else {
		q = sd.model.db.Pool
	}
	query, _, _ := sd.dataset.ToSQL()
	err := pgxscan.Get(sd.model.db.Ctx, q, dst, query)
	return err
}

func (sd *SelectDataset) Query() string {
	query, _, _ := sd.dataset.ToSQL()
	return query
}

func (sd *SelectDataset) WithTx(tx pgx.Tx) *SelectDataset {
	sd.tx = tx
	return sd
}
