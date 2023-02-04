package models

// Result struct
type Result struct {
	Stat    int         `json:"stat" example:"0"`
	Message string      `json:"message" example:"请求信息"`
	Data    interface{} `json:"data" `
}

// SimpleIssue struct
type SimpleIssue struct {
	ID int `json:"id"`
}

// Issue struct
type Issue struct {
	ID         int         `form:"id" json:"id"`
	Field      string      `form:"field" json:"field" example:"这是描述"`
	Value      string      `form:"value" json:"value" example:"这是解决方法"`
	SolvedUser string      `form:"solved_user" json:"solved_user" example:"这是解决人"`
	CreateTime string      `form:"create_time" json:"create_time"`
	Data       interface{} `form:"data" json:"data" swaggertype:"object,string" example:"key:value,key2:value2"` // map[string]interface{}
	Data1      interface{} `form:"data1" json:"data1" swaggertype:"array,string" example:"value1,value2"`        // []string

}
