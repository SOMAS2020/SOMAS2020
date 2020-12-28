import React, { useEffect } from 'react';
import './App.css';
import { HashRouter } from "react-router-dom";
import AppLayout from './containers/AppLayout/AppLayout';
import { go, runGame } from './wasmAPI';

function App() {

  useEffect(() => {
    load()
  }, [])
  const load = async () => {
    console.log(`load`)
    try {
      console.log(`1`)
      const { instance, module } = await WebAssembly.instantiateStreaming(
        fetch(`${process.env.PUBLIC_URL}/SOMAS2020.wasm`), 
        go.importObject
      )
      console.log(`2`)
      // await go.run(instance)
      
      console.log(`3`)
      console.log(`Going to run`)
  
      const res = await runGame()
      console.log(`Finished to run`)
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
