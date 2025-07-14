package test

import (
	"fmt"
	"github.com/MrKrisYu/koi-go-common/config/source/file"
	"github.com/MrKrisYu/koi-go-common/sdk"
	"github.com/MrKrisYu/koi-go-common/sdk/api"
	"github.com/MrKrisYu/koi-go-common/sdk/config"
	"github.com/MrKrisYu/koi-go-common/sdk/i18n/example"
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

	//translator := example.NewDefaultTranslator(example.DefaultLanguage, example.AllowedLanguage)

	engine := gin.Default()
	engine.Use(middleware.RequestId("X-Request-Id")).
		Use(middleware.GinLogger("X-Request-Id")).
		Use(example.AcceptLanguage())

	engine.POST("/testJson", func(c *gin.Context) {
		var req JsonRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var testData struct {
			StrValue     string  `json:"strValue"`
			IntValue     int     `json:"intValue"`
			BooleanValue bool    `json:"booleanValue"`
			FloatValue   float64 `json:"floatValue"`
			StructValue  struct {
				StrValue     string `json:"strValue"`
				IntValue     int    `json:"intValue"`
				BooleanValue bool   `json:"booleanValue"`
			}
			StrArrayValue []string          `json:"strArrayValue"`
			MapValue      map[string]string `json:"mapValue"`
			StrPointValue *string           `json:"strPointValue"`
		}
		//lang := language.Und
		//if value, exists := c.Get(header.AcceptLanguageFlag); exists {
		//	lang = value.(language.Tag)
		//}
		//message := translator.TrWithData(lang, i18n.Message{ID: req.Value, DefaultMessage: "DefaultMessage", Args: map[string]interface{}{"Arg": "测试参数"}})
		//req.Value = message
		c.JSON(http.StatusOK, testData)
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

type Menu struct {
	Name     string
	Path     string
	Children []Menu
}

type TestStruct struct {
	StrValue     string  `json:"strValue"`
	IntValue     int     `json:"intValue"`
	BooleanValue bool    `json:"booleanValue"`
	FloatValue   float64 `json:"floatValue"`
	StructValue  struct {
		StrValue      string            `json:"strValue"`
		IntValue      int               `json:"intValue"`
		BooleanValue  bool              `json:"booleanValue"`
		StrArrayValue []string          `json:"strArrayValue"`
		MapValue      map[string]string `json:"mapValue"`
		StrPointValue *string           `json:"strPointValue"`
	}
	StrArrayValue    []string          `json:"strArrayValue"`
	MapValue         map[string]string `json:"mapValue"`
	StrPointValue    *string           `json:"strPointValue"`
	EmbedArrayStruct []TestStruct      `json:"embedArrayStruct"`
}

func (c *TestController) MixParams(ctx *gin.Context) {
	//var req MixReq
	//var req JsonRequest
	_ = c.MakeContext(ctx)
	//myError := i18n.NewMyErrorWithMessageID(i18n.MessageID{ID: "response.OK", DefaultMessage: "OKOK"})
	//myError.AddMessage(i18n.MessageID{ID: "response.NoOK", DefaultMessage: "OKOK1"})
	//myError.AddMessage(i18n.MessageID{ID: "response.NOK", DefaultMessage: "OKOK2"})
	//c.Error(ctx, response.OK, myError)
	//err := i18n.NewMyError(fmt.Errorf("我就想看看你能不能直接返回这个错误"))
	var testData TestStruct

	//testData := Menu{
	//	Name:     "",
	//	Path:     "",
	//	Children: make([]Menu, 0),
	//}

	c.OK(ctx, testData)
	return
}

func TestConfig(t *testing.T) {
	config.Setup(file.NewSource(file.WithPath("./application-test.yaml")))
	config := sdk.RuntimeContext.GetDefaultConfig()
	fmt.Printf("moduleName = %v \n", config.Get("settings", "moduleName").StringSlice([]string{}))

}
