import React from 'react'

export const TippZeile = ({tip}) => {
    return (
        <div>
            <i>  {tip.Datum} </i> | 
            <b> {tip.Klasse} </b> | 
            <i> {tip.Id} </i> 
        </div>
    )
}
