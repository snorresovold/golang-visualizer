package algos

import "fmt"

func InsertionSort(items []float32) []float32 {
	var n = len(items)
	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			// if the thing before j
			if items[j-1] > items[j] {
				// switch items[j-1] and items[j]
				items[j-1], items[j] = items[j], items[j-1]
				fmt.Println("switching", items[j-1], items[j])
				fmt.Println(items)
			}
			// sends j one step back
			j = j - 1
		}
	}
	return items
}
