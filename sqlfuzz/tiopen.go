package sqlfuzz

import (
	"fuzz_debug_platform/config"
	"github.com/pkg/errors"
	"net/http"
)

func ToggleTiDB() error {
	resp, err := http.Get(config.GetGlobalConf().TiDBWrapperSwitchAddr)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("http resp status error")
	}

	return nil
}
