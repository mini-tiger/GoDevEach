package main

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  first1
 * @Version: 1.0.0
 * @Date: 2021/9/15 上午9:41
 */

func main() {

	df := dataframe.LoadRecords(
		[][]string{
			[]string{"A", "B", "C", "D"},
			[]string{"a", "4", "5.1", "true"},
			[]string{"k", "5", "7.0", "true"},
			[]string{"k", "4", "6.0", "true"},
			[]string{"a", "2", "7.1", "false"},
		},
	)
	df2 := dataframe.LoadRecords(
		[][]string{
			[]string{"A", "F", "D"},
			[]string{"1", "1", "true"},
			[]string{"4", "2", "false"},
			[]string{"2", "8", "false"},
			[]string{"5", "9", "false"},
		},
	)

	df3 := dataframe.LoadRecords(
		[][]string{
			[]string{"A", "F", "D"},
			[]string{"2", "1", "true"},
			[]string{"6", "2", "false"},
			[]string{"8", "8", "false"},
			[]string{"3", "9", "false"},
		},
	)

	// 主键是D列  内连接
	join := df.InnerJoin(df2, "D")
	fmt.Println(join)

	rbind := df2.RBind(df3)
	fmt.Println(rbind)

	//
	mean := func(s series.Series) series.Series {

		floats := s.Float()
		sum := 0.0
		for _, f := range floats {
			sum += f
		}
		//fmt.Println(s.Name,sum)
		return series.Floats(sum / float64(len(floats)))
	}
	df = df.Capply(mean) // 按列 执行函数
	//df=df.Rapply(mean)
	fmt.Println(df)

}
