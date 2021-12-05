import React, { useState } from 'react'
import { TippsPage } from './components/pages/TippsPage/TippsPage'
import { login, logout, signUp } from './api'
import { LandingPage } from './components/pages/TippsPage/LandingPage'
import { EmployeePage } from './components/pages/TippsPage/EmployeePage'

function App() {
  const [auth, setAuth] = useState('')
  const [isPlayer, setIsPlayer] = useState(true)

  const onLogin = (username, password) => {
    login(username, password).then((data) => {
      setAuth(data.auth)
      if (!data.istspieler) {
        setIsPlayer(false)
      }
      console.log(data)
    })
  }

  const onLogout = () => {
    logout(auth).then((data) => {
      setAuth(null)
      console.log(data)
    })
  }

  const onSignUp = (username, password) => {
    signUp(username, password).then((data) => {
      console.log(data)
    })
  }

  return (
    <div>
      {auth ? (
        isPlayer ? (
          <TippsPage auth={auth} onLogout={onLogout} />
        ) : (
          <EmployeePage auth={auth} onLogout={onLogout}/>
        )
      ) : (
        <LandingPage
          onLogin={onLogin}
          onLogout={onLogout}
          onSignUp={onSignUp}
        />
      )}
    </div>
  )
}

export default App
