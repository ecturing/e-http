package ews_error

import "fmt"

type EError struct {
	Msg string
}

func (err *EError) Error() string {
	return fmt.Sprintf("%v", err.Msg)
}
