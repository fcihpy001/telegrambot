package utils

type ConfigData struct {
	Token       string `yaml:"token"`
	WebhookUrl  string `yaml:"webhook_url"`
	CertFile    string `yaml:"cert_file"`
	KeyFile     string `yaml:"key_file"`
	ApiPath     string `yaml:"api_path"`
	URL         string `yaml:"url"`
	DatabaseURL string `yaml:"database_url"`
	RedisURL    string `yaml:"redis_url"`
	Debug       bool   `yaml:"debug"`
	Mysql       Mysql  `yaml:"mysql"`
	Dsn         string `yaml:"dsn"`
}

type Mysql struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	UserName string `yaml:"user_name"`
	Passwd   string `yaml:"passwd"`
}

type RequestFile struct {
	FileName string
}
