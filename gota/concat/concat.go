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

	var n dataframe.DataFrame

	fmt.Println(df)
	fmt.Println(len(df.Records()))
	for i := 2; i <= len(df.Records()); i++ {
		row := df.Records()[:i]
		n = dataframe.LoadRecords(row)
		fmt.Println(n)
	}

}
