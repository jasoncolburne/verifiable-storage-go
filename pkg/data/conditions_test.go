package data_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/data"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/data/clauses"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/data/expressions"
)

func TestComposition(t *testing.T) {
	if err := testComposition(); err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}
}

func testComposition() error {
	condition := clauses.And([]data.ClauseOrExpression{
		expressions.Equal("a", "b"),
		expressions.NotEqual("c", "d"),
		clauses.Or([]data.ClauseOrExpression{
			expressions.GreaterThanOrEqual("e", "f"),
			expressions.Null("g"),
		}),
	})

	if !strings.EqualFold(condition.String(), `(a=? AND c<>? AND (e>=? OR g IS NULL))`) {
		return fmt.Errorf("unexpected string result: %s", condition.String())
	}

	values := condition.Values()
	if len(values) != 3 {
		return fmt.Errorf("unexpected value count: %d", len(values))
	}

	expected := []string{"b", "d", "f"}
	for i, value := range values {
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("value was not a string")
		}

		if !strings.EqualFold(v, expected[i]) {
			return fmt.Errorf("unexpected value at index %d", i)
		}
	}

	return nil
}
