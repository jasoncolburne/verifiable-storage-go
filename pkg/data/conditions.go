package data

type ClauseOrExpression interface {
	String() string
	Values() []any
}

type Clause interface {
	ClauseOrExpression
	SetChildren(children []ClauseOrExpression)
}

type Ordering interface {
	String() string
}
