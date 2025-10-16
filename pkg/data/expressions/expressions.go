package expressions

import (
	"fmt"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
)

type SeparatedColumnAndValueExpression struct {
	separator string
	column    string
	value     any
}

func (s SeparatedColumnAndValueExpression) String() string {
	return fmt.Sprintf("%s%s?", s.column, s.separator)
}

func (s SeparatedColumnAndValueExpression) Values() []any {
	return []any{s.value}
}

type EqualExpression struct {
	SeparatedColumnAndValueExpression
}

func Equal(column string, value any) *EqualExpression {
	return &EqualExpression{
		SeparatedColumnAndValueExpression: SeparatedColumnAndValueExpression{
			separator: "=",
			column:    column,
			value:     value,
		},
	}
}

type NotEqualExpression struct {
	SeparatedColumnAndValueExpression
}

func NotEqual(column string, value any) *NotEqualExpression {
	return &NotEqualExpression{
		SeparatedColumnAndValueExpression: SeparatedColumnAndValueExpression{
			separator: "<>",
			column:    column,
			value:     value,
		},
	}
}

type GreaterThanExpression struct {
	SeparatedColumnAndValueExpression
}

func GreaterThan(column string, value any) *GreaterThanExpression {
	return &GreaterThanExpression{
		SeparatedColumnAndValueExpression: SeparatedColumnAndValueExpression{
			separator: ">",
			column:    column,
			value:     value,
		},
	}
}

type GreaterThanOrEqualExpression struct {
	SeparatedColumnAndValueExpression
}

func GreaterThanOrEqual(column string, value any) *GreaterThanOrEqualExpression {
	return &GreaterThanOrEqualExpression{
		SeparatedColumnAndValueExpression: SeparatedColumnAndValueExpression{
			separator: ">=",
			column:    column,
			value:     value,
		},
	}
}

type LessThanExpression struct {
	SeparatedColumnAndValueExpression
}

func LessThan(column string, value any) *LessThanExpression {
	return &LessThanExpression{
		SeparatedColumnAndValueExpression: SeparatedColumnAndValueExpression{
			separator: "<",
			column:    column,
			value:     value,
		},
	}
}

type LessThanOrEqualExpression struct {
	SeparatedColumnAndValueExpression
}

func LessThanOrEqual(column string, value any) *LessThanOrEqualExpression {
	return &LessThanOrEqualExpression{
		SeparatedColumnAndValueExpression: SeparatedColumnAndValueExpression{
			separator: "<=",
			column:    column,
			value:     value,
		},
	}
}

/////////////////////////////

type ColumnLiteralExpression struct {
	prefix string
	column string
	suffix string
}

func (c ColumnLiteralExpression) String() string {
	return fmt.Sprintf("%s%s%s", c.prefix, c.column, c.suffix)
}

func (ColumnLiteralExpression) Values() []any {
	return []any{}
}

type NullExpression struct {
	ColumnLiteralExpression
}

func Null(column string) *NullExpression {
	return &NullExpression{
		ColumnLiteralExpression: ColumnLiteralExpression{
			prefix: "",
			column: column,
			suffix: " IS NULL",
		},
	}
}

type NotNullExpression struct {
	ColumnLiteralExpression
}

func NotNull(column string) *NotNullExpression {
	return &NotNullExpression{
		ColumnLiteralExpression: ColumnLiteralExpression{
			prefix: "",
			column: column,
			suffix: " IS NOT NULL",
		},
	}
}

/////////////////////////////

type AnyExpression struct {
	column     string
	values     []any
	anyBuilder data.AnyBuilder
}

func Any(column string, values []any, anyBuilder data.AnyBuilder) *AnyExpression {
	return &AnyExpression{
		column:     column,
		values:     values,
		anyBuilder: anyBuilder,
	}
}

func (e AnyExpression) String() string {
	return e.anyBuilder.String(e.column, e.values)
}

func (e AnyExpression) Values() []any {
	return e.anyBuilder.Values(e.values)
}
