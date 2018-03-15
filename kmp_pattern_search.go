package util

// PatternMatchOfString 字符传匹配函数，使用KMP算法
func PatternMatchOfString(pattern string, target string) int {
	nextArray := ComputeNextArray(pattern)
	patternIndex := 0
	targetIndex := 0
	for patternIndex < len(pattern) && targetIndex < len(target) {
		if target[targetIndex] == pattern[patternIndex] {
			targetIndex++
			patternIndex++
		} else if patternIndex == 0 {
			targetIndex++
		} else {
			patternIndex = nextArray[patternIndex-1] + 1
		}
	}
	if patternIndex == len(pattern) {
		return targetIndex - patternIndex
	}
	return -1
}

// ComputeNextArray KMP算法中计算pattern字符串的overlay数组
func ComputeNextArray(pattern string) []int {
	retNextArray := []int{}
	for i := 0; i < len(pattern); i++ {
		retNextArray = append(retNextArray, -1)
	}

	for idx := range pattern {
		if idx == 0 {
			retNextArray[idx] = -1
		} else {
			preStoreForward := retNextArray[idx-1]
			for preStoreForward >= 0 && pattern[idx] != pattern[preStoreForward+1] {
				preStoreForward = retNextArray[preStoreForward]
			}
			if pattern[idx] == pattern[preStoreForward+1] {
				retNextArray[idx] = preStoreForward + 1
			} else {
				retNextArray[idx] = -1
			}
		}
	}
	return retNextArray
}

// func main() {
// 	pattern := "annacanna"
// 	overlayArray := ComputeNextArray(pattern)
// 	target := "annbcdanacadsannannacanna"
// 	matchIndex := PatternMatchOfString(pattern, target)
// 	fmt.Println(overlayArray, " match index:", matchIndex)
// }
