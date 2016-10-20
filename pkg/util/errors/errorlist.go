package errors

import "fmt"

type ErrorList interface {
	Error() string
	Len() int
	Append(...error)
}

type errorList struct {
	l []error
}

func NewErrorList() ErrorList {
	return &errorList{[]error{}}
}

func (l *errorList) Error() string {
	if len(l.l) == 0 {
		return ""
	}

	retBytes := []byte{'{'}
	for i, err := range l.l {
		retBytes = append(retBytes, []byte(fmt.Sprintf("%d:%s", i+1, err.Error()))...)
		if i < l.Len()-1 {
			retBytes = append(retBytes, []byte{','}...)
		}
	}
	retBytes = append(retBytes, []byte{'}'}...)

	return string(retBytes)
}

func (l *errorList) Len() int {
	if l == nil {
		return 0
	}

	return len(l.l)
}

func (l *errorList) Append(err ...error) {
	if len(err) > 0 {
		l.l = append(l.l, err...)
	}
}
