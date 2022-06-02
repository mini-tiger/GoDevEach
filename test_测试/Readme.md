```
这里介绍几个常用的参数：
-bench regexp 执行相应的 benchmarks，例如 -bench=.；
-cover 开启测试覆盖率；
-run regexp 只运行 regexp 匹配的函数，例如 -run=Array 那么就执行包含有 Array 开头的函数；
-v 显示测试的详细命令。


测试文件 以  _test.go 结尾 

```

### 功能
```

go test -v first_test.go
go test -v /data/work/go/GoDevEach/test_测试/

压力函数名  Test 开头
func TestSum(t *testing.T) {
```


### 压力

```
first_Benchmark_test.go

mem
go test -v -bench=. -benchmem first_Benchmark_test.go

cpu
go test -bench=. -run=none -benchmem -cpuprofile=cpu.pprof first_Benchmark_test.go

time
go test -v -bench=. -benchtime=5s first_Benchmark_test.go

count 运行次数 
go test -v -bench=. -benchtime=50x first_Benchmark_test.go

压力函数名  Benchmark_
    func Benchmark_ByteString(b *testing.B) {
	var x = []byte("Hello")
	x = append(x, []byte(" World!")...)
	for i := 0; i < b.N; i++ {
		_ = *(*string)(unsafe.Pointer(&x))
	}
}
```

