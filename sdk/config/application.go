package config

type Application struct {
	ReadTimeout  int    `json:"readTimeout"`  // Gin读数据超时时间，单位秒
	WriteTimeout int    `json:"writeTimeout"` // Gin写数据超时时间，单位秒
	Host         string `json:"host"`         // 监听主机号
	Port         int64  `json:"port"`         // 监听端口号
	Name         string `json:"name"`         // 应用名称
	Mode         string `json:"mode"`         // 开发模式： dev-开发环境， prod- 生成环境， test-测试环境
	DemoMsg      string `json:"demoMsg"`      // 启动的提示消息
}

var ApplicationConfig = new(Application)

func getDefaultAppConfig() *Application {
	return &Application{
		ReadTimeout:  600,
		WriteTimeout: 600,
		Host:         "localhost",
		Port:         80,
		Name:         "Koi-ApplicationConfig",
		Mode:         "dev",
		DemoMsg:      "Koi-ApplicationConfig started.......",
	}
}

func (a *Application) Key() string {
	return "application"
}

func (a *Application) Setup() {
	defaultAppConfig := getDefaultAppConfig()
	if a.ReadTimeout <= 0 {
		a.ReadTimeout = defaultAppConfig.ReadTimeout
	}
	if a.WriteTimeout <= 0 {
		a.WriteTimeout = defaultAppConfig.WriteTimeout
	}
	if len(a.Host) == 0 {
		a.Host = defaultAppConfig.Host
	}
	if a.Port <= 0 {
		a.Port = defaultAppConfig.Port
	}
	if len(a.Mode) == 0 {
		a.Mode = defaultAppConfig.Mode
	}
}
