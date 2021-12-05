import React from 'react'
import { Footer } from '../../layout/Footer/Footer'
import { Header } from '../../layout/Header/Header'

export const EmployeePage = (auth, onLogout) => {
    return (
        <div>
            <Header auth={auth} onLogout={onLogout}/>
            
            <Footer/>
        </div>
    )
}
