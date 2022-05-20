goSrc := $(shell find . -name "*.go") # 执行shell

matchAllRoutes.docker: $(goSrc) # 前置条件,xxx 如果没有main.go则报错  ,前置条件通常是一组文件名，之间用空格分隔。它指定了"目标"是否重新构建的判断标准
	CGO_ENABLED=0 GOOS=linux go build -o $@  # $@指 当前目标(matchAllRoutes.docker)，就是Make命令当前构建的那个目标


.PHONY: docker # 声明docker 是"伪目标"之后，make就不会去检查是否存在一个叫做docker的文件，而是每次运行都执行对应的命令。像.PHONY这样的内置目标名还有不少
docker: Dockerfile #前置条件 文件Dockerfile
ifneq ($(dtag),) # 其比較引數“arg1”和“arg2”的值是否相同，如果不同，則為真。和“ifeq”類似。 xxx 判断变量 dtag 是否为空
	docker build --tag matchallroutes:$(dtag) .
else
	# Missing dtag, try again. Example: make docker dtag=1.2.3
	exit 1
endif
