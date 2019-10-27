package sqldebug

import "github.com/pingcap/errors"

// Notify receives digestCode and isBug information
func Notify(digestCode string, isBug bool) error {
	traceInfo, err := Collect(digestCode, isBug)
	if err != nil {
		return errors.Trace(err)
	}
	if err := Trace(traceInfo); err != nil {
		return errors.Trace(err)
	}
	return nil
}
