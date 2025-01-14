package qb

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type OffsetPagingBuilder struct {
	sortBy       []string
	sortOrder    sortOrder
	limit        int
	currentPage  int
	totalRecords int
}

func NewOffsetPaging(limit int, currentPage, totalRecords int, sortOrder string, sortBy ...string) *OffsetPagingBuilder {
	return &OffsetPagingBuilder{
		sortBy:       sortBy,
		sortOrder:    sortOrderFromString(sortOrder),
		limit:        limit,
		currentPage:  currentPage,
		totalRecords: totalRecords,
	}
}

func (p *OffsetPagingBuilder) Build(tx *gorm.DB) *gorm.DB {
	orders := make([]string, len(p.sortBy))
	for i := range p.sortBy {
		orders[i] = fmt.Sprintf("%s %s", p.sortBy[i], p.sortOrder)
	}
	if len(orders) == 0 {
		orders = []string{"id asc"}
	}

	return tx.
		Order(strings.Join(orders, ",")).
		Limit(p.limit).
		Offset((p.currentPage - 1) * p.limit)
}

type OffsetPaging struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	PerPage     int `json:"per_page"`
}

func (p *OffsetPagingBuilder) Paging() any {
	totalPage := p.totalRecords / p.limit
	if p.totalRecords%p.limit != 0 {
		totalPage++
	}

	return OffsetPaging{
		CurrentPage: p.currentPage,
		TotalPage:   totalPage,
		PerPage:     p.limit,
	}
}
