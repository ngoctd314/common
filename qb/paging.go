package qb

// contains the common pagination logic

type sortOrder string

const (
	sortOrderASC  sortOrder = "asc"
	sortOrderDESC sortOrder = "desc"
)

func sortOrderFromString(sortOrderStr string) sortOrder {
	switch sortOrderStr {
	case string(sortOrderASC):
		return sortOrderASC
	case string(sortOrderDESC):
		return sortOrderDESC
	default:
		return sortOrderASC
	}
}

// cursor pagination common
