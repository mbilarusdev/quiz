package service_errors

import "fmt"

type DuplicateError struct {
	ID int
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("Entity with id=%v already exist\n", e.ID)
}
