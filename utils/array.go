package utils

func ContainsString(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func ArrayRemoveElement(slice interface{}, ele interface{}) interface{} {
	aSlice := slice.([]interface{})
	for i, v := range aSlice {
		if v == ele {
			slice = append(aSlice[:i], aSlice[i+1:]...)
		}
	}
	return aSlice
}
