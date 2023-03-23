package test_test

import (
	"GinDemoTest/modules"
	"GinDemoTest/routers"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

var _ = Describe("loginTestDesc", func() {

	var login1 modules.Login
	var login2 modules.Login

	BeforeEach(func() {
		login1 = modules.Login{
			User:     "abc",
			Password: "123",
		}
		//
		login2 = modules.Login{
			User:     "taojun",
			Password: "123",
		}
	})

	Describe("Login views Test ", func() {
		r := SetUpRouter()
		//r.GET("/", HomepageHandler)
		//req, _ := http.NewRequest("GET", "/", nil)
		routers.LoadRoute(r)

		Context("With more than 300 pages", func() {
			It("should be a novel 1", func() {
				jsonValue, _ := json.Marshal(login1)
				req, err := http.NewRequest("POST", "/hg/api/logindemo", bytes.NewBuffer(jsonValue))

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				//fmt.Println(w)
				Expect(err).NotTo(HaveOccurred())
				Expect(w.Result().StatusCode).To(Equal(200))
			})
			It("should be a novel 2", func() {
				jsonValue, _ := json.Marshal(login2)
				req, err := http.NewRequest("POST", "/hg/api/logindemo", bytes.NewBuffer(jsonValue))
				Expect(err).Should(BeNil())
				Expect(err).NotTo(HaveOccurred())
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				result, _ := ioutil.ReadAll(w.Body)
				var ret map[string]interface{}
				err = json.Unmarshal(result, &ret)
				Expect(err).Should(BeNil())
				//fmt.Println(ret)

				Expect(ret).To(HaveKey("statusText"))
			})
		})

	})
})
