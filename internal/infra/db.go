package infra

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

func NewDBConnection(dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return conn, nil
}
