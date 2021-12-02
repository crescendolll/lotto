import React from 'react'
import { TippsTabelle } from '../../forms/TippsTabelle/TippsTabelle'
import { TippabgabeButton } from '../../forms/TippabgabeButton/TippabgabeButton'
import { AnsichtSwitch } from '../../forms/AnsichtSwitch/AnsichtSwitch'

export const TippsContainer = ({tipps, onLogin}) => {
    return (
        <div> 
            <TippabgabeButton onLogin={onLogin}/>
            <AnsichtSwitch/>
            <TippsTabelle tipps={tipps}/>
        </div>
    )
}
