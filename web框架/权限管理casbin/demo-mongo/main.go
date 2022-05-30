package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/mongodb-adapter/v3"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// https://github.com/casbin/mongodb-adapter/blob/master/adapter_test.go
// http://mrlch.cn/archives/1421
func main() {
	// Initialize a MongoDB adapter with NewAdapterWithClientOption:
	// The adapter will use custom mongo client options.
	// custom database name.
	// default collection name 'casbin_rule'.
	mongoClientOption := mongooptions.Client().ApplyURI("mongodb://cc:cc@172.22.50.25:27021/?authSource=cmdb&authMechanism=SCRAM-SHA-256")
	databaseName := "casbin"
	a, err := mongodbadapter.NewAdapterWithClientOption(mongoClientOption, databaseName)
	// Or you can use NewAdapterWithCollectionName for custom collection name.
	if err != nil {
		panic(err)
	}

	e, err := casbin.NewEnforcer("/data/work/go/GoDevEach/web框架/权限管理casbin/demo-mongo/rbac_model.conf", a)
	if err != nil {
		panic(err)
	}

	// Load the policy from DB.
	e.LoadPolicy()

	// xxx 添加策略名 admin 对 data1 的read,write权限
	// xxx 重复 添加 return false
	log.Println("---------add----------------")
	AddPolicy(e)

	log.Println("----------关联 策略 ---------------")
	//// xxx 添加组gname1 关联 策略 admin 权限
	result, err := e.AddGroupingPolicy("gname1", "admin")
	if err != nil {
		panic(err)
	}

	result, err = e.AddGroupingPolicies([][]string{[]string{"gname2", "admin"}})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Check the permission.
	log.Println("-------------check-------------")
	CheckPolicy(e)

}

func AddPolicy(e *casbin.Enforcer) {

	// xxx 重复 添加 return false
	result, err := e.AddPolicy("admin", "data1", "read")
	if err != nil {
		panic(err)
	}
	result, err = e.AddPolicy("admin", "data1", "write")
	if err != nil {
		panic(err)
	}
	fmt.Printf("添加策略名 admin :%v\n", result)

	// Save the policy back to DB.
	e.SavePolicy()
}
func CheckPolicy(e *casbin.Enforcer) {

	checkbool, err := e.Enforce("admin", "data1", "read")
	if err != nil {
		panic(err)
	}
	fmt.Println(checkbool)

	checkbool, err = e.Enforce("gname1", "data1", "write")
	if err != nil {
		panic(err)
	}
	fmt.Println(checkbool)
	// 策略名 检查 权限
	fmt.Println(e.HasPolicy("admin", "data1", "write"))
	// 组名 包含 策略名
	fmt.Println(e.HasGroupingPolicy("gname1", "admin"))

	//获取策略中的指定 角色继承规则
	fmt.Println(e.GetFilteredNamedGroupingPolicy("g", 0, "gname1"))

	//获取策略中的所有角色继承规则
	fmt.Println(e.GetNamedGroupingPolicy("g"))

}
