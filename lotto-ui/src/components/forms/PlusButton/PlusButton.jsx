import React, { useState } from 'react'
import { Modal } from '../../layout/Modal'
import { login } from '../../../api'

export const PlusButton = (props) => {
  // const openModal = () => {
  //   return <Modal />
  // }
  // const [message, setMessage] = useState('no data yet')

  return (
    <button
      {...props}
      className='text-xl font-bold border border-primary py-1 px-2.5 m-2 mr-4 rounded-full bg-hellgrau hover:bg-primary_neon hover:text-hellgrau'
    >
      +
    </button>
  )
}
