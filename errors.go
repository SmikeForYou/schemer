package schemer

import "fmt"

var UnknownDriverError = func(driver string) error {
	return fmt.Errorf("unknown driver: %s", driver)
}
