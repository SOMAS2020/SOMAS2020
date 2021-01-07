import React from 'react'
import './App.css'
import { HashRouter } from 'react-router-dom'
import AppLayout from './containers/AppLayout/AppLayout'
import { LoadingStateProvider } from './contexts/loadingState'
import Loading from './components/Loading/Loading'

function App() {
  const a = 'aoeu'
  return (
    <LoadingStateProvider>
      <Loading />
      <HashRouter>
        <AppLayout />
      </HashRouter>
    </LoadingStateProvider>
  )
}

export default App
