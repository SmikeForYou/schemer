package dialects

import (
	"database/sql"
	"fmt"
	schemer "schemer"
)

var PostgresMetadataDefaultQuery = `
	SELECT c.table_catalog,
	       c.table_schema,
	       c.table_name,
	       c.column_name,
	       c.udt_name as data_type,
	       c.is_nullable,
	       fk.foreign_table_schema ,
	       fk.foreign_table_name, 
	       fk.foreign_column_name,
	       (pk is not null) as is_primary_key
	FROM information_schema.columns c
	LEFT JOIN (
		SELECT tc.table_schema,
			   tc.constraint_name,
			   tc.table_name,
			   kcu.column_name,
			   ccu.table_schema AS foreign_table_schema,
			   ccu.table_name   AS foreign_table_name,
			   ccu.column_name  AS foreign_column_name,
			   tc.constraint_type as constraint_type		
		FROM information_schema.table_constraints AS tc
				 JOIN information_schema.key_column_usage AS kcu
					  ON tc.constraint_name = kcu.constraint_name
						  AND tc.table_schema = kcu.table_schema
				 JOIN information_schema.constraint_column_usage AS ccu
					  ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY'
	) fk on c.table_schema = fk.table_schema and c.table_name = fk.table_name and c.column_name = fk.column_name
	LEFT JOIN (
	    SELECT tc.table_schema,
			   tc.constraint_name,
			   tc.table_name,
			   kcu.column_name
		FROM information_schema.table_constraints AS tc
				 JOIN information_schema.key_column_usage AS kcu
					  ON tc.constraint_name = kcu.constraint_name
						  AND tc.table_schema = kcu.table_schema
				 JOIN information_schema.constraint_column_usage AS ccu
					  ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'PRIMARY KEY'
	) pk on c.table_schema = pk.table_schema and c.table_name = pk.table_name and c.column_name = pk.column_name
	where c.table_schema not in ('pg_catalog', 'information_schema')
`

type PostgresDialect struct{}

func (p PostgresDialect) QueryMetadata(conn *sql.DB) (schemer.Metadata, error) {
	rows, err := conn.Query(PostgresMetadataDefaultQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	metadata := schemer.Metadata(make([]schemer.MetadataRow, 0))
	var nullable string
	for rows.Next() {
		var entry schemer.MetadataRow
		err = rows.Scan(
			&entry.TableCatalog,
			&entry.TableSchema,
			&entry.TableName,
			&entry.ColumnName,
			&entry.DataType,
			&nullable,
			&entry.ForeignKeyTableSchema,
			&entry.ForeignKeyTableName,
			&entry.ForeignKeyColumnName,
			&entry.IsPrimaryKey,
		)
		if err != nil {
			return nil, err
		}
		if nullable == "YES" {
			entry.IsNullable = true
		} else {
			entry.IsNullable = false
		}
		metadata = append(metadata, entry)
	}
	for _, simpleMetadata := range newSimpleMetadataSlice(metadata).Tables() {
		columnTypes, err := p.selectColumnTypes(
			conn,
			simpleMetadata.tableCatalog,
			simpleMetadata.tableSchema,
			simpleMetadata.tableName,
		)
		if err != nil {
			return nil, err
		}
		for _, columntType := range columnTypes {
			for _, m := range metadata {
				if m.TableCatalog == simpleMetadata.tableCatalog &&
					m.TableSchema == simpleMetadata.tableSchema &&
					m.TableName == simpleMetadata.tableName &&
					m.ColumnName == columntType.Name() {
					m.ReflectType = columntType.ScanType()
				}
			}
		}
	}

	return metadata, nil
}

func (p PostgresDialect) selectColumnTypes(conn *sql.DB, catalog, schema, table string) ([]*sql.ColumnType, error) {
	query := fmt.Sprintf("SELECT * FROM \"%s\".\"%s\".\"%s\" LIMIT 0", catalog, schema, table)
	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	return rows.ColumnTypes()

}
