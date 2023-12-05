package pkg

import (
	// "strings"
	"github.com/pkg/errors"
)

type MyError struct {
	funcName string
}

func NewMyError(fn string) *MyError{
	myerr := &MyError{funcName: fn}
	return myerr

}

func (er *MyError) Wrap(err error, errorMessage string) error{
	if err == nil {
		return errors.Wrap(errors.New(er.funcName), errorMessage)
	}
	return errors.Wrap(err, errorMessage)
}