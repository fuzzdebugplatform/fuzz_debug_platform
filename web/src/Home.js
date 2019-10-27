import React from 'react'
import './App.css'

import { Link } from 'react-router-dom'

function Home() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>TiDB Hackathon 2019</h1>
        <h3>基于路径统计的 SQL Bug Root Cause 自动化分析</h3>
        <div>
          {/* <Link className="App-link" to="/sqlfuzz">SQL Fuzz</Link> */}
          <a className="App-link" href="./sql_fuzz.html">
            SQL Fuzz
          </a>
          <Link className="App-link" to="/sqldebug">
            SQL Debug
          </Link>
          <Link className="App-link" to="/sqldebug_2">
            SQL Debug 2
          </Link>
        </div>
      </header>
    </div>
  )
}

export default Home
