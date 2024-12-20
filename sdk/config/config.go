package config

import (
	"encoding/json"
	"fmt"
	"github.com/MrKrisYu/koi-go-common/config"
	"github.com/MrKrisYu/koi-go-common/config/source"
	"github.com/MrKrisYu/koi-go-common/sdk"
	"log"
)

var (
	ExtendConfig interface{}
	_cfg         *Settings
)

type Config struct {
	ApplicationConfig *Application          `json:"application"`
	LogConfig         *Log                  `json:"log"`
	Ssl               *Ssl                  `json:"ssl"`
	Jwt               *Jwt                  `json:"jwt"`
	Database          *Database             `json:"database"`
	Databases         *map[string]*Database `json:"databases"`
	Cache             *Cache                `json:"cache"`
	Queue             *Queue                `json:"queue"`
	Locker            *Locker               `json:"locker"`
	Extend            interface{}           `json:"extend"`
}

// 多db改造
func (e *Config) multiDatabase() {
	if len(*e.Databases) == 0 {
		*e.Databases = map[string]*Database{
			"*": e.Database,
		}

	}
}

// Setup 载入配置文件
func Setup(s source.Source,
	fs ...func()) {
	_cfg = &Settings{
		Root: Config{
			ApplicationConfig: ApplicationConfig,
			Ssl:               SslConfig,
			LogConfig:         LogConfig,
			Jwt:               JwtConfig,
			Database:          DatabaseConfig,
			Databases:         &DatabasesConfig,
			Cache:             CacheConfig,
			Queue:             QueueConfig,
			Locker:            LockerConfig,
			Extend:            ExtendConfig,
		},
		callbacks: fs,
	}
	var err error
	runtimeConfig, err := config.NewConfig(
		config.WithSource(s),
		config.WithEntity(_cfg),
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("New config object fail: %s", err.Error()))
	}
	sdk.RuntimeContext.SetDefaultConfig(runtimeConfig)

	_cfg.Init()
}

type Settings struct {
	// entities
	Root      Config `json:"settings"`
	callbacks []func()
}

func NewSettings() *Settings {
	return &Settings{
		Root: Config{
			ApplicationConfig: ApplicationConfig,
			LogConfig:         LogConfig,
		},
	}
}

func (s *Settings) runCallback() {
	for i := range s.callbacks {
		s.callbacks[i]()
	}
}

func (s *Settings) OnChange() {
	s.init()
	//log.Println("config change and reload")
}

func (s *Settings) Init() {
	s.init()
	//log.Println("config init")
}

func (s *Settings) init() {
	s.Root.ApplicationConfig.Setup()
	s.Root.multiDatabase()
	s.Root.LogConfig.Setup()
	s.runCallback()
}

func (s *Settings) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "[error] Json Marshaling error"
	}
	return string(bytes)
}
