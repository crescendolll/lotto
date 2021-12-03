import React, { useState } from 'react'
import { Button } from '../../Button'
import { Footer } from '../../layout/Footer/Footer'
import { Header } from '../../layout/Header/Header'

export const LandingPage = ({ onLogin, onLogout }) => {
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
      <Header onLogout={onLogout}/>
      <form onSubmit={onSubmit} className='my-9'>
        <input type='text' size='20' onChange={onChangeUsername} className=' mx-2 border rounded-md border-dunkelgrau'/>
        <input type='password' size='20' onChange={onChangePassword} className=' mx-2 border rounded-md border-dunkelgrau'/>
        <Button>Submit</Button>
      </form>
      <Footer />
    </div>
  )
}
