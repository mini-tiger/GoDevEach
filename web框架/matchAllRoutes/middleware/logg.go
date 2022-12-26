package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Logmiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		//Start timer
		start := time.Now()
		//fmt.Println(start)
		path := c.Request.URL.Path
		//raw := c.Request.URL.RawQuery

		// Process request
		// xxx c.Next以上是请求时的运行
		c.Next() // 注意 next()方法的作用是跳过该调用链去直接后面的中间件以及api路由
		// xxx c.Next以下是返回时的运行
		//end:=time.Now()
		//fmt.Println(end)
		var b interface{}
		b, _ = c.Get("reqbody")

		log.Printf("%s | %s | %s | %d | %s | reqBody: %v",
			path,
			c.ClientIP(),
			c.Request.Method,
			c.Writer.Status(),
			time.Now().Sub(start),
			b,
		)

	}
	// Log only when path is not being skipped
	//
	//	param := LogFormatterParams{
	//		Request: c.Request,
	//		isTerm:  isTerm,
	//		Keys:    c.Keys,
	//	}
	//
	//	// Stop timer
	//	param.TimeStamp = time.Now()
	//	param.Latency = param.TimeStamp.Sub(start)
	//
	//	param.ClientIP = c.ClientIP()
	//	param.Method = c.Request.Method
	//	param.StatusCode = c.Writer.Status()
	//	param.ErrorMessage = c.Errors.ByType(ErrorTypePrivate).String()
	//
	//	param.BodySize = c.Writer.Size()
	//
	//	if raw != "" {
	//		path = path + "?" + raw
	//	}
	//
	//	param.Path = path
	//
	//	fmt.Fprint(out, formatter(param))

}
