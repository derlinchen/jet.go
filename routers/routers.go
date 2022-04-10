package routers

import (
	"encoding/json"
	"jet/bean"
	"jet/controller"
	"log"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// SetupRouter 配置路由信息
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.NoRoute(NoFound)
	router.NoMethod(NoFound)
	router.Use(Recover)
	router.Use(Logger)
	baseDic := router.Group("/wms/baseDic")
	{
		baseDic.POST("/saveBaseDic", controller.SaveBaseDic)
		baseDic.GET("/getBaseDic", controller.GetBaseDic)
		baseDic.POST("/searchBaseDic", controller.SearchBaseDic)
		baseDic.DELETE("/deleteBaseDic", controller.DeleteBaseDic)
		baseDic.POST("/updateBaseDic", controller.UpdateBaseDic)
	}
	baseLocation := router.Group("/wms/baseLocation")
	{
		baseLocation.POST("/saveBaseLocation", controller.SaveBaseLocation)
	}
	return router
}

func NoFound(c *gin.Context) {
	bean.NewResult(c).Error("404", "未找到接口")
	return
}

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v\n", r)
			debug.PrintStack()
			bean.NewResult(c).Error("500", "服务器内部错误")
		}
	}()

	c.Next()
}

func Logger(c *gin.Context) {
	f, err := os.OpenFile("log/jet.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		_, err := os.Stat("log/")
		if os.IsNotExist(err) {
			os.MkdirAll("log/", os.ModePerm)
		}
		os.Create("log/jet.log")
	}
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	uri := c.Request.RequestURI
	method := c.Request.Method
	log.Printf("请求:%s \t%s", method, uri)

	contentType := c.ContentType()
	if binding.MIMEJSON == contentType {
		paramMap := make(map[string]interface{})
		c.ShouldBindBodyWith(&paramMap, binding.JSON)
		paramJson, err := json.Marshal(paramMap)
		if err == nil {
			log.Printf("请求参数:%s", paramJson)
		}

	}

	c.Next()
}
