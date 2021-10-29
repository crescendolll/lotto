import React from 'react'
import { TippsTabelle } from '../../forms/TippsTabelle/TippsTabelle'
import { TippabgabeButton } from '../../forms/TippabgabeButton/TippabgabeButton'
import { AnsichtSwitch } from '../../forms/AnsichtSwitch/AnsichtSwitch'

export const TippsContainer = (props) => {
    const tipps=props.tipps
    return (
        <div> 
            <TippabgabeButton/>
            <AnsichtSwitch/>
            <TippsTabelle tipps={tipps}/>
        </div>
    )
}
