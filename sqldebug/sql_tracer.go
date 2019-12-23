package sqldebug

import (
	"fmt"
	"fuzz_debug_platform/config"
	"github.com/pingcap/errors"
	"io/ioutil"
	"path/filepath"
	"sort"
	"sync"
)

// CodeBlock contains information of code block
type CodeBlock struct {
	FilePath  string   `json:"filePath"`
	StartLine int64    `json:"startLine"`
	EndLine   int64    `json:"endLine"`
	Score     float64  `json:"score"`
	Count     int      `json:"count"`
	SQLs      []*string `json:"-"`
}

// String generates unique id of a codeBlock
func (codeBlock *CodeBlock) String() string {
	return formatCodeBlockUid(codeBlock.FilePath, codeBlock.StartLine, codeBlock.EndLine)
}

func formatCodeBlockUid(filepath string, startLine, endLine int64) string {
	return fmt.Sprintf("%s:%d-%d", filepath, startLine, endLine)
}

// CodeBlockCounter counters frequency of a codeBlock
type CodeBlockCounter struct {
	CodeBlock
	Counter int
}

func (counter *CodeBlockCounter) Count() {
	counter.Counter += 1
}

// SQLTracer records normal path counter and bug path counter
type SQLTracer struct {
	sync.RWMutex
	NormalPathCounters map[string]*CodeBlockCounter
	BugPathCounters    map[string]*CodeBlockCounter
}

func Trace(info *SQLTraceInfo) error {
	sqlTracer.Lock()
	defer sqlTracer.Unlock()

	var counters = sqlTracer.NormalPathCounters
	if info.IsBug {
		counters = sqlTracer.BugPathCounters
	}

	for _, blockTrace := range info.Trace {
		filePath := blockTrace.File
		for _, block := range blockTrace.Line {
			codeBlock := CodeBlock{
				FilePath:  filePath,
				StartLine: block[0],
				EndLine:   block[1],
				SQLs:      []*string{},
			}
			uid := codeBlock.String()
			if _, prs := counters[uid]; !prs {
				counters[uid] = &CodeBlockCounter{
					CodeBlock: codeBlock,
					Counter:   0,
				}
			}
			counters[uid].Count()
			if info.IsBug {
				counters[uid].SQLs = append(counters[uid].SQLs, &info.Sql)
			}
		}
	}

	return nil
}

type CodeBlocks []CodeBlock

func (c CodeBlocks) Len() int      { return len(c) }
func (c CodeBlocks) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c CodeBlocks) Less(i, j int) bool {
	return c[i].StartLine < c[j].StartLine
}

type CodeBlockPos struct {
	FilePath   string     `json:"filePath"`
	CodeBlocks CodeBlocks `json:"codeBlocks"`
	Content    string     `json:"content"`
}

type CodePos []CodeBlockPos

func (c CodePos) Len() int      { return len(c) }
func (c CodePos) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c CodePos) Less(i, j int) bool {
	iCount := 0
	jCount := 0

	for _, codeBlock := range c[i].CodeBlocks {
		iCount += codeBlock.Count
	}
	iCountDensity := float64(iCount) / float64(len(c[i].CodeBlocks))

	for _, codeBlock := range c[j].CodeBlocks {
		jCount += codeBlock.Count
	}

	jCountDensity := float64(jCount) / float64(len(c[i].CodeBlocks))
	if iCountDensity != jCountDensity {
		return iCountDensity > jCountDensity
	}
	return c[i].FilePath > c[j].FilePath
}

func Summarize() (CodePos, error) {
	sqlTracer.RLock()
	defer sqlTracer.RUnlock()

	codeMap := make(map[string]*CodeBlockPos)
	var result CodePos

	for uid, counter := range sqlTracer.BugPathCounters {
		normalCounter, ok := sqlTracer.NormalPathCounters[uid]
		var normalCount = 0

		if ok {
			normalCount = normalCounter.Counter
		}

		score := float64(counter.Counter) / float64(counter.Counter+normalCount)
		if _, prs := codeMap[counter.FilePath]; !prs {
			code, err := readSourceCode(counter.FilePath)
			if err != nil {
				return nil, errors.Trace(err)
			}
			codeMap[counter.FilePath] = &CodeBlockPos{
				FilePath:   counter.FilePath,
				CodeBlocks: []CodeBlock{},
				Content:    code,
			}
		}
		codeBlockPos := codeMap[counter.FilePath]

		codeBlockPos.CodeBlocks = append(codeBlockPos.CodeBlocks, CodeBlock{
			FilePath:  counter.FilePath,
			StartLine: counter.StartLine,
			EndLine:   counter.EndLine,
			Score:     score,
			Count:     counter.Counter,
			SQLs:      append(counter.SQLs),
		})
	}

	for _, value := range codeMap {
		sort.Sort(value.CodeBlocks)
		result = append(result, *value)
	}

	sort.Sort(result)

	if len(result) > config.GetGlobalConf().TopN {
		result = result[:config.GetGlobalConf().TopN]
	}

	return result, nil
}

func Flush() {
	sqlTracer = SQLTracer{
		NormalPathCounters: make(map[string]*CodeBlockCounter),
		BugPathCounters:    make(map[string]*CodeBlockCounter),
	}
}

func readSourceCode(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(config.GetGlobalConf().TiDBSourceDir, filePath))
	if err != nil {
		return "", errors.Trace(err)
	}
	return string(data), nil
}

//func readSourceCodeBlock(filePath string, startLine, endLine int64) (string, error) {
//	data, err := ioutil.ReadFile(config.GetGlobalConf().TiDBSourceDir + filePath)
//	if err != nil {
//		return "", errors.Trace(err)
//	}
//
//	content := string(data)
//
//	tmp := strings.Split(content, "\n")
//	tmp = tmp[startLine-1 : endLine]
//
//	return strings.Join(tmp, "\n"), nil
//}

var sqlTracer SQLTracer

func init() {
	sqlTracer = SQLTracer{
		NormalPathCounters: make(map[string]*CodeBlockCounter),
		BugPathCounters:    make(map[string]*CodeBlockCounter),
	}
}

func BugSqls(filepath string, startLine, endLine int64) []string {
	sqlTracer.RLock()
	defer sqlTracer.RUnlock()

	counter, ok := sqlTracer.BugPathCounters[formatCodeBlockUid(filepath, startLine, endLine)]
	if !ok {
		return nil
	}

	sqlStrs := make([]string, 0, len(counter.SQLs))
	for _, strPoint := range counter.SQLs {
		sqlStrs = append(sqlStrs, *strPoint)
	}

	return sqlStrs
}