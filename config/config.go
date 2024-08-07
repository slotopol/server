package config

import (
	"time"
)

var (
	// compiled binary version, sets by compiler with command
	//    go build -ldflags="-X 'github.com/slotopol/server/config.BuildVers=%buildvers%'"
	BuildVers string

	// compiled binary build date, sets by compiler with command
	//    go build -ldflags="-X 'github.com/slotopol/server/config.BuildTime=%buildtime%'"
	BuildTime string
)

type CfgJwtAuth struct {
	AccessTTL    time.Duration `json:"access-ttl" yaml:"access-ttl" mapstructure:"access-ttl"`
	RefreshTTL   time.Duration `json:"refresh-ttl" yaml:"refresh-ttl" mapstructure:"refresh-ttl"`
	AccessKey    string        `json:"access-key" yaml:"access-key" mapstructure:"access-key"`
	RefreshKey   string        `json:"refresh-key" yaml:"refresh-key" mapstructure:"refresh-key"`
	NonceTimeout time.Duration `json:"nonce-timeout" yaml:"nonce-timeout" mapstructure:"nonce-timeout"`
}

// CfgWebServ is web server settings.
type CfgWebServ struct {
	// List of network origins (IPv4 addresses, IPv4 CIDRs, IPv6 addresses or IPv6 CIDRs) from which to trust request's headers that contain alternative client IP when `(*gin.Engine).ForwardedByClientIP` is `true`.
	TrustedProxies []string `json:"trusted-proxies" yaml:"trusted-proxies" mapstructure:"trusted-proxies"`
	// List of address:port values for non-encrypted connections. Address is skipped in most common cases, port only remains.
	PortHTTP []string `json:"port-http" yaml:"port-http" mapstructure:"port-http"`
	// Maximum duration for reading the entire request, including the body.
	ReadTimeout time.Duration `json:"read-timeout" yaml:"read-timeout" mapstructure:"read-timeout"`
	// Amount of time allowed to read request headers.
	ReadHeaderTimeout time.Duration `json:"read-header-timeout" yaml:"read-header-timeout" mapstructure:"read-header-timeout"`
	// Maximum duration before timing out writes of the response.
	WriteTimeout time.Duration `json:"write-timeout" yaml:"write-timeout" mapstructure:"write-timeout"`
	// Maximum amount of time to wait for the next request when keep-alives are enabled.
	IdleTimeout time.Duration `json:"idle-timeout" yaml:"idle-timeout" mapstructure:"idle-timeout"`
	// Controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line, in bytes.
	MaxHeaderBytes int `json:"max-header-bytes" yaml:"max-header-bytes" mapstructure:"max-header-bytes"`
	// Maximum duration to wait for graceful shutdown.
	ShutdownTimeout time.Duration `json:"shutdown-timeout" yaml:"shutdown-timeout" mapstructure:"shutdown-timeout"`
}

type CfgXormDrv struct {
	// Provides driver name to create XORM engine.
	// It can be "sqlite3" or "mysql".
	DriverName string `json:"driver-name" yaml:"driver-name" mapstructure:"driver-name"`
	// Data source name for 'club' database to create XORM engine.
	// For sqlite3 it should be database file name (slot-club.sqlite),
	// for mysql it should match to pattern user:password@/slot_club.
	ClubSourceName string `json:"club-source-name" yaml:"club-source-name" mapstructure:"club-source-name"`
	// Data source name for 'spin' database to create XORM engine.
	// For sqlite3 it should be database file name (slot-spin.sqlite),
	// for mysql it should match to pattern user:password@/slot_spin.
	SpinSourceName string `json:"spin-source-name" yaml:"spin-source-name" mapstructure:"spin-source-name"`
	// Duration between flushes of SQL batching buffers.
	SqlFlushTick time.Duration `json:"sql-flush-tick" yaml:"sql-flush-tick" mapstructure:"sql-flush-tick"`
	// Maximum number of users at each club to group wallets updates to database.
	// If size is 1, update will be sequential with error code expecting.
	BankBufferSize int `json:"bank-buffer-size" yaml:"bank-buffer-size" mapstructure:"bank-buffer-size"`
	// Maximum size of buffer for spinlog items to group inserting to database.
	SpinlogBufferSize int `json:"spinlog-buffer-size" yaml:"spinlog-buffer-size" mapstructure:"spinlog-buffer-size"`
}

type CfgGameplay struct {
	// Maximum value to add to wallet by one transaction.
	AdjunctLimit float64 `json:"adjunct-limit" yaml:"adjunct-limit" mapstructure:"adjunct-limit"`
	// Maximum number of spin attempts at bad bank balance.
	MaxSpinAttempts int `json:"max-spin-attempts" yaml:"max-spin-attempts" mapstructure:"max-spin-attempts"`
}

// Config is common service settings.
type Config struct {
	CfgJwtAuth  `json:"authentication" yaml:"authentication" mapstructure:"authentication"`
	CfgWebServ  `json:"web-server" yaml:"web-server" mapstructure:"web-server"`
	CfgXormDrv  `json:"database" yaml:"database" mapstructure:"database"`
	CfgGameplay `json:"gameplay" yaml:"xorm" mapstructure:"gameplay"`
}

// Instance of common service settings.
// Inits default values if config is not found.
var Cfg = &Config{
	CfgJwtAuth: CfgJwtAuth{
		AccessTTL:    1 * 24 * time.Hour,
		RefreshTTL:   3 * 24 * time.Hour,
		AccessKey:    "skJgM4NsbP3fs4k7vh0gfdkgGl8dJTszdLxZ1sQ9ksFnxbgvw2RsGH8xxddUV479",
		RefreshKey:   "zxK4dUnuq3Lhd1Gzhpr3usI5lAzgvy2t3fmxld2spzz7a5nfv0hsksm9cheyutie",
		NonceTimeout: 150 * time.Second,
	},
	CfgWebServ: CfgWebServ{
		TrustedProxies:    []string{"127.0.0.0/8"},
		PortHTTP:          []string{":8080"},
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ShutdownTimeout:   15 * time.Second,
	},
	CfgXormDrv: CfgXormDrv{
		DriverName:        "sqlite3",
		ClubSourceName:    "slot-club.sqlite",
		SpinSourceName:    "slot-spin.sqlite",
		SqlFlushTick:      2500 * time.Millisecond,
		BankBufferSize:    25,
		SpinlogBufferSize: 50,
	},
	CfgGameplay: CfgGameplay{
		AdjunctLimit:    100000,
		MaxSpinAttempts: 300,
	},
}
