import React from 'react'
import { Modal } from '../../layout/Modal'

export const TippabgabeButton = () => {
  const openModal = () => {
    return <Modal />
  }
  return (
    <div>
      <button onClick={() => openModal} className='border px-2 rounded-full'>
        +
      </button>
      Tipp abgeben
    </div>
  )
}
