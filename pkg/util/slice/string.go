package array

func RemoveSliceElement(l []string, target string) []string {
	exec := true
	rangeSlice(l, &exec, func(i int, str string) error {

		if target != str {
			return nil
		}

		if hasOneElement(l) {
			l = []string{}
		} else {

			if isFirstElement(l, i) {
				l = l[i+1:]
			} else if isLastElement(l, i) {
				l = l[:i]
			} else {
				l = append(l[:i], l[i+1:]...)
			}
		}

		exec = false
		return nil
	})

	return l
}

func rangeSlice(l []string, execute *bool, cb func(index int, v string) error) (err error) {
	for i, v := range l {
		if !*execute {
			return nil
		}
		if err = cb(i, v); err != nil {
			return
		}
	}

	return
}

func hasOneElement(l []string) bool {
	return len(l) == 1
}
func isFirstElement(l []string, index int) bool {
	return 0 == index
}

func isLastElement(l []string, index int) bool {
	lastIndex := len(l) - 1
	return lastIndex == index
}
