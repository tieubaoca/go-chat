package utils

func ContainsString(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func ArrayIntRemoveElement(slice []int, ele int) []int {
	for i, v := range slice {
		if v == ele {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
