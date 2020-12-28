import React, { useEffect } from 'react';
import './App.css';
import { HashRouter } from "react-router-dom";
import AppLayout from './containers/AppLayout/AppLayout';

function App() {
  useEffect(() => {

  })

  const load = async () => {
    const { instance, module } = await WebAssembly.instantiateStreaming(
      fetch(`main.wasm`), 
      window.go.importObject
    )
  }

  return (
    <HashRouter>
      <AppLayout />
    </HashRouter>
  );
}

export default App;
