package view

import (
	"encoding/json"
	"fuzz_debug_platform/sqldebug"
	"net/http"
)

func CodePos() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		codePos, _ := sqldebug.Summarize()
		json.NewEncoder(writer).Encode(codePos)
	}
}
