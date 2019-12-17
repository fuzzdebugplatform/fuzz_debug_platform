# SqlDebug

## Toy Arount with It by Docker Compose

You should install [Docker Compose](https://docs.docker.com/compose/install/)
 on your computer first.

```bash
$ git clone https://github.com/fuzzdebugplatform/docker-compose.git
$ cd docker-compose
$ docker-compose up -d
Creating docker-compose_pltform_1 ... done
Creating docker-compose_tw_1      ... done
Creating docker-compose_mysql_1   ... done
```

The test above is about window functions, with yy 
file [exampleyy/windows.yy](exampleyy/windows.yy).

You can use `docker logs` to see test process:

```bash
$ docker logs docker-compose_pltform_1
2019/12/17 11:49:35 staring prepare data in db
2019/12/17 11:49:53 prepare data in db ok
2019/12/17 11:49:53 tidb statistic funtion open ok
2019/12/17 11:49:53 starting generate query
2019/12/17 11:50:04 fuzz ok
```

When you see "fuzz ok", you can open `localhost:43000`
to see analysis（you can alse open this page advance, but you'd
 bettor flush the page after fuzz ok）:
 
![quick experience](img/quick.gif)




## Run with Binary

 - First, you should wrap the tidb version you are interesting with  [tidb-wrapper](https://github.com/fuzzdebugplatform/tidb-wrapper)

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

 - build web page

```bash
cd web
yarn run build
```

 - run test
 
```bash
./bin/platform -Y exampleyy/subquery_test.yy --dsn1 "root:@tcp(127.0.0.1:4000)/randgen"  --dsn2 "root:@tcp(127.0.0.1:4001)/randgen" -Q 100 --debug -W "web/build"
```

visit `localhost:3000` to see analysis.