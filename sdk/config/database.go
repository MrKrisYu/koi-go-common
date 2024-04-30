package config

type Database struct {
	Driver          string             `json:"driver"`
	Source          string             `json:"source"`
	ConnMaxIdleTime int                `json:"connMaxIdleTime"`
	ConnMaxLifeTime int                `json:"connMaxLifeTime"`
	MaxIdleConns    int                `json:"maxIdleConns"`
	MaxOpenConns    int                `json:"maxOpenConns"`
	Registers       []DBResolverConfig `json:"registers"`
}

type DBResolverConfig struct {
	Sources  []string `json:"sources"`
	Replicas []string `json:"replicas"`
	Policy   string   `json:"policy"`
	Tables   []string `json:"tables"`
}

var (
	DatabaseConfig  = new(Database)
	DatabasesConfig = make(map[string]*Database)
)
