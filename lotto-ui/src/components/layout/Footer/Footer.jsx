import React from 'react'
import { ReactComponent as Logo } from '../../../lotto_logo.svg'
import { Link } from '../../Link'
import {ReactComponent as Github} from '../../../Octicons-mark-github.svg'

export const Footer = () => {
  return (
    <div className='bottom-0 fixed w-full'>
      <div className='bg-hellgrau text-primary text-sm grid justify-items-center space-y-8 '>
        <Link link='/'>
          <Logo title='Kleeblatt' className='w-44 align-center mt-8 p-4' />
        </Link>
        <div className='bg-hellgrau p-4 rounded space-x-8'>
          <a href='/' className=''>
            For Customers
          </a>
          <a href='/'>For Employers</a>
        </div>
        <Link link='https://github.com/crescendolll/lotto'>
          <Github title='Github Repo' className='w-44 align-center p-4' />
        </Link>
        <div className='text-xs pb-20'>
          Coded with React and Go and a lot of ❤ in Berlin ©2021 Lotto-App
        </div>
      </div>
    </div>
  )
}
