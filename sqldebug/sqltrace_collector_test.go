package sqldebug

import (
	"encoding/json"
	config2 "fuzz_debug_platform/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	data, err := ioutil.ReadFile("testcase/show_database.json")
	if err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
	var sql SQLTraceInfo

	if err := json.Unmarshal(data, &sql); err != nil {
		t.Fatalf("unexpect error: %s", err)
	}
}

func TestCollect(t *testing.T) {
	data, _ := ioutil.ReadFile("testcase/show_database.json")
	var expect SQLTraceInfo
	if err := json.Unmarshal(data, &expect); err != nil {
		t.Fatalf("unexpect error occured: %s", err)
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
	}))
	defer func() { testServer.Close() }()

	config := config2.GetGlobalConf()
	config.TiDBTraceServerAddr = testServer.URL + "/"

	got, err := Collect("123", false)
	if err != nil {
		t.Fatalf("unexpect error occured: %s", err)
	}

	if !reflect.DeepEqual(*got, expect) {
		gotStr, _ := json.Marshal(got)
		expectStr, _ := json.Marshal(expect)
		t.Fatalf("expect %s, got %s", gotStr, expectStr)
	}
}
