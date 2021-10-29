import React from 'react'

export const Header = () => {
  return (
    <div className='bg-hellgrau my-9 font-bold text-primary flex justify-between px-6 items-center'>
      <div className='w-44'>
        <div className='h-16 w-16 bg-primary hover:bg-primary_neon align-center'>logo</div>
      </div>
      <nav className='content-around'>
        <ul className='text-sm flex px-2 '>
          <li className='pr-3 hover:text-primary_neon'>
            <a href='/'>Konto</a>
          </li>
          <li className='pr-3 hover:text-primary_neon'>
            <a href='/'>Tipps abgeben</a>
          </li>
          <li className='hover:text-primary_neon'>
            <a href='/'>Statistik</a>
          </li>
        </ul>
      </nav>
      <div className='w-44'> 
        <a
          className='bg-primary text-hellgrau px-4 py-2 rounded hover:bg-primary_neon'
          href='/'
        >
          Kontoverwaltung
        </a>
      </div>
    </div>
  )
}
