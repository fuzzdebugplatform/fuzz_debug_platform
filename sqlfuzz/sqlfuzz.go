package sqlfuzz

import (
	"encoding/hex"
	"fuzz_debug_platform/sqldebug"
	"github.com/pingcap/go-randgen/compare"
	"github.com/pingcap/go-randgen/gendata"
	"github.com/pingcap/go-randgen/grammar"
	"github.com/pingcap/go-randgen/grammar/sql_generator"
	"github.com/spaolacci/murmur3"
	"log"
	"sync"
)

// Succ Fail Counter
type SFCounter struct {
	Succ int
	Fail int
}

type AlterCounter struct {
	// production number,alter number ->counter
	AlterMap map[[2]int]*SFCounter
	Mu       sync.RWMutex
}

type ProductionCounter struct {
	// production number -> counter
	ProductionMap map[int]*SFCounter
	Mu sync.RWMutex
}

var ACounter = &AlterCounter{AlterMap:make(map[[2]int]*SFCounter)}
var PCounter = &ProductionCounter{ProductionMap:make(map[int]*SFCounter)}

func Fuzz(yy string, dsn1 string, dsn2 string, queries int, debug bool) {

	db1, err := compare.OpenDBWithRetry("mysql", dsn1)
	if err != nil {
		log.Fatalf("%s connect error %v\n", dsn1, err)
	}

	db2, err := compare.OpenDBWithRetry("mysql", dsn2)
	if err != nil {
		log.Fatalf("%s connect error %v\n", dsn2, err)
	}

	log.Println("staring prepare data in db")

	ddls, keyf, err := gendata.ByZz("")
	if err != nil {
		log.Fatalf("get ddl error %v\n", err)
	}

	errSql, err := compare.ExecSqlsInDbs(ddls, db1, db2)
	if err != nil {
		log.Fatalf("sql %s err %v\n", errSql, err)
	}

	log.Println("prepare data in db ok")

	// open the tidb statistical function
	if debug {
		err = ToggleTiDB()
		if err != nil {
			log.Fatalf("open tidb statistic fail, %v\n", err)
		}
		log.Println("tidb statistic funtion open ok")
	}

	log.Println("starting generate query")

	sqlIter, err := grammar.NewIter(yy, "query", 500, keyf, false)
	if err != nil {
		log.Fatalf("get iterator error %v\n", err)
	}

	err = sqlIter.Visit(sql_generator.FixedTimesVisitor(func(_ int, sql string) {
		consistent, _, _ := compare.BySql(sql, db1, db2, false)
		if debug {
			sqldebug.Notify(digest(sql), !consistent)
		}
		if consistent {
			PCounter.Mu.Lock()
			for _, prod := range sqlIter.PathInfo().ProductionSet.Productions {
				if c := PCounter.ProductionMap[prod.Number]; c == nil {
					PCounter.ProductionMap[prod.Number] = &SFCounter{}
				} else {
					c.Succ++
				}
			}
			PCounter.Mu.Unlock()

			ACounter.Mu.Lock()
			for _, seq := range sqlIter.PathInfo().SeqSet.Seqs {
				if c := ACounter.AlterMap[[2]int{seq.PNumber, seq.SNumber}]; c == nil {
					ACounter.AlterMap[[2]int{seq.PNumber, seq.SNumber}] = &SFCounter{}
				} else {
					c.Succ++
				}
			}
			ACounter.Mu.Unlock()
		} else {
			PCounter.Mu.Lock()
			for _, prod := range sqlIter.PathInfo().ProductionSet.Productions {
				if c := PCounter.ProductionMap[prod.Number]; c == nil {
					PCounter.ProductionMap[prod.Number] = &SFCounter{}
				} else {
					c.Fail++
				}
			}
			PCounter.Mu.Unlock()

			ACounter.Mu.Lock()
			for _, seq := range sqlIter.PathInfo().SeqSet.Seqs {
				if c := ACounter.AlterMap[[2]int{seq.PNumber, seq.SNumber}]; c == nil {
					ACounter.AlterMap[[2]int{seq.PNumber, seq.SNumber}] = &SFCounter{}
				} else {
					c.Fail++
				}
			}
			ACounter.Mu.Unlock()
		}
	}, queries))

	if err != nil {
		log.Fatalf("visit error %v\n", err)
	}

	log.Println("fuzz ok")
}

func digest(sql string) string {
	mur32 := murmur3.New32()
	mur32.Write([]byte(sql))
	encode := mur32.Sum(nil)
	return hex.EncodeToString(encode)
}
