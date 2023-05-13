package main

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
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
	// Change column C with a new one
	mut := df.Mutate(
		series.New([]string{"a", "b", "c", "d"}, series.String, "C"),
	)
	// Add a new column E
	mut2 := df.Mutate(
		series.New(make([]int, df.Nrow()), series.Int, "E"),
	)

	fmt.Println(df)
	fmt.Println(mut)
	fmt.Println(mut2)
}
