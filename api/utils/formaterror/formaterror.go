package formaterror

import (
	"errors"
	"fmt"
)

func FormatError(field string) error {
	return errors.New(fmt.Sprintf("Error on %s field", field))
}
