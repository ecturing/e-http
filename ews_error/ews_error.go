package ews_error

import "fmt"

type E_error struct{
	Msg string
}

func (err *E_error) Error() string{
	return fmt.Sprintf("%v",err.Msg)
}