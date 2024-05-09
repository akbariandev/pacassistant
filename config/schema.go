package config

type Config[T any] struct {
	Address string   `yaml:"address" json:"address"`
	Domain  string   `yaml:"domain" json:"domain"`
	Origins []string `yaml:"origins" json:"origins"`
	Grpc    struct {
		Port            int    `yaml:"port" json:"port"`
		CertFilePath    string `yaml:"cert_file_path" json:"cert_file_path"`
		CertKeyFilePath string `yaml:"cert_key_file_path" json:"cert_key_file_path"`
	} `yaml:"grpc" json:"grpc"`
	Rest struct {
		Port            int    `yaml:"port" json:"port"`
		CertFilePath    string `yaml:"cert_file_path" json:"cert_file_path"`
		CertKeyFilePath string `yaml:"cert_key_file_path" json:"cert_key_file_path"`
	} `yaml:"rest" json:"rest"`
	ThirdPartyAPI struct {
		Port int `yaml:"port" json:"port"`
	} `yaml:"third_party_api" json:"third_party_api"`
	Websocket struct {
		Port            int    `yaml:"port" json:"port"`
		CertFilePath    string `yaml:"cert_file_path" json:"cert_file_path"`
		CertKeyFilePath string `yaml:"cert_key_file_path" json:"cert_key_file_path"`
	} `yaml:"websocket" json:"websocket"`
	GraphQL struct {
		Port int `yaml:"port" json:"port"`
	} `yaml:"graphql" json:"graphql"`
	Development bool          `yaml:"development" json:"development"`
	GrpcClients []*GrpcClient `yaml:"grpc_clients" json:"grpc_clients"`
	Database    *Database     `yaml:"database" json:"database"`
	Logging     *Logging      `yaml:"logging" json:"logging"`
	Telegram    *Telegram     `yaml:"telegram" json:"telegram"`
	ExtraData   T             `yaml:"extra_data" json:"extra_data"`
}

type GrpcClient struct {
	Name           string `yaml:"name" json:"name"`
	Address        string `yaml:"address" json:"address"`
	Port           int    `yaml:"port" json:"port"`
	SocketFilePath string `yaml:"socket_file_path" json:"socket_file_path"`
	CertCAFilePath string `yaml:"cert_ca_file_path" json:"cert_ca_file_path"`
}

type Database struct {
	Mongodb struct {
		URI          string `yaml:"uri" json:"uri"`
		DatabaseName string `yaml:"database_name" json:"database_name"`
	} `yaml:"mongodb" json:"mongodb"` // Mongodb URI address
	MySQL struct {
		URI          string `yaml:"uri" json:"uri"`
		DatabaseName string `yaml:"database_name" json:"database_name"`
	} `yaml:"mysql" json:"mysql"` // MySQL URI address
	Postgres struct {
		URI          string `yaml:"uri" json:"uri"`
		DatabaseName string `yaml:"database_name" json:"database_name"`
	} `yaml:"postgres" json:"postgres"` // Postgres URI address
	Redis struct {
		Address  string `yaml:"address" json:"address"`
		Username string `yaml:"username" json:"username"`
		Password string `yaml:"password" json:"password"`
		Database int    `yaml:"database" json:"database"`
	} `yaml:"redis" json:"redis"` // Redis URI address
}

type Logging struct {
	Debug        bool   `yaml:"debug" json:"debug"`
	Handler      uint8  `yaml:"handler" json:"handler"` // Handler 0= console handler, 1= text handler, 2= json handler
	EnableCaller bool   `yaml:"enable_caller" json:"enable_caller"`
	SentryDSN    string `yaml:"sentry_dsn" json:"sentry_dsn"`
}

type Telegram struct {
	ApiKey string `yaml:"api_key" json:"apiKey"`
}
