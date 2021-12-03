import React from 'react'
import { ReactComponent as Logo} from '../../../lotto_logo.svg'
import { Link } from '../../Link'


export const Footer = () => {
  return (
    <div className='bg-hellgrau text-primary text-sm grid justify-items-center space-y-8'>
      <Link link='/'><Logo title='Kleeblatt' className='w-44 align-center mt-8' /></Link>
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
