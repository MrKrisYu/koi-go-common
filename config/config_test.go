package config

import (
	"encoding/json"
	"fmt"
	"github.com/MrKrisYu/koi-go-common/config/encoder/yaml"
	"github.com/MrKrisYu/koi-go-common/config/entity"
	"github.com/MrKrisYu/koi-go-common/sdk"
	"testing"
)

func TestConfig(t *testing.T) {
	Setup(
		WithPath("./application.yaml"),
		WithEncoder(yaml.NewEncoder()),
	)
	applicationConfig, err := GetEntity[*entity.Application](Global, entity.ApplicationConfig.Key())
	if err != nil {
		t.Error(err)
	}
	bytes, err := json.Marshal(applicationConfig)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("sdkRuntimeContext = ", sdk.RuntimeContext)
	fmt.Println("GetEntityFromApplicationKey = ", string(bytes))
	fmt.Println("configuration = ", Global.Settings.String())

}
