package view

import (
	"encoding/json"
	"fuzz_debug_platform/sqldebug"
	"net/http"
)

func CodePos() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		codePos, _ := sqldebug.Summarize()
		json.NewEncoder(writer).Encode(codePos)
	}
}
