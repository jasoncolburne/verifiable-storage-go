package orderings

import "fmt"

type AscendingOrdering struct {
	column string
}

func Ascending(column string) *AscendingOrdering {
	return &AscendingOrdering{column: column}
}

func (o AscendingOrdering) String() string {
	return fmt.Sprintf("ORDER BY %s ASC", o.column)
}

type DescendingOrdering struct {
	column string
}

func Descending(column string) *DescendingOrdering {
	return &DescendingOrdering{column: column}
}

func (o DescendingOrdering) String() string {
	return fmt.Sprintf("ORDER BY %s DESC", o.column)
}
