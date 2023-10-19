package main

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap"
	"log"
	"strings"
)

func main() {
	//abc()
	Example_userAuthentication1()
}

func Example_userAuthentication1() {
	// The username and password we want to check
	// 用来认证的用户名和密码
	username := "tao.jun"
	//password := "Taojun207"
	fmt.Println(username)
	// 用来获取查询权限的 bind 用户.如果 ldap 禁止了匿名查询,那我们就需要先用这个帐户 bind 以下才能开始查询
	// bind 的账号通常要使用完整的 DN 信息.例如 cn=manager,dc=example,dc=org
	// 在 AD 上,则可以用诸如 mananger@example.org 的方式来 bind
	bindusername := "cn=gitadm,ou=serverusers,ou=21vianet,dc=21vianet,dc=com"
	bindpassword := "21VIAnet@G!t157"

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "21vianet.com", 3268))
	if err != nil {
		log.Fatal(err)
	}
	//l.Debug = true
	defer l.Close()

	// Reconnect with TLS
	// 建立 StartTLS 连接,这是建立纯文本上的 TLS 协议,允许您将非加密的通讯升级为 TLS 加密而不需要另外使用一个新的端口.
	// 邮件的 POP3 ,IMAP 也有支持类似的 StartTLS,这些都是有 RFC 的
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatal(err)
	}

	// First bind with a read only user
	// 先用我们的 bind 账号给 bind 上去
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Fatal(err)
	}

	// Search for the given username
	// 这样我们就有查询权限了,可以构造查询请求了
	searchRequest := ldap.NewSearchRequest(
		// 这里是 basedn,我们将从这个节点开始搜索
		//"ou=21vianet,dc=21vianet,dc=com",
		"dc=21vianet,dc=com",
		// 这里几个参数分别是 scope, derefAliases, sizeLimit, timeLimit,  typesOnly
		// 详情可以参考 RFC4511 中的定义,文末有链接
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		// 这里是 LDAP 查询的 Filter.这个例子例子,我们通过查询 uid=username 且 objectClass=organizationalPerson.

		//"(objectClass=group)",  // xxx 所有组
		//fmt.Sprintf("(&(objectClass=person))"), //xxx 所有用户
		fmt.Sprintf("(&(objectCategory=user)(sAMAccountName=%s))", ldap.EscapeFilter(username)), // xxx 根据属性筛选
		// xxx 这里是查询返回的属性,以数组形式提供.如果为空则会返回所有的属性
		nil,
		//[]string{"dn", "cn", "uid"},
		nil,
	)
	// 好了现在可以搜索了,返回的是一个数组
	//sr, err := l.Search(searchRequest)
	sr, err := l.SearchWithPaging(searchRequest, 100)
	if err != nil {
		log.Fatal(err)
	}

	for _, group := range sr.Entries {
		//if strings.Contains(group.DN,"OU=21vianet"){

		//if strings.Contains(group.DN, "陈志") {
		//	for _, singleUser := range group.Attributes {
		//		fmt.Printf("name:%s value:%s\n", singleUser.Name, singleUser.Values)
		//	}
		//	//} else {
		//	//	for _, singleUser := range group.Attributes {
		//	//		fmt.Printf("name:%s value:%s\n", singleUser.Name, singleUser.Values)
		//	//	}
		//fmt.Printf("%+v\n", group)
		//fmt.Println(group.GetAttributeValue("description"))
		//fmt.Println(group.GetAttributeValue("mail"))
		num := group.GetAttributeValue("description")
		mail := group.GetAttributeValue("mail")
		cn := group.GetAttributeValue("cn")

		if strings.Contains(mail, "tao.jun") || strings.Contains("sAMAccountName", "matt") ||
			strings.Contains(cn, "Matt") {
			fmt.Println(num)
			fmt.Println(mail)
			fmt.Println(cn)
			for _, a := range group.Attributes {
				fmt.Printf("%+v\n", a)
			}
		}

		//}

	}

	// 如果没有数据返回或者超过1条数据返回,这对于用户认证而言都是不允许的.
	// 前这意味着没有查到用户,后者意味着存在重复数据
	fmt.Println(len(sr.Entries))

}
