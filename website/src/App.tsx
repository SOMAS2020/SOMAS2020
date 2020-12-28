import React, { useEffect } from 'react';
import './App.css';
import { HashRouter } from "react-router-dom";
import AppLayout from './containers/AppLayout/AppLayout';
import { go } from './wasm';

function App() {

  useEffect(() => {
    load()
  }, [])
  const load = async () => {
    try {
      const { instance, module } = await WebAssembly.instantiateStreaming(
        fetch(`${process.env.PUBLIC_URL}/main.wasm`), 
        go.importObject
      )
      await go.run(instance)
  
      // @ts-ignore
      window.sayHelloJS("hello")
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
