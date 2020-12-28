import React, { useEffect } from 'react';
import './App.css';
import { HashRouter } from "react-router-dom";
import AppLayout from './containers/AppLayout/AppLayout';
import {  runGame } from './wasmAPI';

function App() {

  useEffect(() => {
    load()
  }, [])
  const load = async () => {
    try {
      const res = await runGame()
      console.log(res)
    }
    catch (err) {
      console.error(err) 
    }
  }

  return (
    <HashRouter>
      <AppLayout />
    </HashRouter>
  );
}

export default App;
