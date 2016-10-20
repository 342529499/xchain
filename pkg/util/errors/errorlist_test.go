package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestNewErrorList(t *testing.T) {
	errs := NewErrorList()
	if errs.Error() != "" || errs.Len() != 0 {
		t.Fatal()
	}

	var err1, err2 = errors.New("test errlist err1"), errors.New("test errlist err2")
	errs.Append(err1, err2)

	if errs.Error() != fmt.Sprintf("{1:%s,2:%s}", err1.Error(), err2.Error()) || errs.Len() != 2 {
		t.Fatal(errs.Error())
	}
}
