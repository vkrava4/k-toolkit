package util

func EqualSliceS(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, e1 := range s1 {
		if e1 != s2[i] {
			return false
		}
	}

	return true
}
