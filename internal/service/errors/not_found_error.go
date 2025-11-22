package service_errors

import "fmt"

type NotFoundError struct {
	ID int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Entity with id=%v not found\n", e.ID)
}
