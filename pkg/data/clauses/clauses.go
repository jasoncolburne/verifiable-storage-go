package clauses

import (
	"fmt"
	"strings"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
)

type AndClause struct {
	Clause
}

func And(children []data.ClauseOrExpression) *AndClause {
	return clause(children, &AndClause{})
}

func (c AndClause) String() string {
	return c.string("AND")
}

type OrClause struct {
	Clause
}

func Or(children []data.ClauseOrExpression) *OrClause {
	return clause(children, &OrClause{})
}

func (c OrClause) String() string {
	return c.string("OR")
}

type Clause struct {
	children []data.ClauseOrExpression
}

func clause[T data.Clause](children []data.ClauseOrExpression, t T) T {
	t.SetChildren(children)
	return t
}

func (c Clause) string(separator string) string {
	children := []string{}
	for _, child := range c.children {
		children = append(children, child.String())
	}

	return fmt.Sprintf("(%s)", strings.Join(children, " "+separator+" "))
}

func (c Clause) Values() []any {
	children := []any{}

	for _, child := range c.children {
		children = append(children, child.Values()...)
	}

	return children
}

func (c *Clause) SetChildren(children []data.ClauseOrExpression) {
	c.children = children
}
