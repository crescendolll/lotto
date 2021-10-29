import React from 'react'
import { Header } from '../../layout/Header/Header'
import { Footer } from '../../layout/Footer/Footer'
import { TippsContainer } from '../../layout/TippsContainer/TippsContainer'

export const TippsPage = (props) => {
    const tipps = props.tipps

    return (
        <div>
            <Header/>
            <TippsContainer tipps={tipps}/>
            <Footer/>
        </div>
    )
}
