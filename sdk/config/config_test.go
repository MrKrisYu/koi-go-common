package config

import (
	"fmt"
	"github.com/MrKrisYu/koi-go-common/config/source/file"
	"github.com/MrKrisYu/koi-go-common/sdk"
	"testing"
)

func TestConfig(t *testing.T) {
	configPath := "./application-example.yaml"
	Setup(file.NewSource(file.WithPath(configPath)))

	fmt.Println(sdk.RuntimeContext)
	s := _cfg.String()
	fmt.Println(s)
}
