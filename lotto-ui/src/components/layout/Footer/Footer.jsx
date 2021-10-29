import React from 'react'

export const Footer = () => {
  return (
    <div className='bg-hellgrau text-primary text-sm grid justify-items-center space-y-8'>
      <div className='pt-8'>logo</div>
      <div className='bg-hellgrau p-4 rounded space-x-8'>
        <a href='/' className=''>For Customers</a>
        <a href='/'>For Employers</a>
      </div>
      <div> 
          socials
      </div>
      <div className='text-xs pb-20'>
        Coded with React and Go and a lot of ❤ in Berlin ©2021 Lotto-App
      </div>
    </div>
  )
}
