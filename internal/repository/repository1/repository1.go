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

var someTableName = "table"

// Example of implementing repository method
func (r *LogsRepository) R1Method(ctx context.Context) (error) {
	id := 1 // bad for production!!
	q := r.sb.		
		Select("1").
		From(someTableName).
		Where(sq.Eq{"room_id": id}).
		Limit(1)

	query, args, err := q.ToSql()
	if err != nil {
		return false, err
	}

	var one int
	err = r.conn.QueryRow(ctx, query, args...).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return nil
}
