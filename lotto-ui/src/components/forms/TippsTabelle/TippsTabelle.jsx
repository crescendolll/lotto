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
        <div className='border-2 border-dashed border-primary bg-grau w-1/3'>
            <TippsScrollbar/>
            <TippsHeader/>
            {tipps.map(tipp => {return <TippZeile tipp={tipp}/>})}
        </div>
    )
}
