package paw

import "strings"

// isInArray 判斷目標字串是否是在陣列中
func isInArray(list *[]string, s string) (isIn bool) {

	if len(*list) == 0 {
		return false
	}

	isIn = false
	for _, f := range *list {

		if f == s {
			isIn = true
			break
		}
	}

	return isIn
}

// isInSuffix 判斷目標字串的末尾是否含有陣列中指定的字串
func isInSuffix(list *[]string, s string) (isIn bool) {

	isIn = false
	for _, f := range *list {
		if strings.TrimSpace(f) != "" && strings.HasSuffix(s, f) {
			isIn = true
			break
		}
	}

	return isIn
}

// isAllEmpty 判斷陣列各元素是否是空字串或空格
func isAllEmpty(list *[]string) (isEmpty bool) {

	if len(*list) == 0 {
		return true
	}

	isEmpty = true
	for _, f := range *list {
		if strings.TrimSpace(f) != "" {
			isEmpty = false
			break
		}
	}

	return isEmpty
}
