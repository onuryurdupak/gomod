package slice

import "strings"

func FindIndex(slice []string, find string) (index int, ok bool) {
	index = -1
	ok = false
	for i, e := range slice {
		if e == find {
			index = i
			ok = true
			return
		}
	}
	return
}

func FindSubtring(slice []string, contains string) (element string, ok bool) {
	element = ""
	ok = false
	for _, e := range slice {
		if strings.Contains(e, contains) {
			element = e
			ok = true
			return
		}
	}
	return
}

func RemoveString(slice *[]string, toRemove string) (removed bool) {
	removed = false
	foundIndex := -1

	for i, e := range *slice {
		if e == toRemove {
			foundIndex = i
			removed = true
			break
		}
	}

	if !removed {
		return
	}

	tempA := (*slice)[0:foundIndex]
	tempB := (*slice)[foundIndex+1:]

	*slice = append(tempA, tempB...)
	return
}

func RemoveSubstring(slice *[]string, contains string) (completeStr string, removed bool) {
	completeStr = ""
	removed = false
	foundIndex := -1

	for i, e := range *slice {
		if strings.Contains(e, contains) {
			foundIndex = i
			completeStr = (*slice)[i]
			removed = true
			break
		}
	}

	if !removed {
		return
	}

	tempA := (*slice)[0:foundIndex]
	tempB := (*slice)[foundIndex+1:]

	*slice = append(tempA, tempB...)
	return
}
