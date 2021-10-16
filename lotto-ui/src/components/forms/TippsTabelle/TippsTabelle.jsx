import React from 'react'
import { TippZeile } from '../TippZeile/TippZeile'
import { TippsScrollbar } from '../TippsScrollbar/TippsScrollbar'
import { TippsHeader } from '../TippsHeader/TippsHeader'
//import {nutzer, tipps, ziehungen, auszahlungen} from '../../mockdata.json'


//const Nutzer = nutzer;
// const Tipps = tipps;
// const Ziehungen = ziehungen;
// const Auszahlungen = auszahlungen;


export const TippsTabelle = (props) => {
    const tipps=props.tipps
    return (
        <div>
            <TippsScrollbar/>
            <TippsHeader/>
            {tipps.map(tipp => {return <TippZeile tipp={tipp}/>})}
        </div>
    )
}
