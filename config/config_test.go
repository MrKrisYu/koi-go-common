package config

import (
	"encoding/json"
	"fmt"
	"github.com/MrKrisYu/koi-go-common/config/source/file"
	"testing"
)

func TestConfig(t *testing.T) {
	c, err := NewConfig()
	if err != nil {
		t.Error(err)
	}
	err = c.Load(file.NewSource(file.WithPath("./application.yaml")))
	if err != nil {
		t.Error(err)
	}
	m := c.Map()
	jsonbtyes, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("DefaultConfig = %s\n", jsonbtyes)
	fmt.Printf("diy params = %s\n", c.Get("settings", "diy", "servicePath").String(""))
	fmt.Printf("diy params = %s \n", c.Get("settings", "diy", "moduleName").StringSlice([]string{}))
}
