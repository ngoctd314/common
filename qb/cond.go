package qb

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Cond struct {
	operator string
	query    string
	params   []any
}

func (cond *Cond) Build(tx *gorm.DB) *gorm.DB {
	if notEmptyString(cond.query) {
		tx = tx.Where(cond.query, cond.params...)
	}

	return tx
}

func (cond *Cond) not() *Cond {
	sb := strings.Builder{}
	sb.WriteString("NOT(")
	sb.WriteString(cond.query)
	sb.WriteString(")")

	cond.query = sb.String()
	cond.operator = ""
	return cond
}

func (cond *Cond) append(operator string, conditions ...*Cond) *Cond {
	// wrap previous condition in () if needed
	builder := strings.Builder{}
	if cond.operator != operator && notEmptyString(cond.operator) &&
		notEmptyString(cond.query) {
		cond.query = wrapQuery(cond.query)
	}
	builder.WriteString(cond.query)
	cond.operator = operator
	for _, condition := range conditions {
		if notEmptyString(condition.query) {
			if condition.operator != operator && notEmptyString(condition.operator) {
				condition.query = wrapQuery(condition.query)
			}
			if builder.Len() > 0 {
				builder.WriteString(operator)
			}
			builder.WriteString(condition.query)
			cond.params = append(cond.params, condition.params...)
		}
	}
	cond.query = builder.String()
	return cond
}

func wrapQuery(query string) string {
	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(query)
	sb.WriteString(")")

	return sb.String()
}

// Eq represents "field = value".
func Eq(field string, value interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" = ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// NotEq represents "field <> value".
func NotEq(field string, value any) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" <> ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// Gt represents "field > value".
func Gt(field string, value interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" > ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// Gte represents "field >= value".
func Gte(field string, value interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" >= ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// Lt represents "field < value".
func Lt(field string, value interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" < ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// Lte represents "field <= value".
func Lte(field string, value interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" <= ?")

	return &Cond{
		query:  sb.String(),
		params: []interface{}{value},
	}
}

func appendSliceWhereIN(value any, values ...any) []any {
	var results []any
	if value == nil {
		return values
	}

	rv := reflect.ValueOf(value)
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		for i := range rv.Len() {
			rvi := rv.Index(i)
			if rvi.CanInterface() {
				results = append(results, rvi.Interface())
			}
		}
	} else {
		results = append(results, value)
	}

	results = append(results, values...)

	return results
}

// In represents "field IN (value...)".
// Examples:
// + In("id", []int{1,2,3})
// + In("id", []any{1,2,3})
// + In("id", 1,2,3)
func In(field string, value any, values ...any) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" IN (?)")

	return &Cond{
		query:  sb.String(),
		params: []any{appendSliceWhereIN(value, values...)},
	}
}

// NotIn represents "field NOT IN (value...)".
func NotIn(field string, values ...interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" NOT IN (?)")

	return &Cond{
		query:  sb.String(),
		params: []any{values},
	}
}

// Like represents "field LIKE value".
func Like(field string, value string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" LIKE ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// NotLike represents "field NOT LIKE value".
func NotLike(field string, value string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" NOT LIKE ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// IsNull represents "field IS NULL".
func IsNull(field string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" IS NULL")

	return &Cond{
		query:  sb.String(),
		params: []any{},
	}
}

// IsNotNull represents "field IS NOT NULL".
func IsNotNull(field string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" IS NOT NULL")

	return &Cond{
		query:  sb.String(),
		params: []any{},
	}
}

// Between represents "field BETWEEN lower AND upper".
func Between(field string, lower, upper string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" BETWEEN ? AND ?")

	return &Cond{
		query:  sb.String(),
		params: []any{lower, upper},
	}
}

// NotBetween represents "field NOT BETWEEN lower AND upper".
func NotBetween(field string, lower, upper string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString("NOT BETWEEN ? AND ?")

	return &Cond{
		query:  sb.String(),
		params: []any{lower, upper},
	}
}

// And will Join simple a slice of condition into a condition with AND WiseFunc for where statement
func And(conds ...*Cond) *Cond {
	object := &Cond{}
	return object.append(" AND ", conds...)
}

// Raw represents and raw query
func Raw(query string, params []interface{}) *Cond {
	object := &Cond{
		query:  query,
		params: params,
	}
	return object
}

// Or will Join simple a slice of condition into a condition with OR WiseFunc for where statement
func Or(conds ...*Cond) *Cond {
	object := &Cond{}
	return object.append(" OR ", conds...)
}

func Not(cond *Cond) *Cond {
	if cond == nil {
		return new(Cond)
	}

	sb := strings.Builder{}
	sb.WriteString("NOT(")
	sb.WriteString(cond.query)
	sb.WriteString(")")

	cond.query = sb.String()
	cond.operator = ""
	return cond
}
