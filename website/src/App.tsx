import React from 'react';
import './App.css';
import { HashRouter } from "react-router-dom";
import AppLayout from './containers/AppLayout/AppLayout';
import { LoadingStateProvider } from './contexts/loadingState';
import Loading from './components/Loading/Loading';

function App() {

  // useEffect(() => {
  //   load()
  // }, [])
  // const load = async () => {
  //   try {
  //     const res = await runGame()
  //     console.log(res)
  //   }
  //   catch (err) {
  //     console.error(err) 
  //   }
  // }

  return (
    <LoadingStateProvider>
      <Loading />
      <HashRouter>
        <AppLayout />
      </HashRouter>
    </LoadingStateProvider>
  );
}

export default App;
