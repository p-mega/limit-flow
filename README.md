### go web限流中间件

目前实现了gin的中间件

使用方式：
```
go get -u github.com/p-mega/limit-flow/filter
```
```go
func main() {
	router := gin.Default()
	router.Use(filter.GinLimitFlow(2*1000, 10, 300*1000))
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	router.Run(":8000")
}
```
其中GinLimitFlow()第一个参数表示计算的时间范围，
第二个参数表示一个ip该时间范围内一个ip的最大请求次数
第三个参数表示超过频率的ip被封禁的时间，单位：（ms）
