import React, { useState, useEffect, useMemo } from 'react'
import { Link } from 'react-router-dom'
import * as d3 from 'd3'

const hueScale = d3
  .scaleLinear()
  .domain([0, 1.0])
  .range([100, 0])

const range = (start, stop, step = 1) =>
  Array.from({ length: (stop - start) / step + 1 }, (_, i) => start + i * step)

console.log(range(256, 258))

function SQLDebug2() {
  const [codePos, setCodePos] = useState([])
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
      fetch('./new_codepos.json')
        .then(res => res.json())
        .then(codePos => {
          console.log(codePos)
          codePos = codePos.map(file => {
            return {
              filePath: file.filePath,
              codeBlocks: file.codeBlocks.filter(block => block.score >= 0.6),
              lines: file.content.split('\n')
            }
          })
          console.log(codePos)
          setCodePos(codePos)
        })
    }
    fetchCodes()
  }, [])

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

  return (
    <div className="SQLDebug">
      <div className="SQLDebug-header">
        <h1>SQL Debug</h1>
        <Link to="/">home</Link>
      </div>
      {codePos.map((file, fileIdx) => (
        <div className="SQLDebug-file" key={file.filePath}>
          <span>{file.filePath}</span>

          {file.codeBlocks.map((block, blockIdx) => (
            <div
              className="SQLDebug-block"
              key={`${file.filePath}_${blockIdx}`}
            >
              {blockIdx > 0 &&
                block.startLine - file.codeBlocks[blockIdx - 1].endLine > 1 && (
                  <div className="SQLDebug-line">
                    <code>
                      {file.codeBlocks[blockIdx - 1].endLine + 1}...
                      {block.startLine - 1}
                    </code>
                  </div>
                )}
              {/* 前面扩展两行，第一个 block 始终扩展，其余看情况 */}
              {blockIdx === 0 &&
                range(block.startLine - 2, block.startLine - 1).map(lineNum =>
                  renderLine(file, blockIdx, lineNum, false)
                )}
              {blockIdx > 0 &&
                range(block.startLine - 2, block.startLine - 1).map(lineNum =>
                  renderLine(file, blockIdx, lineNum, false)
                )}
              {/* 高亮内容 */}
              {range(block.startLine, block.endLine).map(lineNum =>
                renderLine(file, blockIdx, lineNum)
              )}
              {/* 后面扩展两行 */}
              {range(block.endLine + 1, block.endLine + 2).map(lineNum =>
                renderLine(file, blockIdx, lineNum, false)
              )}
            </div>
          ))}
        </div>
      ))}
    </div>
  )
}

export default SQLDebug2
