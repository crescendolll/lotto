import React from 'react'

export const Button = (props) => {
  return (
    <button className='bg-primary text-hellgrau px-4 py-2 rounded hover:bg-primary_neon'>
      {props.children}
    </button>
  )
}
