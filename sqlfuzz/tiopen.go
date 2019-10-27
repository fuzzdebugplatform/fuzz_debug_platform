package sqlfuzz

import (
	"github.com/pkg/errors"
	"net/http"
)

func ToggleTiDB() error {
	resp, err := http.Get("http://localhost:43222/switch")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("http resp status error")
	}

	return nil
}
