package g

/**
 * @Author: Tao Jun
 * @Description: g
 * @File:  config
 * @Version: 1.0.0
 * @Date: 2021/4/13 下午2:57
 */

var (
	IndexName = "subject"
	TypeName  = "online"
	Servers   = []string{"http://172.16.71.31:9200"}
)

type Subject struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
}
