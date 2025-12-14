package gcvalidation

import "regexp"

var (
	emailRegex = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
)
