import React from 'react'

export const TippZeile = (props) => {
    const tipp=props.tipp
    return (
        <div>
            <i>  {tipp.datum} </i> | 
            <b> {tipp.ziehung} </b> | 
            <i> {tipp.id} </i> 
        </div>
    )
}
