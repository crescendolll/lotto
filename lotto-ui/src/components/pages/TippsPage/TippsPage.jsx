import React, { useEffect, useState } from 'react'
import { Header } from '../../layout/Header/Header'
import { Footer } from '../../layout/Footer/Footer'
import { TippsTabelle } from '../../forms/TippsTabelle/TippsTabelle'
import { TippabgabeButton } from '../../forms/TippabgabeButton/TippabgabeButton'
import { AnsichtSwitch } from '../../forms/AnsichtSwitch/AnsichtSwitch'
import { getTips, getOpenGames, getClosedGames, submitTip } from '../../../api'
import { Button } from '../../Button'

export const TippsPage = ({ auth, onLogout }) => {
  // const [username, setUsername] = useState(null)

  const [tiplist, setTips] = useState([])

  const [openGames, setOpenGames] = useState([])

  const [date, setDate] = useState('')
  const [tipsubmission, setTipsubmission] = useState('')

  const [showTipSubmit, setShowTipSubmit] = useState(false)

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

  const onOpenSubmitTip = () => {
    setShowTipSubmit(true)
  }

  const onSubmit = (e) => {
    e.preventDefault()
    submitTip(auth, tipsubmission, date).then((data) => console.log(data))
    setShowTipSubmit(false)
    setDate('')
    setTipsubmission('')
  }

  const onChangeDate = (event) => {
    setDate(event.target.value)
  }

  const onChangeTipsubmission = (event) => {
    setTipsubmission(event.target.value)
  }

  return (
    <div>
      <Header auth={auth} onLogout={onLogout} />
      <Button onClick={() => onGetOpenGames(auth)}>get open Games</Button>

      {showTipSubmit ? (
        <form onSubmit={onSubmit} className='my-9 bg-grau max-w-min min-w-max p-2'>
          <Button>Submit</Button>
          <input
            type='date'
            size='20'
            onChange={onChangeDate}
            className=' mx-2 border rounded-md border-dunkelgrau'
          />
          <input
            type='text'
            size='20'
            onChange={onChangeTipsubmission}
            className=' mx-2 border rounded-md border-dunkelgrau'
          />
        </form>
      ) : (
        <div className='bg-grau p-2 max-w-min min-w-max my-9'>
          <Button onClick={onOpenSubmitTip}>+</Button> Submit a Tip
        </div>
      )}
      {openGames.length > 0 ? <TippsTabelle tips={openGames} /> : null}
      <AnsichtSwitch />
      <TippsTabelle tips={tiplist} />
      <Footer />
    </div>
  )
}
