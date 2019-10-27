package sqldebug

import (
	"encoding/json"
	"testing"
)

var infos = []*SQLTraceInfo{
	{
		IsBug: true,
		Sql:   "select * from sbtest;",
		Trace: []SqlTraceBlock{
			{
				File: "server/tokenlimiter.go",
				Line: []SQLTraceBlockLine{
					[]int64{21, 24},
					[]int64{27, 29},
				},
			},
		},
	},
	{
		IsBug: false,
		Sql:   "select * from sbtest2;",
		Trace: []SqlTraceBlock{
			{
				File: "server/buffered_read_conn.go",
				Line: []SQLTraceBlockLine{
					[]int64{21, 21},
				},
			},
		},
	},
	{
		IsBug: true,
		Sql:   "select * from sbtest3;",
		Trace: []SqlTraceBlock{
			{
				File: "server/buffered_read_conn.go",
				Line: []SQLTraceBlockLine{
					[]int64{29, 31},
					[]int64{33, 38},
					[]int64{21, 21},
				},
			},
			{
				File: "server/tokenlimiter.go",
				Line: []SQLTraceBlockLine{
					[]int64{21, 24},
				},
			},
		},
	},
}

func TestTrace(t *testing.T) {
	defer Flush()
	for _, info := range infos {
		if err := Trace(info); err != nil {
			t.Fatalf("unexpect error occured: %s", err)
		}
	}
}

func TestSummarize(t *testing.T) {
	defer Flush()
	for _, info := range infos {
		if err := Trace(info); err != nil {
			t.Fatalf("unexpect error occured: %s", err)
		}
	}
	codePos, err := Summarize()
	if err != nil {
		t.Fatalf("unexpect error occured: %s", err)
	}
	data, err := json.Marshal(codePos)
	if err != nil {
		t.Fatalf("unexpect error occured: %s", err)
	}
	expect := `[{"filePath":"server/tokenlimiter.go","codeBlocks":[{"filePath":"server/tokenlimiter.go","startLine":21,"endLine":24,"score":1,"count":2,"sqls":[]},{"filePath":"server/tokenlimiter.go","startLine":27,"endLine":29,"score":1,"count":1,"sqls":[]}],"content":"// Copyright 2015 PingCAP, Inc.\n//\n// Licensed under the Apache License, Version 2.0 (the \"License\");\n// you may not use this file except in compliance with the License.\n// You may obtain a copy of the License at\n//\n//     http://www.apache.org/licenses/LICENSE-2.0\n//\n// Unless required by applicable law or agreed to in writing, software\n// distributed under the License is distributed on an \"AS IS\" BASIS,\n// See the License for the specific language governing permissions and\n// limitations under the License.\n\npackage server\n\n// Token is used as a permission to keep on running.\ntype Token struct {\n}\n\n// TokenLimiter is used to limit the number of concurrent tasks.\ntype TokenLimiter struct {\n\tcount uint\n\tch    chan *Token\n}\n\n// Put releases the token.\nfunc (tl *TokenLimiter) Put(tk *Token) {\n\ttl.ch \u003c- tk\n}\n\n// Get obtains a token.\nfunc (tl *TokenLimiter) Get() *Token {\n\treturn \u003c-tl.ch\n}\n\n// NewTokenLimiter creates a TokenLimiter with count tokens.\nfunc NewTokenLimiter(count uint) *TokenLimiter {\n\ttl := \u0026TokenLimiter{count: count, ch: make(chan *Token, count)}\n\tfor i := uint(0); i \u003c count; i++ {\n\t\ttl.ch \u003c- \u0026Token{}\n\t}\n\n\treturn tl\n}\n"},{"filePath":"server/buffered_read_conn.go","codeBlocks":[{"filePath":"server/buffered_read_conn.go","startLine":21,"endLine":21,"score":0.5,"count":1,"sqls":[]},{"filePath":"server/buffered_read_conn.go","startLine":29,"endLine":31,"score":1,"count":1,"sqls":[]},{"filePath":"server/buffered_read_conn.go","startLine":33,"endLine":38,"score":1,"count":1,"sqls":[]}],"content":"// Copyright 2017 PingCAP, Inc.\n//\n// Licensed under the Apache License, Version 2.0 (the \"License\");\n// you may not use this file except in compliance with the License.\n// You may obtain a copy of the License at\n//\n//     http://www.apache.org/licenses/LICENSE-2.0\n//\n// Unless required by applicable law or agreed to in writing, software\n// distributed under the License is distributed on an \"AS IS\" BASIS,\n// See the License for the specific language governing permissions and\n// limitations under the License.\n\npackage server\n\nimport (\n\t\"bufio\"\n\t\"net\"\n)\n\nconst defaultReaderSize = 16 * 1024\n\n// bufferedReadConn is a net.Conn compatible structure that reads from bufio.Reader.\ntype bufferedReadConn struct {\n\tnet.Conn\n\trb *bufio.Reader\n}\n\nfunc (conn bufferedReadConn) Read(b []byte) (n int, err error) {\n\treturn conn.rb.Read(b)\n}\n\nfunc newBufferedReadConn(conn net.Conn) *bufferedReadConn {\n\treturn \u0026bufferedReadConn{\n\t\tConn: conn,\n\t\trb:   bufio.NewReaderSize(conn, defaultReaderSize),\n\t}\n}\n"}]`
	if string(data) != expect {
		t.Fatalf("expect %s, got %s", expect, string(data))
	}
}
