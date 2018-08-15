package config

type Config struct {
	DbSource          DataSource `json:"data_source,omitempty"`
	LevelDbPath       string     `json:"level_db_path,omitempty"`
	ServerCert        string     `json:"server_cert,omitempty"`
	ServerKey         string     `json:"server_key,omitempty"`
	ClientCert        string     `json:"client_cert,omitempty"`
	GrpcSerPort       string     `json:"grpc_ser_port,omitempty"`
	EtcdEndpoints     []string   `json:"etcd_endpoints,omitempty"`
	VoucherName       string     `json:"voucher_name,omitempty"`
	ManagerIpPort     string     `json:"manager_ip_port,omitempty"`
	DepositUrl        string     `json:"deposit_url,omitempty"`
	WithDrawUrl       string     `json:"withdraw_url,omitempty"`
	WithDrawTxUrl     string     `json:"withdraw_tx_url,omitempty"`
	TokenChangeUrl    string     `json:"token_change_url,omitempty"`
	RegistApprovalUrl string     `json:"regist_approval_url,omitempty"`
	CompanionServer   string     `json:"companion_server,omitempty"`
	VoucherServer     string     `json:"voucher_server,omitempty"`
	Assets            string     `json:"assets_url,omitempty"`
	TradeHistory      string     `json:"trade_history_url,omitempty"`
}

type DataSource struct {
	DriverName string `json:"driver_name,omitempty"`
	Url        string `json:"url,omitempty"`
	MaxIdle    int    `json:"max_idle,omitempty"`
	MaxConn    int    `json:"max_conn,omitempty"`
	AliasName  string `json:"alias_name,omitempty"`
	Debug      bool   `json:"debug,omitempty"`
}

var GConfig *Config
