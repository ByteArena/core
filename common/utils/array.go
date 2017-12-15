package utils

func IsStringInArray(array []string, p string) bool {
	for _, v := range array {
		if v == p {
			return true
		}
	}

	return false
}
