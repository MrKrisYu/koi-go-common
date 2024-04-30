package config

type Ssl struct {
	KeyStr string `json:"key"`
	Pem    string `json:"pem"`
	Enable bool   `json:"enable"`
	Domain string `json:"domain"`
}

var SslConfig = new(Ssl)
