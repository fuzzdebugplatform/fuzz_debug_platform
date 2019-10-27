import React, { useState, useEffect, useMemo } from 'react'
import { Link } from 'react-router-dom'
import * as d3 from 'd3'

function SQLDebug() {
  const [orderBy, setOrderBy] = useState('score')
  const [scoreThresh, setScoreThresh] = useState(0.7)

  const [codeblocks, setCodeblocks] = useState([])
  const brighterScale = useMemo(() => {
    const countArr = codeblocks.map(block => block.count)
    const minCount = d3.min(countArr) || 1
    const maxCount = d3.max(countArr) || 100
    // return d3.scaleSequential(d3.interpolateOrRd).domain([minCount, maxCount])
    return d3
        .scaleLinear()
        .domain([minCount, maxCount])
        .range([0.7 * 100, 0.3 * 100])
  }, [codeblocks])
  const hueScale = useMemo(() => {
    const scoreArr = codeblocks.map(block => block.score)
    const minScore = d3.min(scoreArr) || 0
    const maxScore = d3.max(scoreArr) || 1
    return d3
        .scaleLinear()
        .domain([0.0, 1.0])
        .range([100, 0])
  }, [codeblocks])

  useEffect(() => {
    let timer
    function fetchBlocks() {
      fetch('./code_block.json')
      fetch('http://localhost:43000/codepos')
          .then(res => res.json())
          .then(data => {
            console.log(data)
            // data.sort((a, b) => (a.score < b.score ? 1 : -1))
            data = data.filter(item => item.score >= scoreThresh)
            data.sort((a, b) => {
              if (orderBy === 'score') {
                return a.score < b.score ? 1 : -1
              } else {
                return a.count < b.count ? 1 : -1
              }
            })
            console.log('sorted:', data)
            setCodeblocks(data)
          })
    }
    fetchBlocks()
    timer = setInterval(fetchBlocks, 5 * 1000)
    //timer = setInterval(fetchBlocks, 5 * 1000)
    return () => clearInterval(timer)
  }, [orderBy, scoreThresh])

  return (
      <div className="SQLDebug">
        <div className="SQLDebug-header">
          <h1>SQL Debug</h1>
          <Link to="/">home</Link>
        </div>
        <div className="SQLDebug-options">
          <div>
            排序：
            <input
                type="radio"
                name="order"
                value="score"
                checked={orderBy === 'score'}
                onChange={() => setOrderBy('score')}
            ></input>
            失败率
            <input
                type="radio"
                name="order"
                value="failed_count"
                checked={orderBy === 'failed_count'}
                onChange={() => setOrderBy('failed_count')}
            ></input>
            失败次数
          </div>
          <div>
            按 score 过滤：
            <input
                className="slider"
                type="range"
                min={0}
                max={1.0}
                step={0.1}
                onChange={e => setScoreThresh(e.target.value)}
                value={scoreThresh}
            />
            &nbsp;&nbsp;{scoreThresh}
          </div>
        </div>
        {codeblocks.map((block, _blockIdx) => (
            <div
                key={`${block.filePath}_${block.score}_${block.count}`}
                key={`${block.filePath}_${block.score}_${block.count}_${_blockIdx}`}
                className="SQLDebug-block"
            >
          <span>
            {block.filePath}: (score: {block.score}, failed count:
            {block.count})
          </span>
              {block.codeBlock.content.split('\n').map((line, lineIdx) => (
                  <div key={lineIdx + ''} className="SQLDebug-line">
                    <code>{block.codeBlock.startLine + lineIdx}:</code>
                    <pre>
                <code>
                  {line.slice(0, line.length - line.trimStart().length)}
                </code>
                <code
                    style={{
                      backgroundColor: `hsl(${hueScale(
                          block.score
                      )}, 100%, ${brighterScale(block.count)}%)`
                    }}
                >
                  {line.trimStart()}
                </code>
              </pre>
                  </div>
              ))}
            </div>
        ))}
      </div>
  )
}

export default SQLDebug

