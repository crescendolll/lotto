import React from 'react'

export const NavBar = (props) => {
  return (
    <nav className='content-around'>
      <ul className='text-sm flex px-2'>
        {props.children.map((child) => (
          <li className='pr-3'>{child}</li>
        ))}
      </ul>
    </nav>
  )
}
