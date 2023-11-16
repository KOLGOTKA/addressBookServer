package pkg

import (
	"runtime"

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
	where string
	what string
}

func (er *Errorer) Regist(error_place string) {
	er.where = error_place
}
