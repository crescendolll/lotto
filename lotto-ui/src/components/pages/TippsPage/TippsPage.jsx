import React, { useEffect, useState } from 'react'
import { Header } from '../../layout/Header/Header'
import { Footer } from '../../layout/Footer/Footer'
import { TippsTabelle } from '../../forms/TippsTabelle/TippsTabelle'
import { TippabgabeButton } from '../../forms/TippabgabeButton/TippabgabeButton'
import { AnsichtSwitch } from '../../forms/AnsichtSwitch/AnsichtSwitch'
import { getTips, getOpenGames, getClosedGames } from '../../../api'
import { Button } from '../../Button'

export const TippsPage = ({ auth, onLogout }) => {
  // const [username, setUsername] = useState(null)

  const [tiplist, setTips] = useState([])

  const [openGames, setOpenGames] = useState([])

  useEffect(() => {
    getTips(auth, '', '').then((data) => {
      return setTips(data.statistik)
    })
  }, [])

  const onGetOpenGames = (auth) => {
    console.log(auth)
    getOpenGames(auth).then((data) => {
      return setOpenGames(data.ziehungstage)
    })
  }

  return (
    <div>
      <Header auth={auth} onLogout={onLogout}/>
      <Button onClick={() => onGetOpenGames(auth)}>get open Games</Button>
      {openGames.length > 0 ? <TippsTabelle tips={openGames} /> : null}
      <AnsichtSwitch />
      <TippsTabelle tips={tiplist} />
      <Footer />               
    </div>
  )
}
