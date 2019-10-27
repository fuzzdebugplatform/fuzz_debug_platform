package sqldebug

import (
	"encoding/json"
	"fuzz_debug_platform/config"
	"github.com/pingcap/errors"
	"net/http"
)

type SQLTraceBlockLine []int64

type SqlTraceBlock struct {
	File string              `json:"file"`
	Line []SQLTraceBlockLine `json:"line"`
}

type SQLTraceInfo struct {
	IsBug bool
	Sql   string          `json:"sql"`
	Trace []SqlTraceBlock `json:"trace"`
}

// Collect fetches sqlTraceInfo from TiDB server
func Collect(digestCode string, isBug bool) (*SQLTraceInfo, error) {
	resp, err := http.Get(config.GetGlobalConf().TiDBTraceServerAddr + digestCode)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var traceInfo SQLTraceInfo

	if err := decoder.Decode(&traceInfo); err != nil {
		return nil, errors.Trace(err)
	}

	traceInfo.IsBug = isBug
	return &traceInfo, nil
}
