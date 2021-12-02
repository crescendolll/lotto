import React from 'react'
import { TippZeile } from '../TippZeile/TippZeile'
import { TippsScrollbar } from '../TippsScrollbar/TippsScrollbar'
import { TippsHeader } from '../TippsHeader/TippsHeader'


export const TippsTabelle = ({tips}) => {
    return (
        <div className='border-2 border-dashed border-primary bg-grau w-1/3'>
            <TippsScrollbar/>
            <TippsHeader/>
            {tips.map(tip => {return <TippZeile tip={tip}/>})}
        </div>
    )
}
