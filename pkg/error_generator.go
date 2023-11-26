package pkg

import (
	// "strings"
	"github.com/pkg/errors"
)

type MyError struct {
	// Where    string
	// Funct    string
	funcName string
}

func NewMyError(fn string) *MyError{
	myerr := &MyError{funcName: fn}
	return myerr

}

func (er *MyError) Wrap(err error, errorMessage string) error{
	if err == nil {
		// if errorMessage == "" {
		// 	return errors.New(er.funcName)
		// }
		return errors.Wrap(errors.New(er.funcName), errorMessage)
	}
	return errors.Wrap(err, errorMessage)
	// strings.Join([]string{er.Where, errorPlace}, ": ")
}


// func (er *Errorer) Add(errorPlace string) {
// 	er.Where = strings.Join([]string{er.Where, errorPlace}, ": ")
// }

// func (er *Errorer) RegFunct(errorBase string) {
// 	er.Funct = errorBase
// }

// func (er *Errorer) GetError() string {
// 	result := er.Funct + er.Where
// 	return result
// }
