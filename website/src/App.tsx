import React from 'react';
import './App.css';
import { HashRouter } from "react-router-dom";
import AppLayout from './containers/AppLayout/AppLayout';

function App() {
  return (
    <HashRouter>
      <AppLayout />
    </HashRouter>
  );
}

export default App;
