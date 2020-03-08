package cerr

import (
	"fmt"
)

type ErrDuplicatedID string

func (e ErrDuplicatedID) Error() string {
	return fmt.Sprintf("%q already exists", string(e))
}

type HttpStatusError struct {
	Code int
	Body []byte
}

func (e HttpStatusError) Error() string {
	return fmt.Sprintf("status code %d: %s", e.Code, e.Body)
}
