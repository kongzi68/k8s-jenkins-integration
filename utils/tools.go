package utils

import (
	"regexp"
	"sort"
)

// in 语句
func In(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}

// 正则匹配字符串是否包含数组中的元素
func RegMatchStringArray(target string, strArray [3]string) bool {
	for _, val := range strArray {
		pattern := regexp.MustCompile(val + `.?`)
		ret := pattern.MatchString(target)
		if ret {
			return true
		}
	}
	return false
}

// 去重
func Unique(arr []string) []string {
	sort.Strings(arr)
	var arr_len int = len(arr) - 1
	for ; arr_len > 0; arr_len-- {
		// 拿最后项与前面的各项逐个(自后向前)进行比较
		for j := arr_len - 1; j >= 0; j-- {
			if arr[arr_len] == arr[j] {
				arr = append(arr[:arr_len], arr[arr_len+1:]...)
				break
			}
		}

		/*
		   // 或拿最后项与前面的各项逐个(自前向后)进行比较
		   for j := 0; j < arr_len; j++ {
		     if arr[arr_len] == arr[j] {
		     	// fmt.Printf("arr_len=%d equals j=%d\n ", arr[arr_len], arr[j])
		     	// 如果存在重复项，则将重复项删除，并重新给数组赋值
		       arr = append(arr[:arr_len], arr[arr_len + 1:]...)
		       break
		     }
		   }
		*/
	}
	return arr
}
