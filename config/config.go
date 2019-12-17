package config

import "sync/atomic"

type Config struct {
	TiDBSourceDir        string  `toml:"tidb-source-dir", json:"tidb-source-dir"`
	TiDBWrapperSwitchAddr string `toml:"tidb-wrapper-switch-addr", json:"tidb-wrapper-switch-addr"`
	TiDBTraceServerAddr  string  `toml:"tidb-trace-server-addr", json:"tidb-trace-server-addr"`
	TopN                 int     `toml:"top-n", json: "top-n"`
	FailureRateThreshold float64 `toml:"failure-rate-threshold", json:"failure-rate-threshold"`
}

var defaultConf = Config{
	TiDBSourceDir:        "/home/dqyuan/append/language/Go/projects/hackthon/tidb-bad/",
	TiDBWrapperSwitchAddr: "http://localhost:43222/switch",
	TiDBTraceServerAddr:  "http://localhost:43222/trace/",
	TopN:                 10000,
	FailureRateThreshold: 0.7,
}

var globalConf = atomic.Value{}

func GetGlobalConf() *Config {
	return globalConf.Load().(*Config)
}

func init() {
	globalConf.Store(&defaultConf)
}
