import { React, useState } from 'react'
import { Modal } from '../Modal'
import { Link } from '../../Link'
import { Button } from '../../Button'
import { NavBar } from '../../NavBar'
import { ReactComponent as Logo } from '../../../lotto_logo.svg'

export const Header = () => {
  // const [modal, setModal] = useState(false)

  const openModal = () => {
    //setModal(true)
  }
  return (
    <div>
      <div className='bg-hellgrau mb-9 font-bold text-primary flex justify-between p-6 items-center'>
        <Link link='/'><Logo title='Kleeblatt' className='w-44 align-center' /></Link>

        <NavBar>
          <Link link='/'>Konto</Link>
          <Link link='/'>Tipps abgeben</Link>
          <Link link='/'>Statistik</Link>
        </NavBar>

        <div className='w-44 align-center'>
          <Button onClick={() => openModal}>Kontoverwaltung</Button>
        </div>
      </div>
      <Modal type='Kontoverwaltung' onClose={() => openModal()}>
        <Button onClick={() => openModal(false)}>Close</Button>
      </Modal>
    </div>
  )
}
