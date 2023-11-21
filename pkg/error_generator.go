package pkg

import (
	"runtime"
	"strings"
	"github.com/pkg/errors"
)

func ErrorGenerate(text string) error {
	pc, _, _, ok := runtime.Caller(1)

	if ok {
		funcInfo := runtime.FuncForPC(pc)
		if funcInfo != nil {
			return errors.Errorf("%s: %s", funcInfo.Name(), text)
		} else {
			return errors.Errorf("Unknown function: %s", text)
		}
	}
	return errors.Errorf("Unable to retrieve caller information")

}

type Errorer struct {
	Where string
	Funct string
}

func (er *Errorer) Add(error_place string) {
	er.Where = strings.Join([]string{er.Where, error_place}, ": ")
}

func (er *Errorer) RegFunct(error_base string) {
	er.Funct = error_base
}

func (er *Errorer) GetError() string{
	result := er.Funct + er.Where
	return result
}