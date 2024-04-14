package entity

import "encoding/json"

type Settings struct {
	// entities
	Root config `json:"settings"`
}

type config struct {
	ApplicationConfig *Application `json:"application"`
	LogConfig         *Log         `json:"log"`
}

func NewSettings() *Settings {
	return &Settings{
		Root: config{
			ApplicationConfig: ApplicationConfig,
			LogConfig:         LogConfig,
		},
	}
}

func (s *Settings) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "[error] Json Marshaling error"
	}
	return string(bytes)
}

func (s *Settings) Init() {
	s.Root.ApplicationConfig.Setup()
	s.Root.LogConfig.Setup()
}
