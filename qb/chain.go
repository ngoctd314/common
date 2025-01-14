package qb

import "gorm.io/gorm"

type Chain struct {
	selectField []string
	omitField   []string
	limit       *int
	offset      *int
	cond        *Cond
	builders    []Builder
}

func New() *Chain {
	return new(Chain)
}

func (c *Chain) Build(tx *gorm.DB) *gorm.DB {
	if c.selectField != nil {
		tx = tx.Select(c.selectField)
	}
	if c.omitField != nil {
		tx = tx.Omit(c.omitField...)
	}
	if c.limit != nil {
		tx = tx.Limit(*c.limit)
	}
	if c.offset != nil {
		tx = tx.Offset(*c.offset)
	}
	if c.cond != nil {
		tx = c.cond.Build(tx)
	}
	for _, b := range c.builders {
		tx = b.Build(tx)
	}

	return tx
}

func (c *Chain) Select(selectField ...string) *Chain {
	if selectField != nil {
		c.selectField = selectField
	}
	return c
}

func (c *Chain) Omit(omitField ...string) *Chain {
	if omitField != nil {
		c.omitField = omitField
	}
	return c
}

func (c *Chain) Limit(limit int) *Chain {
	c.limit = &limit
	return c
}

func (c *Chain) Offset(offset int) *Chain {
	c.offset = &offset
	return c
}

func (c *Chain) Where(cond *Cond) *Chain {
	c.cond = cond
	return c
}

func (c *Chain) Associate(builders ...Builder) *Chain {
	c.builders = builders
	return c
}
