import React, { useState } from 'react'
import { Button } from '../../Button'
import { Footer } from '../../layout/Footer/Footer'
import { Header } from '../../layout/Header/Header'

export const LandingPage = ({ onLogin, onLogout, onSignUp }) => {
  const [username, setUsername] = useState(null)
  const [password, setPassword] = useState(null)
  const [showRegister, setShowRegister] = useState(false)

  const onSubmit = (event) => {
    event.preventDefault()
    if (showRegister) {
      onSignUp(username, password)
    }
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
      <Header onLogout={onLogout} />
      {showRegister ? (
        <div>
          <form onSubmit={onSubmit} className='my-9'>
            {'Register: '}
            <input
              type='text'
              size='20'
              onChange={onChangeUsername}
              className=' mx-2 border rounded-md border-dunkelgrau'
            />
            <input
              type='password'
              size='20'
              onChange={onChangePassword}
              className=' mx-2 border rounded-md border-dunkelgrau'
            />
            <Button>Register</Button>
          </form>
          <button onClick={() => setShowRegister(false)}>
            Click here to register
          </button>
        </div>
      ) : (
        <div>
          <form onSubmit={onSubmit} className='my-9'>
            {'Login: '}
            <input
              type='text'
              size='20'
              onChange={onChangeUsername}
              className=' mx-2 border rounded-md border-dunkelgrau'
            />
            <input
              type='password'
              size='20'
              onChange={onChangePassword}
              className=' mx-2 border rounded-md border-dunkelgrau'
            />
            <Button>Submit</Button>
          </form>
          <button onClick={() => setShowRegister(true)}>
            Click here to register
          </button>
        </div>
      )}
      <Footer />
    </div>
  )
}
