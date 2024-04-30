package config

type Jwt struct {
	Secret  string `json:"secret"`
	Timeout int64  `json:"timeout"`
}

var JwtConfig = new(Jwt)
