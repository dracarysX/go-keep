package main

import (
	"fmt"
)

func quickSort(source []int) (result []int) {
	if len(source) <= 1 {
		result = source[:]
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
func merge(left, right []int) (result []int) {
	if len(right) == 0 {
		result = left
		return
	}
	for i, j := range left {
		if len(right) == 0 {
			result = append(result, left[i:]...)
			break
		}
		for _, y := range right {
			if y >= j {
				break
			}
			result = append(result, y)
			right = right[1:]
		}
		result = append(result, j)
	}
	if len(right) > 0 {
		result = append(result, right...)
	}
	return
}

func mergeSort(source []int) (result []int) {
	if len(source) <= 1 {
		result = source
		return
	}
	t := len(source) / 2
	left := mergeSort(source[:t])
	right := mergeSort(source[t:])
	result = merge(left, right)
	return
}

func main() {
	source := []int{10, 7, 1, 3, 5, 9, 2, 4, 8, 6}
	target := quickSort(source)
	fmt.Println(target)
	fmt.Println(mergeSort(source))
}
