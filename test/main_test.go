package test

import (
	"fmt"
	"github.com/MrKrisYu/koi-go-common/config/source/file"
	"github.com/MrKrisYu/koi-go-common/sdk/config"
	"github.com/MrKrisYu/koi-go-common/sdk/middleware"
	"github.com/gin-gonic/gin"
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

	err := engine.Run("localhost:8080")
	if err != nil {
		fmt.Println("Error running:", err.Error())
	}
	fmt.Println("Server shutdown.....")

	//select {}
}
