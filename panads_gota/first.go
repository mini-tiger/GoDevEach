package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gota/gota/dataframe"
)

// https://pkg.go.dev/github.com/go-gota/gota/dataframe
// https://github.com/go-gota/gota
func main() {
	csvStr := `
Country,Date,Age,Amount,Id
"United States",2012-02-01,50,112.1,01234
"United States",2012-02-01,32,321.31,54321
"United Kingdom",2012-02-01,17,18.2,12345
"United States",2012-02-01,32,321.31,54322
"United Kingdom",2012-02-01,NA,18.2,12345
"United States",2012-02-01,32,321.31,54323
"United States",2012-02-01,32,321.31,54324
Spain,2012-02-01,66,555.42,00241
`
	df := dataframe.ReadCSV(strings.NewReader(csvStr))
	//fmt.Println(df)
	// Group
	groups := df.GroupBy("Country")                                                                                                          // Group by column "key1", and column "key2"
	aggre := groups.Aggregation([]dataframe.AggregationType{dataframe.Aggregation_MAX, dataframe.Aggregation_MIN}, []string{"Amount", "Id"}) // Maximum value in column "values",  Minimum value in column "values2"
	fmt.Println(aggre)

	df1 := df.Rename("Origin", "Country"). // 列重命名
						Filter(dataframe.F{Colname: "Age", Comparator: "<", Comparando: 50}). // filter
						Filter(dataframe.F{Colname: "Origin", Comparator: "==", Comparando: "United States"}).
						Select([]string{"Id", "Origin", "Date", "Age"}). // show fields
						Subset([]int{1, 3})                              // 返回 1,2  两条记录

	if df.Err != nil {
		log.Fatal("Oh noes!")
	}
	fmt.Println(df1)
}
