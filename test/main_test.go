package test

import (
	"fmt"
	"github.com/MrKrisYu/koi-go-common/config/source/file"
	"github.com/MrKrisYu/koi-go-common/sdk/api"
	"github.com/MrKrisYu/koi-go-common/sdk/config"
	"github.com/MrKrisYu/koi-go-common/sdk/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"testing"
)

type JsonRequest struct {
	Value string `json:"value"`
}

func TestGinLogger(t *testing.T) {
	config.Setup(file.NewSource(file.WithPath("./application.yaml")))

	engine := gin.Default()
	engine.Use(middleware.GinLogger("X-Request-Id"))

	engine.POST("/testJson", func(c *gin.Context) {
		var req JsonRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, req)
	})

	controller := TestController{}
	engine.GET("/testMixParams", controller.MixParams)

	err := engine.Run("localhost:8080")
	if err != nil {
		fmt.Println("Error running:", err.Error())
	}
	fmt.Println("Server shutdown.....")

	//select {}
}

type MixReq struct {
	QueryParam1  string `form:"param1" binding:"required"`
	QueryParam2  string `form:"param2" binding:"required"`
	QueryParam3  string `form:"param3" binding:"required"`
	HeaderParam1 string `header:"header1" binding:"required"`
}

type TestController struct {
	api.Api
}

func (c *TestController) MixParams(ctx *gin.Context) {
	var req MixReq
	err := c.MakeContext(ctx).
		Bind(ctx, &req, binding.Header, binding.Query)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"data": "RequestError: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": req})
	return
}
