package utils

import "strconv"

func StringToBool(name string) bool {
	result, err := strconv.ParseBool(name)
	if err != nil {
		return false
	}
	return result
}

func StringToInt(num string) int {
	result, err := strconv.Atoi(num)
	if err != nil {
		return 0
	}
	return result
}
