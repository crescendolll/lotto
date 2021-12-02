import React, { useState, useEffect } from 'react'
import { TippsPage } from './components/pages/TippsPage/TippsPage'
import { login, logout } from './api'
import { LandingPage } from './components/pages/TippsPage/LandingPage'

function App() {
  const [auth, setAuth] = useState('21jqOg4He7nM35H2NDzOJ8BWoje')
  const onLogin = (username, password) => {
    login(username, password).then((data) => {
      setAuth(data.auth)
      console.log(data.auth)
    })
  }

  const onLogout = () => {
    logout(auth).then((data) => setAuth(null))
  }

  return (
    <div>
      {auth ? <TippsPage auth={auth} /> : <LandingPage onLogin={onLogin} onLogout={onLogout} />}
    </div>
  )
}

export default App