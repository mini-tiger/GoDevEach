package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// User 一个用户拥有多种语言，使用 `user_languages` 作为连接表
type User struct {
	gorm.Model
	Name      string
	Phone     string
	Languages []*Language `gorm:"many2many:user_languages;"`
}

// Language 一种语言属于多个用户，使用 `user_languages` 作为连接表
type Language struct {
	gorm.Model
	Name  string
	Users []*User `gorm:"many2many:user_languages;"`
}

// main 多对多
func main() {
	db, err := gorm.Open("mysql", "root:hello@tcp(172.16.71.17:3306)/itgo_monitor?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		println("err", err)
	}
	defer db.Close()

	// 启用Logger，显示详细日志
	// db.LogMode(true)

	// 创建表
	// db.CreateTable(&User{})
	// db.CreateTable(&Language{})
	// 检查表是否存在
	// db.HasTable(&User{})
	// 删除表
	// db.DropTable(&User{})
	// 不存在表时创建
	db.AutoMigrate(&User{}, &Language{})

	// ========
	// = 创建 =
	// ========
	// 插入-重复创建并关联
	newUser := User{Name: "Jinzhu", Phone: "15612341234", Languages: []*Language{&Language{Name: "中文"}, &Language{Name: "英语"}}}
	db.Create(&newUser)
	fmt.Printf("%+v\n", newUser)
	// 使用语言关联用户-重复创建并关联
	// language := Language{Name: "法语", Users: []*User{&User{Name: "Jinzhu1", Phone: "19912341234"}, &User{Name: "Jinzhu2", Phone: "16612341234"}}}
	// db.Create(&language)
	// fmt.Println(language)

	// 开始已有语言集选项-不必重复添加
	// languages := []*Language{&Language{Name: "中文"}, &Language{Name: "英语"}}
	// for _, language := range languages {
	//  db.Create(language)
	// }
	// 后来用户选择已有语言ID进行关联-还是会更新关联
	// user := User{Name: "Jinzhu", Phone: "15612341234"}
	// languages := []*Language{}
	// db.Where("id IN (?)", []uint{1, 2}).Find(&languages)
	// user.Languages = languages
	// db.Create(&user)
	// fmt.Println(user)

	// 正确使用添加后关联-这个不会更新关联
	// u := &User{Name: "dispaly", Phone: "13412341234"}
	// db.Create(u)
	// 关联一个
	// language := Language{}
	// db.First(&language, 1)
	// db.Model(&u).Association("Languages").Append(language)
	// 关联多个
	// languages := []*Language{}
	// db.Where("id IN (?)", []uint{1, 2}).Find(&languages)
	// db.Model(&u).Association("Languages").Append(languages)
	// fmt.Println(u)

	// 使用语言关联用户
	// language := Language{Name: "德语"}
	// db.Create(&language)
	// language := Language{}
	// db.First(&language, 4)
	// u := User{}
	// db.First(&u, 2)
	// db.Model(&language).Association("Users").Append(u)

	// 用户对象数据
	user := User{}
	db.Where("phone=15612341234").First(&user)
	//fmt.Println(user)
	// ========
	// = 查询 =
	// ========
	// 通过 Related 使用 many to many 关联
	db.Model(&user).Related(&user.Languages, "Languages")
	// 查找匹配的关联
	// db.Debug().Model(&user).Association("Languages").Find(&user.Languages)
	// 预加载分两条查询语句
	// user := User{}
	// db.Debug().Preload("Languages").Find(&user, "id = ?", 2)
	fmt.Print(user.Name + " : ")
	for _, language := range user.Languages {
		fmt.Print(language.Name + "  ")
	}
	fmt.Println()
	// 使用语言查询关联的用户
	// language := Language{}
	// db.First(&language, 1)
	// db.Model(&language).Related(&language.Users, "Users")
	// fmt.Print(language.Name + " : ")
	// for _, user := range language.Users {
	//  fmt.Print(user.Name + "  ")
	// }

	// ========
	// = 更新 =
	// ========
	// 使用创建新关联替换当前关联
	// 单个
	// language := Language{}
	// db.First(&language, 1)
	// db.Model(&user).Association("Languages").Replace(language)
	// 多个
	// languages := []*Language{}
	// db.Where("id IN (?)", []uint{1, 2}).Find(&languages)
	// db.Model(&user).Association("Languages").Replace(languages)
	// fmt.Println(user)

	// ========
	// = 清除 =
	// ========
	// 清空对关联的引用，不会删除关联本身
	// db.Model(&user).Association("Languages").Clear()
	// 删除关联的引用，不会删除关联本身
	// &user.Languages 普通查询为空，需要关联查询得到
	// db.Model(&user).Association("Languages").Delete(&user.Languages)

	// ========
	// = 删除 =
	// ========
	// 删除时，清空引用后再执行删除
	// db.Debug().Delete(&user)
	// 默认软删除，Unscoped()记录删除
	// db.Debug().Unscoped().Delete(&user)
	// fmt.Println(user)

	// ========
	// = 总数 =
	// ========
	// 获取关联的总数
	count := db.Model(&user).Association("Languages").Count()
	fmt.Printf("\r\n关联总数：%d", count)

}
