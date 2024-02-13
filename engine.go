package schemer

import (
	"database/sql"
	"schemer/dialects"
)

type dialectMap map[string]SQLDialect

var defaultDialects = dialectMap{
	"postgres": dialects.PostgresDialect{},
	"pgx":      dialects.PostgresDialect{},
	"pgx/v5":   dialects.PostgresDialect{},
}

type Engine struct {
	dialects dialectMap
	conn     *sql.DB
}

func NewEngine() *Engine {
	return &Engine{dialects: defaultDialects}
}

func (e *Engine) RegisterDialect(name string, dialect SQLDialect) {
	if e.dialects[name] != nil {
		panic("dialect already registered")
	}
	e.dialects[name] = dialect
}

func (e *Engine) GetMetadata(driver string, connString string) (Metadata, error) {
	dialect, ok := e.dialects[driver]
	if !ok {
		return nil, UnknownDriverError(driver)
	}
	var err error
	e.conn, err = sql.Open(driver, connString)
	if err != nil {
		return nil, err
	}
	return dialect.QueryMetadata(e.conn)
}

func (e *Engine) Close() {
	e.conn.Close()
}
