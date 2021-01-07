package longeststring

import "fmt"

func findDuplicate(s string) map[string][]int {
	charMap := make(map[string][]int)
	for i, c := range s {
		char := string(c)
		charMap[char] = append(charMap[char], i)
	}

	return charMap
}

func longString(s1 string, s2 string) string {
	if len(s1) > len(s2) {
		return s1
	} else {
		return s2
	}
}

func findLongString(s string, begin int, end int, charMap map[string][]int) string {
	var nextLongString string

	for _, v := range charMap {
		for i := 0; i < len(v); i++ {
			if i+1 < len(v) {
				if begin <= v[i] && end >= v[i+1] {
					left := findLongString(s, begin, v[i+1]-1, charMap)
					right := findLongString(s, v[i]+1, end, charMap)
					nextLongString = longString(nextLongString, left)
					nextLongString = longString(nextLongString, right)
				}
			}
		}
	}

	if len(nextLongString) == 0 {
		return s[begin : end+1]
	} else {
		return nextLongString
	}
}

func lengthOfLongestSubstring(s string) int {
	longString := findLongString(s, 0, len(s)-1, findDuplicate(s))
	fmt.Println(longString)
	return len(longString)
}
