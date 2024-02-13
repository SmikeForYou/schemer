package dialects

import (
	"schemer"
	"slices"
)

type simpleMetadata struct {
	tableCatalog string
	tableSchema  string
	tableName    string
}

type simpleMetadataSlice []simpleMetadata

func (sm simpleMetadataSlice) Tables() simpleMetadataSlice {
	tables := make(simpleMetadataSlice, 0)
	for _, entry := range sm {
		if slices.Index(tables, entry) == -1 {
			tables = append(tables, entry)
		}
	}
	return tables
}

func newSimpleMetadata(metadata schemer.MetadataRow) simpleMetadata {
	return simpleMetadata{
		tableCatalog: metadata.TableCatalog,
		tableSchema:  metadata.TableSchema,
		tableName:    metadata.TableName,
	}
}

func newSimpleMetadataSlice(metadata schemer.Metadata) simpleMetadataSlice {
	simpleMetadataEntries := make([]simpleMetadata, 0)
	for _, meta := range metadata {
		simpleMetadataEntries = append(simpleMetadataEntries, newSimpleMetadata(meta))
	}
	return simpleMetadataEntries
}
