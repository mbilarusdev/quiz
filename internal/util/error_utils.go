package util

import "strings"

func CheckDublicateErr(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
