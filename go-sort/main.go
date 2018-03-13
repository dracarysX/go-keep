package main

import (
	"fmt"
)

func quickSort(source []int) (result []int) {
	if len(source) <= 1 {
		result = source
		return
	}
	key := source[0]
	var left []int
	var right []int
	for _, i := range source[1:] {
		if i <= key {
			left = append(left, i)
		} else if i > key {
			right = append(right, i)
		}
	}
	left = quickSort(left)
	right = quickSort(right)
	result = append(result, left...)
	result = append(result, key)
	result = append(result, right...)
	return
}

// 合并两个有序数组
func merge(a, b []int) (result []int) {
	for i, j := range a {
		if len(b) == 0 {
			result = append(result, a[i:]...)
			break
		}
		var x, y int
		for x, y = range b {
			if y >= j {
				result = append(result, j)
				break
			}
			result = append(result, y)
		}
		b = b[x:]
	}
	return
}

func mergeSort(source []int) (result []int) {
	if len(source) <= 1 {
		result = source
		return
	}
	return
}

func main() {
	source := []int{10, 7, 1, 3, 5, 9, 2, 4, 8, 6}
	target := quickSort(source)
	fmt.Println(target)
	r := merge([]int{1, 3, 5, 7, 9}, []int{2, 4, 6, 8})
	fmt.Println(r)
}
