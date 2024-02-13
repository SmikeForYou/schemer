package dialects

import (
	"fmt"
)

func ErrorTypeNotDefined(type_ string, meta schemer.MetadataRow) error {
	return fmt.Errorf("type %s is unknown. MetadataRow: %+v", type_, meta)
}
