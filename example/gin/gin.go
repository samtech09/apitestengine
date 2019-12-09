package ginrest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Age  int
	Name string
}

func startGin() {
	r := SetupRouter()
	fmt.Println("Server listening on port 3333...")
	http.ListenAndServe(":3333", r)
}

func SetupRouter() *gin.Engine {
	//to supress debug info from gin
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.New()
	//r.Use(gin.Logger)
	r.Use(gin.Recovery())

	r.GET("/", Index)
	r.GET("/param/:key", ParamTest)
	r.GET("/getjson", ReturnJSON)
	r.POST("/postjson", PostJSON)

	return r
}

func Index(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("Bare minimum API server in go with gin router"))
}

// func Hello(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "<html><body><h1>Hello from golang!</h1></body></html>")
// }

func ParamTest(c *gin.Context) {
	key := c.Param("key")
	c.JSON(http.StatusOK, fmt.Sprintf("key = %s", key))
	//c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf("key = %s", key)))
}

func ReturnJSON(c *gin.Context) {
	p := Person{}
	p.Age = 16
	p.Name = "Mohan"
	c.JSON(http.StatusOK, p)
}

func PostJSON(c *gin.Context) {
	p := Person{}
	err := c.BindJSON(&p)
	if err != nil {
		c.Error(fmt.Errorf("JSON bind failed"))
		return
	}
	c.JSON(http.StatusOK, "OK")
}
