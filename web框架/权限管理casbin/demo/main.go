package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"log"
)

//https://github.com/casbin/casbin/tree/master
// xxx 例子 https://darjun.github.io/2020/06/12/godailylib/casbin/
// xxx 存储数据库 https://casbin.org/docs/en/adapters
// xxx docs https://casbin.org/docs/en/policy-storage
// xxx gin middleware https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/%E5%85%B6%E4%BB%96/%E6%9D%83%E9%99%90%E7%AE%A1%E7%90%86.html

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func main() {
	e, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	check(e, "dajun", "data", "read")
	check(e, "dajun", "data", "write")
	check(e, "lizi", "data", "read")
	check(e, "lizi", "data", "write") // lizi所属角色没有write权限
	//check(e, "root", "data", "write") // root 在policy.csv中有 r.sub == "root"
	fmt.Println(e.HasNamedPolicy("p2", "root", "write"))
	fmt.Println(e.HasNamedGroupingPolicy("g2", "aa", "root"))

}
