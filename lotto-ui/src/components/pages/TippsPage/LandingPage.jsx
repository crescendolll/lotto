import React, { useState } from 'react'
import { Header } from '../../layout/Header/Header'

export const LandingPage = ({ onLogin }) => {
  const [username, setUsername] = useState(null)
  const [password, setPassword] = useState(null)

  const onSubmit = (event) => {
    event.preventDefault()
    onLogin(username, password)
  }

  const onChangeUsername = (event) => {
    setUsername(event.target.value)
  }

  const onChangePassword = (event) => {
    setPassword(event.target.value)
  }

  return (
    <div>
      <Header />
      <form onSubmit={onSubmit}>
        <input type='text' size='20' onChange={onChangeUsername} />
        <input type='password' size='20' onChange={onChangePassword} />
        <button>Submit</button>
      </form>
    </div>
  )
}
