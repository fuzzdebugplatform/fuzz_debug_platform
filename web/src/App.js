import React, { useState } from 'react'
import './App.css'
import { Menu } from 'antd'

import { BrowserRouter as Router, Link, Switch, Route } from 'react-router-dom'
import Home from './Home'
import 'antd/dist/antd.css'
import SQLFuzz from './SQLFuzz'
import SQLDebug2 from './SQLDebug_2'

function App() {
  let pathname = window.location.pathname.slice(1);
  if(pathname === '') {
    pathname = 'home'
  }
  const [curKey, setCurKey] = useState(pathname)
  const handleClick = e => {
    setCurKey(e.key)
  }
  return (
    <Router>
      <Menu mode="horizontal" onClick={handleClick} selectedKeys={[curKey]}>
        <Menu.Item key="home">
          <Link to="/">HOME</Link>
        </Menu.Item>

        <Menu.Item key="sqlfuzz">
          <Link to="/sqlfuzz">SQL_FUZZ</Link>
        </Menu.Item>

        <Menu.Item key="sqldebug_2">
          <Link to="/sqldebug_2">SQL_DEBUG</Link>
        </Menu.Item>
      </Menu>

      <Switch>
        <Route exact path="/">
          <Home />
        </Route>
        <Route path="/sqlfuzz">
          <SQLFuzz />
        </Route>
        <Route path="/sqldebug_2">
          <SQLDebug2 />
        </Route>
      </Switch>
    </Router>
  )
}

export default App
