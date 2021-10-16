import React from 'react'
import './App.css'
import { TippsPage } from './components/pages/TippsPage/TippsPage'
import { tipps } from './mockdata.json'

function App() {
  return (
    <div>
      <TippsPage tipps={tipps}/>           
    </div>
  );
}

export default App;
