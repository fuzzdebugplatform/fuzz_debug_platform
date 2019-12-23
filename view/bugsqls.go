package view

import (
	"encoding/json"
	"fuzz_debug_platform/sqldebug"
	"net/http"
	"net/url"
	"strconv"
)

func BugSqls() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var fPath string
		var startLine, endLine string
		var ok bool
		params := request.URL.Query()
		if fPath, ok = getSingleParam(params, "filepath"); !ok {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if startLine, ok = getSingleParam(params, "startLine"); !ok {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if endLine, ok = getSingleParam(params, "endLine"); !ok {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		startLineNum, err := strconv.ParseInt(startLine, 10, 64)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		endLineNum, err := strconv.ParseInt(endLine, 10, 64)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		sqls := sqldebug.BugSqls(fPath, startLineNum, endLineNum)
		if sqls == nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(writer).Encode(sqls)
	}
}

func getSingleParam(params url.Values, key string) (string, bool) {
	if values, ok := params[key]; ok {
		if len(values) < 1 {
			return "", false
		}

		return values[0], true
	}

	return "", false
}
