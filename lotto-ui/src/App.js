import React, { useState, useEffect } from 'react'
import { TippsPage } from './components/pages/TippsPage/TippsPage'
import { login, logout } from './api'
import { LandingPage } from './components/pages/TippsPage/LandingPage'

function App() {
  const [auth, setAuth] = useState('')
  const onLogin = (username, password) => {
    login(username, password).then((data) => {
      setAuth(data.auth)
      console.log(data.auth)
    })
  }

  const onLogout = () => {
    logout(auth).then((data) => {
      setAuth(null)
      console.log(data)
    })
  }

  return (
    <div>
      {auth ? (
        <TippsPage auth={auth} onLogout={onLogout} />
      ) : (
        <LandingPage onLogin={onLogin} onLogout={onLogout} />
      )}
    </div>
  )
}

export default App
