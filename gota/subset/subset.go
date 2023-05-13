package main

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
)

func main() {
	df := dataframe.LoadRecords(
		[][]string{
			{"A", "B", "C", "D"},
			{"a", "4", "5.1", "true"},
			{"k", "5", "7.0", "true"},
			{"k", "4", "6.0", "true"},
			{"a", "2", "7.1", "false"},
		},
	)
	fmt.Println(df)
	sub := df.Subset([]int{2})
	fmt.Println(sub)
	sub.Elem(0, findColIndex(sub.Names(), "C")).Set(99)
	fmt.Println(sub)
	df.Set(2, sub)
	df.Elem(1, findColIndex(sub.Names(), "A")).Set("xx")
	fmt.Println(df)
}

func findColIndex(cols []string, col string) int {
	for i := 0; i < len(cols); i++ {
		if cols[i] == col {
			return i
		}
	}

	return -1
}
