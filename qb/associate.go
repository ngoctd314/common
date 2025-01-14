package qb

import "gorm.io/gorm"

func Preload(table string, conds ...*Cond) *preload {
	if notEmptyString(table) {
		return &preload{
			table: table,
			cond:  And(conds...),
		}
	}
	return nil
}

type preload struct {
	cond  *Cond
	table string
}

func (p *preload) Build(tx *gorm.DB) *gorm.DB {
	if notEmptyString(p.table) {
		tx = tx.Preload(p.table, func(preloadTx *gorm.DB) *gorm.DB {
			if p.cond != nil {
				preloadTx = p.cond.Build(preloadTx)
			}
			return preloadTx
		})
	}

	return tx
}

func Join(table string, conds ...*Cond) *join {
	if notEmptyString(table) {
		return &join{
			cond:  And(conds...),
			table: table,
		}
	}
	return nil
}

type join struct {
	cond  *Cond
	table string
}

func (p *join) Build(tx *gorm.DB) *gorm.DB {
	if notEmptyString(p.table) {
		tx = tx.Joins(p.table, func(joinTx *gorm.DB) *gorm.DB {
			if p.cond != nil {
				joinTx = p.cond.Build(joinTx)
			}
			return joinTx
		})
	}

	return tx
}
