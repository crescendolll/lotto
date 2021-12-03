import React from 'react'

export const Button = (props) => {
  return (
    <button {...props} className='bg-primary text-hellgrau px-4 py-2 rounded m-2 hover:bg-primary_neon'>
      {props.children}
    </button>
  )
}
