import React from 'react'

export const Link = (props) => {
    return (
        <a href = {props.link} className='hover:text-primary_neon fill-current' >
            {props.children}
        </a>
    )
}
