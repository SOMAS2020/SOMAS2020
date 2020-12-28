import React, { useEffect, useState } from 'react';
import './App.css';
import { HashRouter } from "react-router-dom";
import AppLayout from './containers/AppLayout/AppLayout';

function App() {

  useEffect(() => {
    load()
  }, [])
  const load = async () => {
    try {
      const { instance, module } = await WebAssembly.instantiateStreaming(
        fetch(`${process.env.PUBLIC_URL}/main.wasm`), 
        // @ts-ignore
        window.go.importObject
      )
      // @ts-ignore
      await window.go.run(instance)
  
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
