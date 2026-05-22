package repository1

import (
	// always pass `ctx context.Context` into repos
	"context"

	// always use this query builder for sql statements
	"github.com/Masterminds/squirrel" 
)

type Repository1 struct {
	conn connection
	sb   squirrel.StatementBuilderType
}

func NewRepository1(conn connection) *Repository1 {
	return &Repository1{
		conn: conn,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *LogsRepository) R1Method(ctx context.Context) (error) {
	return nil
}
