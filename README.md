# TiSqlDebug


 - config tidb source dir in `config/config.go`(must consistent with wrapped tidb)
 
```go
var defaultConf = Config{
	TiDBSourceDir:        "/home/dqyuan/append/language/Go/projects/hackthon/tidb-bad/",
	TiDBTraceServerAddr:  "http://localhost:43222/trace/",
	TopN:                 10000,
	FailureRateThreshold: 0.7,
}
``` 

 - start a wrapped tidb
 
```bash
tidb-bad-wrapper$ ./bin/tidb-server -path testdb
```

 - start a normal tidb (another version)
 
```bash
tidb-good$ ./bin/tidb-server -P 4001 -path testdb
```

 - run test
 
```bash
./bin/platform -Y exampleyy/subquery_test.yy --dsn1 "root:@tcp(127.0.0.1:4000)/randgen"  --dsn2 "root:@tcp(127.0.0.1:4001)/randgen" -Q 100 --debug
```

 - run web

```bash
cd web
yarn start
```

visit `localhost:3000`