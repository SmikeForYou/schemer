package schemer

import (
	"github.com/SmikeForYou/dp"
	"reflect"
)

// MetadataRow is a struct that holds full information about a column in a table
type MetadataRow struct {
	TableCatalog          string
	TableSchema           string
	TableName             string
	ColumnName            string
	DataType              string
	IsNullable            bool
	ForeignKeyTableSchema *string
	ForeignKeyTableName   *string
	ForeignKeyColumnName  *string
	IsPrimaryKey          bool
	ReflectType           reflect.Type
}

type Metadata []MetadataRow

func (m Metadata) BuildDBTree() DB {
	db := make(DB)
	for _, catalogGrouped := range dp.GroupBy(m, func(row MetadataRow) any {
		return row.TableCatalog
	}) {
		catalog := &Catalog{
			Name:    catalogGrouped.Key.(string),
			Schemas: make(map[string]*Schema),
		}
		db[catalog.Name] = catalog
		for _, schemaGrouped := range dp.GroupBy(catalogGrouped.Group, func(row MetadataRow) any {
			return row.TableSchema
		}) {
			schema := &Schema{
				Name:    schemaGrouped.Key.(string),
				Tables:  make(map[string]*Table),
				Catalog: catalog,
			}
			catalog.Schemas[schema.Name] = schema
			for _, tableGrouped := range dp.GroupBy(schemaGrouped.Group, func(row MetadataRow) any {
				return row.TableName
			}) {
				table := &Table{
					Name:    tableGrouped.Key.(string),
					Columns: make(map[string]*Column),
					Schema:  schema,
				}
				schema.Tables[table.Name] = table
				for _, row := range tableGrouped.Group {
					column := &Column{
						Name:         row.ColumnName,
						Type:         row.DataType,
						ReflectType:  row.ReflectType,
						IsNullable:   row.IsNullable,
						IsPrimaryKey: row.IsPrimaryKey,
						Table:        table,
					}
					table.Columns[column.Name] = column
				}
			}
		}
	}
	for _, fk := range m {
		if fk.ForeignKeyTableSchema != nil {
			fkTable := db.GetCatalog(fk.TableCatalog).GetSchema(fk.TableSchema).GetTable(fk.TableName)
			fkColumn := fkTable.GetColumn(fk.ColumnName)
			fkTable = db.GetCatalog(*fk.ForeignKeyTableSchema).GetSchema(*fk.ForeignKeyTableSchema).GetTable(*fk.ForeignKeyTableName)
			fkColumn.ForeignKey = fkTable.GetColumn(*fk.ForeignKeyColumnName)
		}
	}
	return db
}

// Column is a struct that holds full information about a column in a table
type Column struct {
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	ReflectType  reflect.Type `json:"-"`
	IsNullable   bool         `json:"is_nullable"`
	IsPrimaryKey bool         `json:"is_primary_key"`
	ForeignKey   *Column      `json:"foreign_key"`
	Table        *Table       `json:"-"`
}

// Table is a struct that holds full information about a table in a schema
type Table struct {
	Name    string             `json:"name"`
	Columns map[string]*Column `json:"columns"`
	Schema  *Schema            `json:"-"`
}

func (t *Table) GetColumn(name string) *Column {
	return t.Columns[name]
}

// Schema is a struct that holds full information about a schema in a catalog
type Schema struct {
	Name    string            `json:"name"`
	Tables  map[string]*Table `json:"tables"`
	Catalog *Catalog          `json:"-"`
}

func (s *Schema) GetTable(name string) *Table {
	return s.Tables[name]
}

// Catalog is a struct that holds full information about a catalog
type Catalog struct {
	Name    string             `json:"name"`
	Schemas map[string]*Schema `json:"schemas"`
}

// GetSchema returns a schema from the catalog
func (c *Catalog) GetSchema(name string) *Schema {
	return c.Schemas[name]
}

// DB is a map of catalogs
type DB map[string]*Catalog

func (d *DB) GetCatalog(name string) *Catalog {
	return (*d)[name]
}
