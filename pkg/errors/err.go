package errors

import (
	"fmt"
)

// common error
type Err struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// for log
func (e Err) Error() string {
	return fmt.Sprintf("[code = %d, msg = %s]", e.Code, e.Msg)
}

// internal error
type IErr struct {
	Err
}

// code from 9xxxx to top
// msg log => origin error | ret => convert
func NewIErr(code int, msg string) error {
	return IErr{
		Err: Err{
			Code: code,
			Msg:  msg,
		},
	}
}

// bussiness error
type BErr struct {
	Err
}

// code from 1xxxx to 3xxxx
// msg log | ret => origin
func NewBErr(code int, msg string) error {
	return BErr{
		Err: Err{
			Code: code,
			Msg:  msg,
		},
	}
}

// client error
type CErr struct {
	Err
}

// code from 8xxxx to 89999
// msg log | ret => origin error
func NewCErr(code int, msg string) error {
	return CErr{
		Err: Err{
			Code: code,
			Msg:  msg,
		},
	}
}
