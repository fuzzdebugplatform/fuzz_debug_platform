import React, { useState, useEffect, useMemo } from 'react'
import * as d3 from 'd3'

const hueScale = d3
  .scaleLinear()
  .domain([0, 1.0])
  .range([100, 0])

const range = (start, stop, step = 1) =>
  Array.from({ length: (stop - start) / step + 1 }, (_, i) => start + i * step)

function SQLDebug2() {
  const [codePos, setCodePos] = useState([])
  const [scoreThresh, setScoreThresh] = useState(0.6)
  const [expands, setExpands] = useState({})

  const brighterScale = useMemo(() => {
    const countArr = codePos
      .map(file => file.codeBlocks.map(block => block.count))
      .flat()
    const min = d3.min(countArr)
    const max = d3.max(countArr)
    return d3
      .scaleLinear()
      .domain([min, max])
      .range([0.7 * 100, 0.3 * 100])
  }, [codePos])

  useEffect(() => {
    function fetchCodes() {
      //fetch('http://localhost:43000/codepos')
      fetch('/api/codepos')
        .then(res => res.json())
        .then(codePos => {
          console.log('原始：', codePos)
          // 合并重合 block
          for (let i = 0; i < codePos.length; i++) {
            const file = codePos[i]
            for (let i = file.codeBlocks.length - 1; i > 0; i--) {
              const curBlock = file.codeBlocks[i]
              const preBlock = file.codeBlocks[i - 1]
              if (
                curBlock.startLine <= preBlock.endLine &&
                curBlock.score === preBlock.score &&
                curBlock.count === preBlock.count
              ) {
                preBlock.endLine = curBlock.endLine
                curBlock.score = 0.0
              }
            }
          }
          console.log('合并后：', codePos)
          codePos = codePos.map(file => {
            return {
              filePath: file.filePath,
              codeBlocks: file.codeBlocks.filter(
                block => block.score >= scoreThresh
              ),
              lines: file.content.split('\n')
            }
          })
          console.log('过滤后：', codePos)
          setCodePos(codePos)
        })
    }
    fetchCodes()
  }, [scoreThresh])

  function renderLine(file, blockIdx, lineNum, highlight = true) {
    lineNum = lineNum - 1 // lineNum in file start from 1, not 0, so we need to substract 1
    const block = file.codeBlocks[blockIdx]
    if (lineNum < 0 || lineNum >= file.lines.length) {
      return null
    }
    return (
      <div key={`${blockIdx}_${lineNum}`} className="SQLDebug-line">
        <code>{lineNum + 1}:</code>
        <pre>
          <code>
            {file.lines[lineNum].slice(
              0,
              file.lines[lineNum].length -
                file.lines[lineNum].trimStart().length
            )}
          </code>
          <code
            style={
              highlight
                ? {
                    backgroundColor: `hsl(${hueScale(
                      block.score
                    )}, 100%, ${brighterScale(block.count)}%)`
                  }
                : {}
            }
          >
            {file.lines[lineNum].trimStart()}
          </code>
        </pre>
      </div>
    )
  }

  function genPrefixExtendRange(file, blockIdx) {
    const block = file.codeBlocks[blockIdx]
    if (blockIdx === 0) {
      return range(block.startLine - 2, block.startLine - 1)
    }
    const preBlock = file.codeBlocks[blockIdx - 1]
    if (block.startLine - preBlock.endLine <= 1) {
      return [-1]
    }
    const start = d3.max([block.startLine - 2, preBlock.endLine])
    return range(start, block.startLine - 1)
  }

  function genSuffixExtendRange(file, blockIdx) {
    const block = file.codeBlocks[blockIdx]
    if (blockIdx === file.codeBlocks.length - 1) {
      return range(block.endLine + 1, block.endLine + 2)
    }
    const nextBlock = file.codeBlocks[blockIdx + 1]
    if (nextBlock.startLine - block.endLine <= 1 + 2) {
      return [-1]
    }
    const end = d3.min([block.endLine + 2, nextBlock.startLine - 2])
    return range(block.endLine + 1, end)
  }

  function handleExpand(e, filePath) {
    e.preventDefault()
    setExpands(prevState => ({
      ...prevState,
      [filePath]: !!!prevState[filePath]
    }))
  }

  return (
    <div className="SQLDebug">
      <div className="SQLDebug-header">
        <h1>SQL Debug</h1>
      </div>
      <div className="SQLDebug-options">
        <div>
          Filter by failure rate:&nbsp;&nbsp;
          <input
            className="slider"
            type="range"
            min={0}
            max={1.0}
            step={0.1}
            onChange={e =>
              setScoreThresh(e.target.value < 0.6 ? 0.6 : e.target.value)
            }
            value={scoreThresh}
          />
          &nbsp;&nbsp;&gt;&nbsp;{scoreThresh}
        </div>
      </div>
      {codePos.map((file, fileIdx) => (
        <div className="SQLDebug-file" key={file.filePath}>
          <span>{file.filePath}</span>
          <a
            href="/"
            className="SQLDebug-expand-icon"
            onClick={e => handleExpand(e, file.filePath)}
          >
            {expands[file.filePath] === true ? ' Collapse' : ' Expand'}
          </a>

          {expands[file.filePath] === true &&
            file.codeBlocks.map((block, blockIdx) => (
              <div
                className="SQLDebug-block"
                key={`${file.filePath}_${blockIdx}`}
              >
                {blockIdx > 0 &&
                  block.startLine - file.codeBlocks[blockIdx - 1].endLine >
                    4 && (
                    <div className="SQLDebug-line">
                      <code>
                        {file.codeBlocks[blockIdx - 1].endLine + 3}...
                        {block.startLine - 3}
                      </code>
                    </div>
                  )}
                {/* 前面扩展两行，第一个 block 始终扩展，其余看情况 */}
                {genPrefixExtendRange(file, blockIdx).map(lineNum =>
                  renderLine(file, blockIdx, lineNum, false)
                )}
                {/* 高亮内容 */}
                {range(block.startLine, block.endLine).map(lineNum =>
                  renderLine(file, blockIdx, lineNum)
                )}
                {/* 后面扩展两行 */}
                {genSuffixExtendRange(file, blockIdx).map(lineNum =>
                  renderLine(file, blockIdx, lineNum, false)
                )}
                <div className="SQLDebug-tooltip">
                  <p>
                    Line {block.startLine} ~ {block.endLine}
                  </p>
                  <p>Failure rate: {block.score.toFixed(2)}</p>
                  <p>Failure count: {block.count}</p>
                </div>
              </div>
            ))}
        </div>
      ))}
    </div>
  )
}

export default SQLDebug2
