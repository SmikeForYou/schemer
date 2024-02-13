package schemer

import "database/sql"

type SQLDialect interface {
	QueryMetadata(*sql.DB) (Metadata, error)
}
