import React from 'react';
import './App.css';

import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";
import Home from './Home';
import SQLFuzz from './SQLFuzz';
import SQLDebug from './SQLDebug';
import SQLDebug2 from './SQLDebug_2';

function App() {
  return (
    <Router>
      <Switch>
        <Route exact path="/">
          <Home/>
        </Route>
        <Route path="/sqlfuzz">
          <SQLFuzz />
        </Route>
        <Route path="/sqldebug">
          <SQLDebug />
        </Route>
        <Route path="/sqldebug_2">
          <SQLDebug2 />
        </Route>
      </Switch>
    </Router>
  );
}

export default App;
