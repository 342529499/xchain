package array

import (
	"reflect"
	"testing"
)

func TestRemoveSliceElement(t *testing.T) {

	inSlices := [][]string{
		[]string{},
		[]string{"a", "b", "c"},
	}
	inRemoveElement := [][]string{
		[]string{"x", "a"},
		[]string{"x", "a", "b", "c"},
	}

	expects := [][][]string{
		[][]string{
			[]string{},
			[]string{},
		},
		[][]string{
			[]string{"a", "b", "c"},
			[]string{"b", "c"},
			[]string{"a", "c"},
			[]string{"a", "b"},
		},
	}

	for i := 0; i < len(inSlices); i++ {
		for j := 0; j < len(inRemoveElement[i]); j++ {
			tmp := []string{}
			tmp = append(tmp, inSlices[i]...)
			tmp = RemoveSliceElement(tmp, inRemoveElement[i][j])

			if !reflect.DeepEqual(tmp, expects[i][j]) {
				t.Fatal("xxx")
			}
		}
	}

}
